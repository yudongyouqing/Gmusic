package player

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/mewkiz/flac"
)

// 为 v1 版本固定一个输出参数，避免频繁创建 Context 造成设备异常
const (
	fixedSampleRate   = 44100
	fixedChannelCount = 2
	fixedBytesPerSamp = 2 // 16-bit
)

// Player 基于 oto v1 的播放器，支持 MP3 与 FLAC（16-bit PCM 输出）
type Player struct {
	mu           sync.Mutex
	context      *oto.Context // Context 单例（仅创建一次）
	player       *oto.Player
	playerInited bool

	currentFile *os.File
	decoder     io.Reader

	isPlaying       bool
	isPaused        bool
	currentPosition float64 // 秒（由已写入字节推算）
	duration        float64 // 秒（估算/计算）
	volume          float32 // 0.0 - 1.0
	currentFilePath string
	bytesPerSec     float64

	initialSkipBytes int64 // 首次播放/跳转时需要丢弃的 PCM 字节数

	stopCh chan struct{}
	doneCh chan struct{}
}

// NewPlayer 创建播放器，固定创建一个 Context，后续重复复用（v1 建议如此）
func NewPlayer() (*Player, error) {
	ctx, err := oto.NewContext(fixedSampleRate, fixedChannelCount, fixedBytesPerSamp, 8192)
	if err != nil {
		return nil, fmt.Errorf("创建音频上下文失败: %w", err)
	}
	return &Player{
		context:     ctx,
		volume:      1.0,
		bytesPerSec: float64(fixedSampleRate * fixedChannelCount * fixedBytesPerSamp),
	}, nil
}

// Play 兼容旧调用：从 0 秒开始
func (p *Player) Play(filePath string) error { return p.playAt(filePath, 0) }

// SeekTo 跳转到指定秒数（近似，按 PCM 字节跳过）
func (p *Player) SeekTo(sec float64) error {
	p.mu.Lock()
	path := p.currentFilePath
	p.mu.Unlock()
	if path == "" {
		return fmt.Errorf("无正在播放的文件")
	}
	if sec < 0 {
		sec = 0
	}
	return p.playAt(path, sec)
}

// playAt 播放指定文件并从 startSec 秒开始
func (p *Player) playAt(filePath string, startSec float64) error {
	p.mu.Lock()
	// 若正在播放，优雅停止并等待播放循环退出
	if p.playerInited || p.isPlaying {
		oldDone := p.doneCh
		if p.stopCh != nil {
			close(p.stopCh)
			p.stopCh = nil
		}
		p.mu.Unlock()
		if oldDone != nil {
			<-oldDone
		}
		p.mu.Lock()
	}

	// 清理旧文件句柄（若有）
	if p.currentFile != nil {
		_ = p.currentFile.Close()
		p.currentFile = nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		p.mu.Unlock()
		return fmt.Errorf("打开文件失败: %w", err)
	}
	p.currentFile = f
	p.currentFilePath = filePath

	dec, dur, _, _, err := p.getDecoder(f, filePath)
	if err != nil {
		_ = f.Close()
		p.currentFile = nil
		p.mu.Unlock()
		return err
	}
	p.decoder = dec
	p.duration = dur
	// bytesPerSec 统一采用固定输出参数
	p.bytesPerSec = float64(fixedSampleRate * fixedChannelCount * fixedBytesPerSamp)

	// 使用已创建的 Context，避免重复创建导致设备异常
	if p.context == nil {
		ctx, err := oto.NewContext(fixedSampleRate, fixedChannelCount, fixedBytesPerSamp, 8192)
		if err != nil {
			p.mu.Unlock()
			return fmt.Errorf("创建音频上下文失败: %w", err)
		}
		p.context = ctx
	}
	p.player = p.context.NewPlayer()
	p.playerInited = true

	// 初始化跳过字节数与当前位置
	if startSec > 0 && p.bytesPerSec > 0 {
		p.initialSkipBytes = int64(startSec * p.bytesPerSec)
		p.currentPosition = startSec
	} else {
		p.initialSkipBytes = 0
		p.currentPosition = 0
	}

	p.isPlaying = true
	p.isPaused = false

	// 为本次播放创建停止/完成通道
	p.stopCh = make(chan struct{})
	p.doneCh = make(chan struct{})

	// 启动播放循环
	stopCh := p.stopCh
	decReader := p.decoder
	pl := p.player
	bps := p.bytesPerSec
	p.mu.Unlock()
	go p.playLoop(stopCh, decReader, pl, bps)
	return nil
}

// getDecoder 根据扩展名选择解码器，返回 PCM io.Reader
func (p *Player) getDecoder(file *os.File, filePath string) (io.Reader, float64, int, int, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".mp3":
		dec, err := mp3.NewDecoder(file)
		if err != nil {
			return nil, 0, 0, 0, fmt.Errorf("MP3 解码失败: %w", err)
		}
		sr := dec.SampleRate()
		ch := 2 // go-mp3 输出为 2 通道 16-bit PCM
		bytesLen := float64(dec.Length())
		bps := float64(sr*ch) * 2
		var dur float64
		if bps > 0 {
			dur = bytesLen / bps
		}
		return dec, dur, sr, ch, nil
	case ".flac":
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return nil, 0, 0, 0, err
		}
		stream, err := flac.New(file)
		if err != nil {
			return nil, 0, 0, 0, fmt.Errorf("FLAC 解析失败: %w", err)
		}
		sr := int(stream.Info.SampleRate)
		ch := int(stream.Info.NChannels)
		bps := int(stream.Info.BitsPerSample)
		reader := &flacPCMReader{stream: stream, bitsPerSample: bps}
		var dur float64
		if stream.Info.NSamples > 0 && stream.Info.SampleRate > 0 {
			dur = float64(stream.Info.NSamples) / float64(stream.Info.SampleRate)
		}
		return reader, dur, sr, ch, nil
	default:
		return nil, 0, 0, 0, fmt.Errorf("不支持的音频格式: %s", ext)
	}
}

