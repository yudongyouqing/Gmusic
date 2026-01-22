package metadata

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/go-mp3"
	"github.com/mewkiz/flac"
	"github.com/yudongyouqing/GMusic/internal/lyrics/search"
	"github.com/yudongyouqing/GMusic/internal/storage"
)

// AudioInfo 音频探测信息（用于接口返回）
type AudioInfo struct {
	Format        string `json:"format"`
	Duration      int    `json:"duration"` // 秒
	DurationText  string `json:"duration_text"`
	SampleRate    int    `json:"sample_rate"` // Hz
	Channels      int    `json:"channels"`
	BitsPerSample int    `json:"bits_per_sample"`
	FilePath      string `json:"file_path"`
}

// ExtractMetadata 从音频文件提取元数据（兼容旧接口）
func ExtractMetadata(filePath string) (*storage.Song, error) {
	return ExtractMetadataWithContext(context.Background(), filePath)
}

// ExtractMetadataWithContext 从音频文件提取元数据（支持 context 取消和超时）
func ExtractMetadataWithContext(ctx context.Context, filePath string) (*storage.Song, error) {
	// 检查取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 设置默认超时（如果 context 没有超时，给一个合理的默认值）
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 打开文件（带 context 检查）
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 读取 metadata（可能耗时，需要检查 context）
	mdChan := make(chan struct {
		md  tag.Metadata
		err error
	}, 1)

	go func() {
		md, err := tag.ReadFrom(file)
		mdChan <- struct {
			md  tag.Metadata
			err error
		}{md, err}
	}()

	var md tag.Metadata
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-mdChan:
		if result.err != nil {
			return nil, fmt.Errorf("读取 metadata 失败: %w", result.err)
		}
		md = result.md
	}

	track, _ := md.Track()

	song := &storage.Song{
		Title:    md.Title(),
		Artist:   md.Artist(),
		Album:    md.Album(),
		FilePath: filePath,
		Duration: 0,
		TrackNum: track,
		Year:     md.Year(),
		Format:   getFormat(filePath),
	}

	// 再次检查取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if pic := md.Picture(); pic != nil {
		song.CoverURL = saveCover(pic.Data, filePath)
	}

	// 计算时长（可能耗时）
	if d := ComputeDurationSecondsWithContext(ctx, filePath); d > 0 {
		song.Duration = d
	}

	return song, nil
}

// ComputeDurationSeconds 计算音频时长（秒），支持 mp3/flac，其他返回 0（兼容旧接口）
func ComputeDurationSeconds(filePath string) int {
	return ComputeDurationSecondsWithContext(context.Background(), filePath)
}

// ComputeDurationSecondsWithContext 计算音频时长（支持 context）
func ComputeDurationSecondsWithContext(ctx context.Context, filePath string) int {
	select {
	case <-ctx.Done():
		return 0
	default:
	}

	info, err := ProbeAudioWithContext(ctx, filePath)
	if err != nil {
		return 0
	}
	return info.Duration
}

// ProbeAudio 探测音频基础信息（兼容旧接口）
func ProbeAudio(filePath string) (*AudioInfo, error) {
	return ProbeAudioWithContext(context.Background(), filePath)
}

// ProbeAudioWithContext 探测音频基础信息（支持 context）
func ProbeAudioWithContext(ctx context.Context, filePath string) (*AudioInfo, error) {
	ai := &AudioInfo{FilePath: filePath, Format: getFormat(filePath)}
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".mp3":
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		dec, err := mp3.NewDecoder(f)
		if err != nil {
			return nil, err
		}
		sr := dec.SampleRate()
		ch := 2
		bps := float64(sr*ch) * 2
		var sec int
		if bps > 0 {
			sec = int(float64(dec.Length()) / bps)
		}
		ai.SampleRate, ai.Channels, ai.BitsPerSample = sr, ch, 16
		ai.Duration = max0(sec)
		ai.DurationText = FormatDuration(ai.Duration)
		return ai, nil
	case ".flac":
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		stream, err := flac.New(f)
		if err != nil {
			return nil, err
		}
		sr := int(stream.Info.SampleRate)
		ch := int(stream.Info.NChannels)
		bps := int(stream.Info.BitsPerSample)
		var sec int
		if stream.Info.SampleRate > 0 {
			sec = int(stream.Info.NSamples / uint64(stream.Info.SampleRate))
		}
		ai.SampleRate, ai.Channels, ai.BitsPerSample = sr, ch, bps
		ai.Duration = max0(sec)
		ai.DurationText = FormatDuration(ai.Duration)
		return ai, nil
	default:
		ai.Duration = 0
		ai.DurationText = "00:00"
		return ai, nil
	}
}

// FormatDuration 将秒格式化为 00:00 形式
func FormatDuration(sec int) string {
	if sec < 0 {
		sec = 0
	}
	m := sec / 60
	s := sec % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}

func max0(v int) int {
	if v < 0 {
		return 0
	}
	return v
}

func getFormat(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".mp3":
		return "mp3"
	case ".flac":
		return "flac"
	case ".wav":
		return "wav"
	case ".aac":
		return "aac"
	default:
		return "unknown"
	}
}

func saveCover(data []byte, audioPath string) string {
	dir := filepath.Join(filepath.Dir(audioPath), ".covers")
	_ = os.MkdirAll(dir, 0755)

	filename := strings.TrimSuffix(filepath.Base(audioPath), filepath.Ext(audioPath)) + ".jpg"
	coverPath := filepath.Join(dir, filename)

	if err := os.WriteFile(coverPath, data, 0644); err != nil {
		return ""
	}
	return coverPath
}

func GetBitRate(filePath string) int { return 320 }

// ExtractLyrics 优先查找外部 .lrc 文件，找不到再尝试读取内嵌歌词，最后尝试网络搜索
func ExtractLyrics(song *storage.Song) (string, error) {
	// 1. 尝试外部 .lrc 文件
	lrcPath := strings.TrimSuffix(song.FilePath, filepath.Ext(song.FilePath)) + ".lrc"
	if data, err := os.ReadFile(lrcPath); err == nil {
		return string(data), nil
	}

	// 2. 尝试内嵌歌词
	file, err := os.Open(song.FilePath)
	if err == nil {
		defer file.Close()
		md, err := tag.ReadFrom(file)
		if err == nil {
			if lyrics := md.Lyrics(); lyrics != "" {
				return lyrics, nil
			}
		}
	}

	// 3. 尝试网络搜索
	if song.Title != "" && song.Artist != "" {
		lyrics, err := search.SearchNetEase(song.Title, song.Artist)
		if err == nil && lyrics != "" {
			// (可选) 缓存到本地
			_ = os.WriteFile(lrcPath, []byte(lyrics), 0644)
			return lyrics, nil
		}
	}

	return "", fmt.Errorf("歌词文件不存在，无内嵌歌词，且网络搜索失败")
}

func ExtractID3v2Info(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info := make(map[string]interface{})

	header := make([]byte, 10)
	_, err = file.Read(header)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(header[:3], []byte("ID3")) {
		version := header[3]
		flags := header[5]
		size := decodeSize(header[6:10])

		info["id3_version"] = version
		info["id3_flags"] = flags
		info["id3_size"] = size

		tagData := make([]byte, size)
		_, err = io.ReadFull(file, tagData)
		if err != nil {
			return nil, err
		}
		info["id3_data"] = string(tagData)
	}
	return info, nil
}

func decodeSize(data []byte) int {
	return (int(data[0]) << 21) | (int(data[1]) << 14) | (int(data[2]) << 7) | int(data[3])
}
