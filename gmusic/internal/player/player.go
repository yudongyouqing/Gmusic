package player

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"sync"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

// Player 基于 oto v1 的简单播放器（当前实现 MP3）
type Player struct {
	mu              sync.Mutex
	currentFile     *os.File
	decoder         io.Reader
	context         *oto.Context
	player          *oto.Player
	playerInited    bool
	isPlaying       bool
	isPaused        bool
	currentPosition float64 // 秒（由已写入字节推算）
	duration        float64 // 秒（估算）
	volume          float32 // 0.0 - 1.0
	currentFilePath string
	bytesPerSec     float64
	bytesWritten    int64
}

// NewPlayer 创建播放器上下文
func NewPlayer() (*Player, error) {
	// 44100Hz, 2 声道, 16-bit, 缓冲区 8192 字节
	ctx, err := oto.NewContext(44100, 2, 2, 8192)
	if err != nil {
		return nil, fmt.Errorf("创建音频上下文失败: %w", err)
	}
	return &Player{
		context: ctx,
		volume:  1.0,
	}, nil
}

// Play 播放文件（目前仅 MP3）
func (p *Player) Play(filePath string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 停止并清理旧的
	if p.playerInited {
		_ = p.player.Close()
		p.playerInited = false
	}
	if p.currentFile != nil {
		_ = p.currentFile.Close()
		p.currentFile = nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	p.currentFile = f
	p.currentFilePath = filePath

	// 仅实现 MP3 解码
	dec, err := mp3.NewDecoder(f)
	if err != nil {
		_ = f.Close()
		p.currentFile = nil
		return fmt.Errorf("MP3 解码失败: %w", err)
	}
	p.decoder = dec

	// 估算参数
	p.bytesPerSec = float64(dec.SampleRate()) * 2 /*channels*/ * 2 /*bytes*/
	if p.bytesPerSec > 0 {
		p.duration = float64(dec.Length()) / p.bytesPerSec
	} else {
		p.duration = 0
	}
	p.bytesWritten = 0
	p.currentPosition = 0
	p.isPlaying = true
	p.isPaused = false

	// 新建底层播放器
	pl := p.context.NewPlayer()
	p.player = pl
	p.playerInited = true

	go p.playLoop()
	return nil
}

// 内部播放循环
func (p *Player) playLoop() {
	buf := make([]byte, 4096)
	for {
		p.mu.Lock()
		playing := p.isPlaying
		paused := p.isPaused
		dec := p.decoder
		pl := p.player
		inited := p.playerInited
		vol := p.volume
		bps := p.bytesPerSec
		p.mu.Unlock()

		if !playing || !inited || dec == nil {
			break
		}
		if paused {
			// 简单暂停：停止写数据
			// 也可 sleep 以降低 CPU
			// 这里不更新进度
			continue
		}

		n, err := dec.Read(buf)
		if n > 0 {
			// 应用音量（16-bit LE）
			if vol < 1.0 {
				applyVolume16LE(buf[:n], vol)
			}
			if _, werr := pl.Write(buf[:n]); werr != nil {
				break
			}
			p.mu.Lock()
			p.bytesWritten += int64(n)
			if bps > 0 {
				p.currentPosition = float64(p.bytesWritten) / bps
			}
			p.mu.Unlock()
		}
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			break
		}
	}

	p.mu.Lock()
	p.isPlaying = false
	if p.playerInited {
		_ = p.player.Close()
		p.playerInited = false
	}
	if p.currentFile != nil {
		_ = p.currentFile.Close()
		p.currentFile = nil
	}
	p.mu.Unlock()
}

// applyVolume16LE 对 16-bit 小端 PCM 应用音量
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

// Pause 暂停
func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isPaused = true
}

// Resume 恢复
func (p *Player) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isPaused = false
}

// Stop 停止
func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isPlaying = false
	p.isPaused = false
	p.currentPosition = 0
	if p.playerInited {
		_ = p.player.Close()
		p.playerInited = false
	}
	if p.currentFile != nil {
		_ = p.currentFile.Close()
		p.currentFile = nil
	}
}

// SetVolume 设置音量 0..1
func (p *Player) SetVolume(volume float32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if volume < 0 {
		volume = 0
	}
	if volume > 1 {
		volume = 1
	}
	p.volume = volume
}

// GetCurrentPosition 当前播放位置（秒）
func (p *Player) GetCurrentPosition() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.currentPosition
}

// GetDuration 总时长（秒）
func (p *Player) GetDuration() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.duration
}

// IsPlaying 是否正在播放
func (p *Player) IsPlaying() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.isPlaying && !p.isPaused
}

// Close 关闭
func (p *Player) Close() error {
	p.Stop()
	return nil
}
