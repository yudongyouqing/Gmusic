package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yudongyouqing/GMusic/internal/lyrics"
	"github.com/yudongyouqing/GMusic/internal/metadata"
	"github.com/yudongyouqing/GMusic/internal/player"
	"github.com/yudongyouqing/GMusic/internal/scanner"
	"github.com/yudongyouqing/GMusic/internal/storage"
	"gorm.io/gorm"
)

var (
	audioPlayer *player.Player
	upgrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// SetupRouter 设置路由
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// 放开本地开发的 CORS（localhost/127.0.0.1:5173；也可直接全放开，便于调试）
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // 本地调试最省心；若要收敛可改为 AllowOrigins 指定 5173 源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// 初始化播放器
	var err error
	audioPlayer, err = player.NewPlayer()
	if err != nil {
		fmt.Printf("播放器初始化失败: %v\n", err)
	}

	// 歌曲相关 API
	router.GET("/api/songs", getSongs(db))
	router.GET("/api/songs/search", searchSongs(db))
	router.GET("/api/songs/:id", getSongByID(db))
	router.POST("/api/songs", addSong(db))

	// 播放控制 API（在 play 里已做详细日志与校验）
	router.POST("/api/player/play", playHandler())
	router.POST("/api/player/pause", pauseHandler())
	router.POST("/api/player/resume", resumeHandler())
	router.POST("/api/player/stop", stopHandler())
	router.POST("/api/player/volume", setVolumeHandler())
	router.GET("/api/player/status", getPlayerStatus())

	// 歌词 API
	router.GET("/api/lyrics/:songID", getLyrics(db))

	// 扫描 API
	router.POST("/api/scan", scanDirectory(db))

	// 封面 API
	router.GET("/api/cover/:songID", getCover(db))

	// WebSocket 实时播放状态
	router.GET("/ws/player", playerWebSocket())

	// 静态文件
	router.Static("/covers", "./covers")

	return router
}

// getSongs 获取所有歌曲
func getSongs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songs, err := storage.GetAllSongs(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"total": len(songs),
			"songs": songs,
		})
	}
}

// searchSongs 搜索歌曲
func searchSongs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("q")
		if keyword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
			return
		}

		songs, err := storage.SearchSongs(db, keyword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"keyword": keyword,
			"total":   len(songs),
			"songs":   songs,
		})
	}
}

// getSongByID 根据 ID 获取歌曲
func getSongByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var song storage.Song
		result := db.First(&song, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "歌曲不存在"})
			return
		}

		c.JSON(http.StatusOK, song)
	}
}

// addSong 添加歌曲
func addSong(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FilePath string `json:"file_path" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 提取元数据
		song, err := metadata.ExtractMetadata(req.FilePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 保存到数据库
		if err := storage.AddSong(db, song); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, song)
	}
}

// playHandler 播放处理（增加日志与文件存在校验）
func playHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FilePath string `json:"file_path" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("/api/player/play -> %s\n", req.FilePath)
		if req.FilePath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file_path 不能为空"})
			return
		}
		if _, err := os.Stat(req.FilePath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("文件不存在或不可读: %v", err)})
			return
		}

		if err := audioPlayer.Play(req.FilePath); err != nil {
			fmt.Printf("播放失败: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "播放开始"})
	}
}

// pauseHandler 暂停处理
func pauseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		audioPlayer.Pause()
		c.JSON(http.StatusOK, gin.H{"message": "已暂停"})
	}
}

// resumeHandler 恢复处理
func resumeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		audioPlayer.Resume()
		c.JSON(http.StatusOK, gin.H{"message": "已恢复"})
	}
}

// stopHandler 停止处理
func stopHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		audioPlayer.Stop()
		c.JSON(http.StatusOK, gin.H{"message": "已停止"})
	}
}

// setVolumeHandler 设置音量
func setVolumeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Volume float32 `json:"volume" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		audioPlayer.SetVolume(req.Volume)
		c.JSON(http.StatusOK, gin.H{"volume": req.Volume})
	}
}

// getPlayerStatus 获取播放器状态
func getPlayerStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"is_playing": audioPlayer.IsPlaying(),
			"position":   audioPlayer.GetCurrentPosition(),
			"duration":   audioPlayer.GetDuration(),
		})
	}
}

// getLyrics 获取歌词
func getLyrics(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songID := c.Param("songID")

		var song storage.Song
		if err := db.First(&song, songID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "歌曲不存在"})
			return
		}

		// 尝试读取 .lrc 文件
		lrcContent, err := metadata.ExtractLyrics(song.FilePath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "歌词文件不存在"})
			return
		}

		// 解析 LRC
		lyricData, err := lyrics.ParseLRC(lrcContent)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, lyricData)
	}
}

// scanDirectory 扫描目录
func scanDirectory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			DirPath string `json:"dir_path" binding:"required"`
			Workers int    `json:"workers"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Workers == 0 {
			req.Workers = 4
		}

		// 创建扫描器
		s := scanner.NewScanner(db)

		// 异步扫描
		go func() {
			result, err := s.ScanDirectoryWithWorkers(req.DirPath, req.Workers)
			if err != nil {
				fmt.Printf("扫描错误: %v\n", err)
				return
			}

			fmt.Printf("扫描完成: 总文件数=%d, 添加=%d, 失败=%d\n",
				result.TotalFiles, result.AddedSongs, result.FailedFiles)
		}()

		c.JSON(http.StatusAccepted, gin.H{
			"message":  "扫描已启动",
			"dir_path": req.DirPath,
		})
	}
}

// getCover 获取封面
func getCover(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songID := c.Param("songID")

		var song storage.Song
		if err := db.First(&song, songID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "歌曲不存在"})
			return
		}

		if song.CoverURL == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "封面不存在"})
			return
		}

		// 检查文件是否存在
		if _, err := os.Stat(song.CoverURL); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "封面文件不存在"})
			return
		}

		c.File(song.CoverURL)
	}
}

// playerWebSocket WebSocket 实时播放状态
func playerWebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Printf("WebSocket 升级失败: %v\n", err)
			return
		}
		defer ws.Close()

		for {
			// 读取客户端消息
			var msg map[string]interface{}
			if err := ws.ReadJSON(&msg); err != nil {
				break
			}

			// 发送播放器状态
			status := map[string]interface{}{
				"is_playing": audioPlayer.IsPlaying(),
				"position":   audioPlayer.GetCurrentPosition(),
				"duration":   audioPlayer.GetDuration(),
			}

			if err := ws.WriteJSON(status); err != nil {
				break
			}
		}
	}
}
