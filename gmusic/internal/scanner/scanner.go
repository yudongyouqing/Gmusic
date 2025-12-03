package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yudongyouqing/GMusic/internal/metadata"
	"github.com/yudongyouqing/GMusic/internal/storage"
	"gorm.io/gorm"
)

// ScanResult 扫描结果
type ScanResult struct {
	TotalFiles  int
	AddedSongs  int
	FailedFiles int
	Errors      []string
}

// Scanner 目录扫描器
type Scanner struct {
	db              *gorm.DB
	supportedFormats map[string]bool
	mu              sync.Mutex
	result          *ScanResult
}

// NewScanner 创建新扫描器
func NewScanner(db *gorm.DB) *Scanner {
	return &Scanner{
		db: db,
		supportedFormats: map[string]bool{
			".mp3":  true,
			".flac": true,
			".wav":  true,
			".aac":  true,
		},
		result: &ScanResult{
			Errors: []string{},
		},
	}
}

// ScanDirectory 扫描目录
func (s *Scanner) ScanDirectory(dirPath string) (*ScanResult, error) {
	s.result = &ScanResult{
		Errors: []string{},
	}

	// 验证目录存在
	info, err := os.Stat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("目录不存在: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("路径不是目录")
	}

	// 递归扫描
	err = s.walkDirectory(dirPath)
	if err != nil {
		return nil, err
	}

	return s.result, nil
}

// walkDirectory 递归遍历目录
func (s *Scanner) walkDirectory(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("读取目录失败: %w", err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			// 递归进入子目录
			s.walkDirectory(fullPath)
		} else {
			// 检查是否是支持的格式
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if s.supportedFormats[ext] {
				s.result.TotalFiles++
				s.processAudioFile(fullPath)
			}
		}
	}

	return nil
}

// processAudioFile 处理单个音频文件
func (s *Scanner) processAudioFile(filePath string) {
	// 检查文件是否已存在
	existing, _ := storage.GetSongByPath(s.db, filePath)
	if existing != nil && existing.ID != 0 {
		// 文件已存在，跳过
		return
	}

	// 提取元数据
	song, err := metadata.ExtractMetadata(filePath)
	if err != nil {
		s.mu.Lock()
		s.result.FailedFiles++
		s.result.Errors = append(s.result.Errors, fmt.Sprintf("%s: %v", filePath, err))
		s.mu.Unlock()
		return
	}

	// 保存到数据库
	err = storage.AddSong(s.db, song)
	if err != nil {
		s.mu.Lock()
		s.result.FailedFiles++
		s.result.Errors = append(s.result.Errors, fmt.Sprintf("保存失败 %s: %v", filePath, err))
		s.mu.Unlock()
		return
	}

	s.mu.Lock()
	s.result.AddedSongs++
	s.mu.Unlock()

	fmt.Printf("✅ 已添加: %s - %s\n", song.Artist, song.Title)
}

// ScanDirectoryAsync 异步扫描目录
func (s *Scanner) ScanDirectoryAsync(dirPath string, callback func(*ScanResult)) {
	go func() {
		result, err := s.ScanDirectory(dirPath)
		if err != nil {
			fmt.Printf("扫描错误: %v\n", err)
			return
		}
		callback(result)
	}()
}

// ScanDirectoryWithWorkers 使用工作池并发扫描
func (s *Scanner) ScanDirectoryWithWorkers(dirPath string, numWorkers int) (*ScanResult, error) {
	s.result = &ScanResult{
		Errors: []string{},
	}

	// 收集所有文件
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if s.supportedFormats[ext] {
				files = append(files, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.result.TotalFiles = len(files)

	// 创建工作池
	fileChan := make(chan string, numWorkers)
	var wg sync.WaitGroup

	// 启动工作 goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				s.processAudioFile(filePath)
			}
		}()
	}

	// 分配文件到工作队列
	go func() {
		for _, file := range files {
			fileChan <- file
		}
		close(fileChan)
	}()

	// 等待所有工作完成
	wg.Wait()

	return s.result, nil
}

