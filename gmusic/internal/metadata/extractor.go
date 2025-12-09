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
	"github.com/yudongyouqing/GMusic/internal/storage"
)

// ExtractMetadata 从音频文件提取元数据
func ExtractMetadata(filePath string) (*storage.Song, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 使用 dhowden/tag 解析（该库不提供时长）
	md, err := tag.ReadFrom(file)
	if err != nil {
		return nil, fmt.Errorf("读取 metadata 失败: %w", err)
	}

	track, _ := md.Track() // 只取当前 track 序号

	song := &storage.Song{
		Title:    md.Title(),
		Artist:   md.Artist(),
		Album:    md.Album(),
		FilePath: filePath,
		Duration: 0, // 先置 0，下面尝试计算
		TrackNum: track,
		Year:     md.Year(),
		Format:   getFormat(filePath),
	}

	// 提取封面
	if pic := md.Picture(); pic != nil {
		song.CoverURL = saveCover(pic.Data, filePath)
	}

	// 补充计算时长
	if d := ComputeDurationSeconds(filePath); d > 0 {
		song.Duration = d
	}

	return song, nil
}

// ComputeDurationSeconds 计算音频时长（秒），支持 mp3/flac，其他返回 0
func ComputeDurationSeconds(filePath string) int {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".mp3":
		f, err := os.Open(filePath)
		if err != nil {
			return 0
		}
		defer f.Close()
		dec, err := mp3.NewDecoder(f)
		if err != nil {
			return 0
		}
		sr := float64(dec.SampleRate())
		ch := 2.0
		bps := sr * ch * 2.0
		if bps <= 0 {
			return 0
		}
		sec := int(float64(dec.Length()) / bps)
		if sec < 0 { sec = 0 }
		return sec
	case ".flac":
		f, err := os.Open(filePath)
		if err != nil { return 0 }
		defer f.Close()
		stream, err := flac.New(f)
		if err != nil { return 0 }
		if stream.Info.SampleRate == 0 { return 0 }
		sec := int(stream.Info.NSamples / uint64(stream.Info.SampleRate))
		if sec < 0 { sec = 0 }
		return sec
	default:
		return 0
	}
}

// getFormat 获取文件格式
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

// saveCover 保存封面图片
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

// GetBitRate 获取比特率（占位实现）
func GetBitRate(filePath string) int {
	return 320
}

// ExtractLyrics 从文件提取歌词（如果有 .lrc 文件）
func ExtractLyrics(audioPath string) (string, error) {
	lrcPath := strings.TrimSuffix(audioPath, filepath.Ext(audioPath)) + ".lrc"
	data, err := os.ReadFile(lrcPath)
	if err != nil {
		return "", fmt.Errorf("歌词文件不存在: %w", err)
	}
	return string(data), nil
}

// 扩展：从 ID3v2 标签中提取更多信息（可选）
func ExtractID3v2Info(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info := make(map[string]interface{})

	// 读取 ID3v2 头
	header := make([]byte, 10)
	_, err = file.Read(header)
	if err != nil {
		return nil, err
	}

	// 检查 ID3v2 标识
	if bytes.Equal(header[:3], []byte("ID3")) {
		version := header[3]
		flags := header[5]
		size := decodeSize(header[6:10])

		info["id3_version"] = version
		info["id3_flags"] = flags
		info["id3_size"] = size

		// 读取 ID3v2 数据
		tagData := make([]byte, size)
		_, err = io.ReadFull(file, tagData)
		if err != nil {
			return nil, err
		}
		info["id3_data"] = string(tagData)
	}
	return info, nil
}

// decodeSize 解码 ID3v2 大小字段（7 位编码）
func decodeSize(data []byte) int {
	return (int(data[0]) << 21) | (int(data[1]) << 14) | (int(data[2]) << 7) | int(data[3])
}
