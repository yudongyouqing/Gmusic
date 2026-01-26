package scanner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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
	db               *gorm.DB
	supportedFormats map[string]bool
	mu               sync.Mutex
	result           *ScanResult
	ctx              context.Context    // 用于取消扫描
	cancel           context.CancelFunc // 取消函数
	pauseChan        chan struct{}      // 暂停信号
	resumeChan       chan struct{}      // 恢复信号
	isPaused         bool               // 暂停状态
	pauseMu          sync.Mutex         // 保护暂停状态
}

// NewScanner 创建新扫描器
func NewScanner(db *gorm.DB) *Scanner {
	ctx, cancel := context.WithCancel(context.Background())
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
		ctx:        ctx,
		cancel:     cancel,
		pauseChan:  make(chan struct{}),
		resumeChan: make(chan struct{}),
	}
}

// NewScannerWithContext 使用外部 context 创建扫描器（推荐）
func NewScannerWithContext(ctx context.Context, db *gorm.DB) *Scanner {
	ctx, cancel := context.WithCancel(ctx)
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
		ctx:        ctx,
		cancel:     cancel,
		pauseChan:  make(chan struct{}),
		resumeChan: make(chan struct{}),
	}
}

// Cancel 取消当前扫描
func (s *Scanner) Cancel() {
	s.cancel()
}

// Pause 暂停扫描
func (s *Scanner) Pause() {
	s.pauseMu.Lock()
	defer s.pauseMu.Unlock()
	if !s.isPaused {
		s.isPaused = true
		// 发送暂停信号（非阻塞）
		select {
		case s.pauseChan <- struct{}{}:
		default:
		}
	}
}

// Resume 恢复扫描
func (s *Scanner) Resume() {
	s.pauseMu.Lock()
	defer s.pauseMu.Unlock()
	if s.isPaused {
		s.isPaused = false
		// 发送恢复信号（非阻塞）
		select {
		case s.resumeChan <- struct{}{}:
		default:
		}
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

// processAudioFile 处理单个音频文件（兼容旧接口）
func (s *Scanner) processAudioFile(filePath string) {
	s.processAudioFileWithContext(context.Background(), filePath)
}

// processAudioFileWithContext 处理单个音频文件（支持 context 取消和超时）
func (s *Scanner) processAudioFileWithContext(ctx context.Context, filePath string) {
	// 检查取消
	select {
	case <-ctx.Done():
		return
	default:
	}

	// 检查文件是否已存在
	existing, _ := storage.GetSongByPath(s.db, filePath)
	if existing != nil && existing.ID != 0 {
		// 文件已存在，跳过
		return
	}

	// 提取元数据（支持 context 超时）
	song, err := metadata.ExtractMetadataWithContext(ctx, filePath)
	if err != nil {
		// 如果是取消错误，不记录为失败
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}
		s.mu.Lock()
		s.result.FailedFiles++
		s.result.Errors = append(s.result.Errors, fmt.Sprintf("%s: %v", filePath, err))
		s.mu.Unlock()
		return
	}

	// 再次检查取消（元数据提取可能耗时）
	select {
	case <-ctx.Done():
		return
	default:
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

// ScanDirectoryWithWorkers 使用工作池并发扫描（支持 context 取消）
func (s *Scanner) ScanDirectoryWithWorkers(ctx context.Context, dirPath string, numWorkers int) (*ScanResult, error) {
	s.result = &ScanResult{
		Errors: []string{},
	}

	// 合并外部 context 和内部 context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 监听内部取消
	go func() {
		select {
		case <-s.ctx.Done():
			cancel()
		case <-ctx.Done():
		}
	}()

	// 收集所有文件（支持取消）
	var files []string
	walkErr := make(chan error, 1)
	go func() {
		defer close(walkErr)
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			// 检查取消
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

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
		walkErr <- err
	}()

	// 等待文件收集完成或取消
	select {
	case err := <-walkErr:
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return s.result, fmt.Errorf("扫描已取消")
			}
			return nil, err
		}
	case <-ctx.Done():
		return s.result, fmt.Errorf("扫描已取消")
	}

	s.result.TotalFiles = len(files)

	// 创建工作池
	fileChan := make(chan string, numWorkers)
	var wg sync.WaitGroup

	// 启动工作 goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1) // 在 goroutine 启动前调用 Add，确保计数正确
		go func(workerID int) {
			// 使用 defer + recover 确保 Done 一定会被调用，即使 panic
			defer func() {
				if r := recover(); r != nil {
					// 记录 panic 信息，但不影响 WaitGroup 的 Done 调用
					s.mu.Lock()
					s.result.Errors = append(s.result.Errors,
						fmt.Sprintf("Worker %d panic: %v", workerID, r))
					s.mu.Unlock()
				}
				wg.Done() // 确保 Done 被调用，无论是否 panic
			}()

			for {
				select {
				case <-ctx.Done():
					return // 退出路径 1: context 取消
				case filePath, ok := <-fileChan:
					if !ok {
						return // 退出路径 2: channel 关闭
					}

					// 检查并处理暂停（持续监听直到恢复或取消）
					s.pauseMu.Lock()
					paused := s.isPaused
					s.pauseMu.Unlock()

					if paused {
						// 已暂停，等待恢复或取消
						select {
						case <-s.resumeChan:
							// 恢复，继续处理
						case <-ctx.Done():
							return
						}
					} else {
						// 检查是否有新的暂停信号（非阻塞）
						select {
						case <-s.pauseChan:
							// 收到暂停信号，更新状态并等待恢复
							s.pauseMu.Lock()
							s.isPaused = true
							s.pauseMu.Unlock()
							select {
							case <-s.resumeChan:
								s.pauseMu.Lock()
								s.isPaused = false
								s.pauseMu.Unlock()
							case <-ctx.Done():
								return
							}
						default:
							// 没有暂停信号，继续处理
						}
					}

					// 处理文件（带 context）
					s.processAudioFileWithContext(ctx, filePath)
				}
			}
		}(i) // 传递 worker ID 用于错误日志
	}

	// 分配文件到工作队列（支持取消）
	go func() {
		defer close(fileChan)
		for _, file := range files {
			select {
			case <-ctx.Done():
				return
			case fileChan <- file:
			}
		}
	}()

	// 等待所有工作完成或取消
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// 正常完成
	case <-ctx.Done():
		// 已取消，等待 worker 退出（给一个超时）
		timeout := time.NewTimer(5 * time.Second)
		defer timeout.Stop()
		select {
		case <-done:
		case <-timeout.C:
			// 超时，强制返回
		}
		return s.result, fmt.Errorf("扫描已取消")
	}

	return s.result, nil
}
