package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Song 代表一首歌曲
type Song struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	FilePath string `json:"file_path"`
	Duration int    `json:"duration"` // 秒
	BitRate  int    `json:"bit_rate"`
	Format   string `json:"format"` // mp3, flac, wav
	CoverURL string `json:"cover_url"`
	TrackNum int    `json:"track_num"`
	Year     int    `json:"year"`
}

// Playlist 代表播放列表
type Playlist struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Songs []Song `gorm:"many2many:playlist_songs;" json:"songs"`
}

// PlayHistory 代表播放历史
type PlayHistory struct {
	ID       uint  `gorm:"primaryKey" json:"id"`
	SongID   uint  `json:"song_id"`
	PlayedAt int64 `json:"played_at"` // Unix timestamp
}

// InitDB 初始化数据库
func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	err = db.AutoMigrate(&Song{}, &Playlist{}, &PlayHistory{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetAllSongs 获取所有歌曲
func GetAllSongs(db *gorm.DB) ([]Song, error) {
	var songs []Song
	result := db.Find(&songs)
	return songs, result.Error
}

// SearchSongs 搜索歌曲
func SearchSongs(db *gorm.DB, keyword string) ([]Song, error) {
	var songs []Song
	result := db.Where("title LIKE ? OR artist LIKE ? OR album LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Find(&songs)
	return songs, result.Error
}

// AddSong 添加歌曲到数据库
func AddSong(db *gorm.DB, song *Song) error {
	return db.Create(song).Error
}

// GetSongByPath 根据文件路径获取歌曲
func GetSongByPath(db *gorm.DB, filePath string) (*Song, error) {
	var song Song
	result := db.Where("file_path = ?", filePath).First(&song)
	return &song, result.Error
}
