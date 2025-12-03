package lyrics

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// LyricLine 代表一行歌词
type LyricLine struct {
	Time     int64  `json:"time"`     // 毫秒
	Text     string `json:"text"`     // 歌词文本
	TimeStr  string `json:"time_str"` // 格式化时间 MM:SS.ms
}

// LyricData 代表完整歌词数据
type LyricData struct {
	Title   string       `json:"title"`
	Artist  string       `json:"artist"`
	Album   string       `json:"album"`
	Lines   []LyricLine  `json:"lines"`
}

// ParseLRC 解析 LRC 格式歌词
func ParseLRC(content string) (*LyricData, error) {
	lyricData := &LyricData{
		Lines: []LyricLine{},
	}

	lines := strings.Split(content, "\n")

	// 时间戳正则表达式 [MM:SS.ms]
	timeRegex := regexp.MustCompile(`\[(\d{2}):(\d{2})\.(\d{2,3})\]`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 解析元数据行 [ti:标题] [ar:歌手] [al:专辑]
		if strings.HasPrefix(line, "[ti:") {
			lyricData.Title = extractMetaValue(line)
			continue
		}
		if strings.HasPrefix(line, "[ar:") {
			lyricData.Artist = extractMetaValue(line)
			continue
		}
		if strings.HasPrefix(line, "[al:") {
			lyricData.Album = extractMetaValue(line)
			continue
		}

		// 解析歌词行
		matches := timeRegex.FindAllStringSubmatchIndex(line, -1)
		if matches == nil {
			continue
		}

		// 提取歌词文本（去掉所有时间戳）
		text := line
		for i := len(matches) - 1; i >= 0; i-- {
			start, end := matches[i][0], matches[i][1]
			text = text[:start] + text[end:]
		}
		text = strings.TrimSpace(text)

		// 为每个时间戳创建一行歌词
		for _, match := range matches {
			minutes, _ := strconv.ParseInt(line[match[2]:match[3]], 10, 64)
			seconds, _ := strconv.ParseInt(line[match[4]:match[5]], 10, 64)
			milliseconds, _ := strconv.ParseInt(line[match[6]:match[7]], 10, 64)

			// 处理毫秒位数（可能是 2 位或 3 位）
			if match[7]-match[6] == 2 {
				milliseconds *= 10 // 如果是 2 位，转换为 3 位
			}

			timeMs := minutes*60*1000 + seconds*1000 + milliseconds

			lyricData.Lines = append(lyricData.Lines, LyricLine{
				Time:    timeMs,
				Text:    text,
				TimeStr: formatTime(timeMs),
			})
		}
	}

	// 按时间排序
	sort.Slice(lyricData.Lines, func(i, j int) bool {
		return lyricData.Lines[i].Time < lyricData.Lines[j].Time
	})

	return lyricData, nil
}

// extractMetaValue 提取元数据值
func extractMetaValue(line string) string {
	start := strings.Index(line, ":") + 1
	end := strings.LastIndex(line, "]")
	if start > 0 && end > start {
		return strings.TrimSpace(line[start:end])
	}
	return ""
}

// formatTime 格式化时间为 MM:SS.ms
func formatTime(ms int64) string {
	totalSeconds := ms / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	milliseconds := ms % 1000

	return fmt.Sprintf("%02d:%02d.%03d", minutes, seconds, milliseconds)
}

// GetLyricAtTime 获取指定时间的歌词
func GetLyricAtTime(lyricData *LyricData, timeMs int64) *LyricLine {
	if len(lyricData.Lines) == 0 {
		return nil
	}

	// 二分查找
	left, right := 0, len(lyricData.Lines)-1
	result := &lyricData.Lines[0]

	for left <= right {
		mid := (left + right) / 2
		if lyricData.Lines[mid].Time <= timeMs {
			result = &lyricData.Lines[mid]
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return result
}

// GetLyricWindow 获取歌词窗口（当前行及前后几行）
func GetLyricWindow(lyricData *LyricData, timeMs int64, window int) []LyricLine {
	if len(lyricData.Lines) == 0 {
		return []LyricLine{}
	}

	// 找到当前行
	currentIdx := 0
	for i, line := range lyricData.Lines {
		if line.Time <= timeMs {
			currentIdx = i
		} else {
			break
		}
	}

	// 计算窗口范围
	start := currentIdx - window
	if start < 0 {
		start = 0
	}

	end := currentIdx + window + 1
	if end > len(lyricData.Lines) {
		end = len(lyricData.Lines)
	}

	return lyricData.Lines[start:end]
}

// 支持逐字歌词的扩展格式
// [00:12.00][00:17.20]歌词内容
// 可以支持更复杂的时间轴

// ParseAdvancedLRC 解析高级 LRC 格式（支持逐字）
func ParseAdvancedLRC(content string) (*LyricData, error) {
	// 这里可以扩展支持逐字歌词
	// 例如: [00:12.00]歌[00:13.00]词[00:14.00]内[00:15.00]容
	return ParseLRC(content)
}

