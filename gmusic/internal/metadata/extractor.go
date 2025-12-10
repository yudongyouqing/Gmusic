package metadata

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/go-mp3"
	"github.com/mewkiz/flac"
	"github.com/yudongyouqing/GMusic/internal/lyrics/search"
	"github.com/yudongyouqing/GMusic/internal/storage"
)

// AudioInfo 音频探测信息（用于接口返回）
type AudioInfo struct {
	Format        string `json:"format"`
	Duration      int    `json:"duration"`       // 秒
	DurationText  string `json:"duration_text"`
	SampleRate    int    `json:"sample_rate"`    // Hz
	Channels      int    `json:"channels"`
	BitsPerSample int    `json:"bits_per_sample"`
	FilePath      string `json:"file_path"`
}

// ExtractMetadata 从音频文件提取元数据
func ExtractMetadata(filePath string) (*storage.Song, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	md, err := tag.ReadFrom(file)
	if err != nil {
		return nil, fmt.Errorf("读取 metadata 失败: %w", err)
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

	if pic := md.Picture(); pic != nil {
		song.CoverURL = saveCover(pic.Data, filePath)
	}

	if d := ComputeDurationSeconds(filePath); d > 0 {
		song.Duration = d
	}

	return song, nil
}

// ComputeDurationSeconds 计算音频时长（秒），支持 mp3/flac，其他返回 0
func ComputeDurationSeconds(filePath string) int {
	info, err := ProbeAudio(filePath)
	if err != nil {
		return 0
	}
	return info.Duration
}

// ProbeAudio 探测音频基础信息（时长/采样率/声道/位深）
func ProbeAudio(filePath string) (*AudioInfo, error) {
	ai := &AudioInfo{FilePath: filePath, Format: getFormat(filePath)}
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".mp3":
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
