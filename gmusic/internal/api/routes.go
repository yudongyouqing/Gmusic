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
	upgrader    = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

// SetupRouter 设置路由
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
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
	router.PUT("/api/songs/:id", updateSong(db)) // 新增：更新歌曲信息

	// 播放控制 API
	router.POST("/api/player/play", playHandler())
	router.POST("/api/player/pause", pauseHandler())
	router.POST("/api/player/resume", resumeHandler())
	router.POST("/api/player/stop", stopHandler())
	router.POST("/api/player/volume", setVolumeHandler())
	router.POST("/api/player/seek", seekHandler())
	router.GET("/api/player/status", getPlayerStatus())

	// 工具 API：补全时长
	router.POST("/api/refresh/durations", refreshDurations(db))

	// 音频信息（时长等）API
	router.GET("/api/audio/:songID/info", audioInfoByID(db))
	router.POST("/api/audio/probe", audioProbeByPath())

	// 歌词/封面/扫描
	router.GET("/api/lyrics/:songID", getLyrics(db))
	router.POST("/api/scan", scanDirectory(db))
	router.GET("/api/cover/:songID", getCover(db))

	// WebSocket 实时播放状态
	router.GET("/ws/player", playerWebSocket())

	// 静态文件
	router.Static("/covers", "./covers")

	return router
}

// =========== 歌曲/搜索 ===========
func getSongs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songs, err := storage.GetAllSongs(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"total": len(songs), "songs": songs})
	}
}

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
		c.JSON(http.StatusOK, gin.H{"keyword": keyword, "total": len(songs), "songs": songs})
	}
}

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

func addSong(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FilePath string `json:"file_path" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		song, err := metadata.ExtractMetadata(req.FilePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := storage.AddSong(db, song); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, song)
	}
}

// updateSong 更新歌曲信息（目前仅 duration）
func updateSong(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var song storage.Song
		if err := db.First(&song, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "歌曲不存在"})
			return
		}

		var req struct {
			Duration *int `json:"duration"` // 使用指针以区分 0 和未提供
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Duration != nil {
			if err := db.Model(&song).Update("duration", *req.Duration).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, song)
	}
}

// =========== 播放控制 ===========
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

func seekHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Position float64 `json:"position" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Position < 0 {
			req.Position = 0
		}
		if err := audioPlayer.SeekTo(req.Position); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"position": req.Position})
	}
}

func pauseHandler() gin.HandlerFunc {
	return func(c *gin.Context) { audioPlayer.Pause(); c.JSON(http.StatusOK, gin.H{"message": "已暂停"}) }
}
func resumeHandler() gin.HandlerFunc {
	return func(c *gin.Context) { audioPlayer.Resume(); c.JSON(http.StatusOK, gin.H{"message": "已恢复"}) }
}
func stopHandler() gin.HandlerFunc {
	return func(c *gin.Context) { audioPlayer.Stop(); c.JSON(http.StatusOK, gin.H{"message": "已停止"}) }
}
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
func getPlayerStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"is_playing": audioPlayer.IsPlaying(), "position": audioPlayer.GetCurrentPosition(), "duration": audioPlayer.GetDuration()})
	}
}

// =========== 音频信息（时长等） ===========
func audioInfoByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("songID")
		var song storage.Song
		if err := db.First(&song, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "歌曲不存在"})
			return
		}
		info, err := metadata.ProbeAudio(song.FilePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, info)
	}
}

func audioProbeByPath() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FilePath string `json:"file_path" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if _, err := os.Stat(req.FilePath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("文件不存在或不可读: %v", err)})
			return
		}
		info, err := metadata.ProbeAudio(req.FilePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, info)
	}
}

// =========== 扫描/歌词/封面 ===========
func refreshDurations(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var songs []storage.Song
		if err := db.Find(&songs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		updated, skipped := 0, 0
		for i := range songs {
			if songs[i].Duration > 0 {
				skipped++
				continue
			}
			sec := metadata.ComputeDurationSeconds(songs[i].FilePath)
			if sec > 0 {
				songs[i].Duration = sec
				if err := db.Model(&songs[i]).Update("duration", sec).Error; err == nil {
					updated++
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"total": len(songs), "updated": updated, "skipped": skipped})
	}
}

func getLyrics(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songID := c.Param("songID")
		var song storage.Song
		if err := db.First(&song, songID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "歌曲不存在"})
			return
		}
		lrcContent, err := metadata.ExtractLyrics(song.FilePath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "歌词文件不存在"})
			return
		}
		lyricData, err := lyrics.ParseLRC(lrcContent)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, lyricData)
	}
}

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
		s := scanner.NewScanner(db)
		go func() {
			result, err := s.ScanDirectoryWithWorkers(req.DirPath, req.Workers)
			if err != nil {
				fmt.Printf("扫描错误: %v\n", err)
				return
			}
			fmt.Printf("扫描完成: 总文件数=%d, 添加=%d, 失败=%d\n", result.TotalFiles, result.AddedSongs, result.FailedFiles)
		}()
		c.JSON(http.StatusAccepted, gin.H{"message": "扫描已启动", "dir_path": req.DirPath})
	}
}

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
		if _, err := os.Stat(song.CoverURL); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "封面文件不存在"})
			return
		}
		c.File(song.CoverURL)
	}
}

// =========== WebSocket ===========
func playerWebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Printf("WebSocket 升级失败: %v\n", err)
			return
		}
		defer ws.Close()
		for {
			var msg map[string]any
			if err := ws.ReadJSON(&msg); err != nil {
				break
			}
			status := map[string]any{"is_playing": audioPlayer.IsPlaying(), "position": audioPlayer.GetCurrentPosition(), "duration": audioPlayer.GetDuration()}
			if err := ws.WriteJSON(status); err != nil {
				break
			}
		}
	}
}