// 播放循环：在收到 stopCh 或读到 EOF/错误时退出；退出后负责清理资源并发出 doneCh
func (p *Player) playLoop(stopCh <-chan struct{}, dec io.Reader, pl *oto.Player, bps float64) {
	defer func() {
		if pl != nil {
			_ = pl.Close()
		}
		p.mu.Lock()
		p.playerInited = false
		p.player = nil
		if p.currentFile != nil {
			_ = p.currentFile.Close()
			p.currentFile = nil
		}
		p.isPlaying = false
		p.mu.Unlock()
		p.mu.Lock()
		done := p.doneCh
		p.mu.Unlock()
		if done != nil {
			close(done)
		}
	}()

	buf := make([]byte, 4096)
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		p.mu.Lock()
		paused := p.isPaused
		vol := p.volume
		skip := p.initialSkipBytes
		p.mu.Unlock()
		if paused {
			select {
			case <-stopCh:
				return
			case <-time.After(80 * time.Millisecond):
				continue
			}
		}

		n, err := dec.Read(buf)
		if n > 0 {
			// 跳过指定字节（用于 Seek）
			if skip > 0 {
				toSkip := int64(n)
				if toSkip > skip {
					toSkip = skip
				}
				slice := buf[:n]
				slice = slice[toSkip:]
				n = len(slice)
				p.mu.Lock()
				p.initialSkipBytes -= toSkip
				p.mu.Unlock()
				if n > 0 {
					copy(buf[:n], slice)
				} else { /* 本轮均被跳过 */
				}
			}

			if n > 0 {
				if vol < 1.0 {
					applyVolume16LE(buf[:n], vol)
				}
				if _, werr := pl.Write(buf[:n]); werr != nil {
					return
				}
				if bps > 0 {
					p.mu.Lock()
					p.currentPosition += float64(n) / bps
					p.mu.Unlock()
				}
			}
		}
		if err == io.EOF {
			return
		}
		if err != nil && err != io.EOF {
			return
		}
	}
}

// 控制函数
func (p *Player) Pause()  { p.mu.Lock(); p.isPaused = true; p.mu.Unlock() }
func (p *Player) Resume() { p.mu.Lock(); p.isPaused = false; p.mu.Unlock() }
func (p *Player) Stop() {
	p.mu.Lock()
	if !p.playerInited && !p.isPlaying {
		p.mu.Unlock()
		return
	}
	oldDone := p.doneCh
	if p.stopCh != nil {
		close(p.stopCh)
		p.stopCh = nil
	}
	p.mu.Unlock()
	if oldDone != nil {
		<-oldDone
	}
}
func (p *Player) SetVolume(volume float32) {
	if volume < 0 {
		volume = 0
	}
	if volume > 1 {
		volume = 1
	}
	p.mu.Lock()
	p.volume = volume
	p.mu.Unlock()
}
func (p *Player) GetCurrentPosition() float64 {
	p.mu.Lock()
	v := p.currentPosition
	p.mu.Unlock()
	return v
}
func (p *Player) GetDuration() float64 { p.mu.Lock(); v := p.duration; p.mu.Unlock(); return v }
func (p *Player) IsPlaying() bool {
	p.mu.Lock()
	v := p.isPlaying && !p.isPaused
	p.mu.Unlock()
	return v
}
func (p *Player) Close() error { p.Stop(); return nil }

// 16-bit 小端 PCM 应用音量
func applyVolume16LE(b []byte, volume float32) {
	if volume <= 0 {
		for i := range b {
			b[i] = 0
		}
		return
	}
	if volume >= 1 {
		return
	}
	for i := 0; i+1 < len(b); i += 2 {
		s := int16(binary.LittleEndian.Uint16(b[i:]))
		fv := float64(s) * float64(volume)
		if fv > math.MaxInt16 {
			fv = math.MaxInt16
		}
		if fv < math.MinInt16 {
			fv = math.MinInt16
		}
		binary.LittleEndian.PutUint16(b[i:], uint16(int16(fv)))
	}
}

// flacPCMReader 将 mewkiz/flac 流转换为交织的 16-bit PCM 字节流（不重采样）
type flacPCMReader struct {
	stream        *flac.Stream
	bitsPerSample int
	buf           []byte
	pos           int
}

func (r *flacPCMReader) Read(p []byte) (n int, err error) {
	for r.pos >= len(r.buf) {
		fr, err := r.stream.ParseNext()
		if err != nil {
			return 0, err
		}
		chs := len(fr.Subframes)
		if chs == 0 {
			return 0, io.EOF
		}
		block := len(fr.Subframes[0].Samples)
		need := chs * block * 2
		if cap(r.buf) < need {
			r.buf = make([]byte, 0, need)
		} else {
			r.buf = r.buf[:0]
		}
		shift := 0
		if r.bitsPerSample > 16 {
			shift = r.bitsPerSample - 16
		}
		for i := 0; i < block; i++ {
			for c := 0; c < chs; c++ {
				v := int32(fr.Subframes[c].Samples[i])
				if shift > 0 {
					v >>= uint(shift)
				}
				if v > math.MaxInt16 {
					v = math.MaxInt16
				}
				if v < math.MinInt16 {
					v = math.MinInt16
				}
				u := uint16(int16(v))
				r.buf = append(r.buf, byte(u), byte(u>>8))
			}
		}
		r.pos = 0
	}
	n = copy(p, r.buf[r.pos:])
	r.pos += n
	return n, nil
}
