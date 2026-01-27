package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	neteaseSearchAPI = "http://music.163.com/api/search/get/"
	neteaseLyricAPI  = "http://music.163.com/api/song/lyric"
)

// NetEaseSearchResponse 搜索接口的响应结构
type NetEaseSearchResponse struct {
	Result struct {
		Songs []struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"songs"`
	} `json:"result"`
}

// NetEaseLyricResponse 歌词接口的响应结构
type NetEaseLyricResponse struct {
	Lrc struct {
		Lyric string `json:"lyric"`
	} `json:"lrc"`
}

// SearchNetEase 在网易云音乐搜索歌词
func SearchNetEase(title, artist string) (string, error) {
	// 1. 搜索歌曲ID
	id, err := searchSongID(title, artist)
	if err != nil {
		return "", fmt.Errorf("在网易云音乐中找不到匹配的歌曲: %w", err)
	}

	// 2. 根据ID获取歌词
	return getLyricByID(id)
}

// searchSongID 根据歌曲名和歌手名搜索歌曲ID
func searchSongID(title, artist string) (int, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	formData := url.Values{}
	formData.Set("s", title+" "+artist)
	formData.Set("type", "1") // 1 表示搜索歌曲
	formData.Set("limit", "1")

	req, _ := http.NewRequest("POST", neteaseSearchAPI, strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var searchResult NetEaseSearchResponse
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return 0, err
	}

	if len(searchResult.Result.Songs) > 0 {
		return searchResult.Result.Songs[0].ID, nil
	}

	return 0, fmt.Errorf("未找到歌曲")
}

// getLyricByID 根据歌曲ID获取歌词
func getLyricByID(id int) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	uri := fmt.Sprintf("%s?id=%d&lv=1&kv=1&tv=-1", neteaseLyricAPI, id)

	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("Referer", "http://music.163.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var lyricResult NetEaseLyricResponse
	if err := json.Unmarshal(body, &lyricResult); err != nil {
		return "", err
	}

	if lyricResult.Lrc.Lyric != "" {
		return lyricResult.Lrc.Lyric, nil
	}

	return "", fmt.Errorf("未找到歌词")
}
