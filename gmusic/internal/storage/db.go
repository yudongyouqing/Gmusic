// Package storage 提供 GMusic 的持久化存储层实现。
//
// 职责概览：
// 1) 定义核心数据模型（Song、Playlist、PlayHistory）。
// 2) 封装数据库初始化（基于 GORM，默认使用本地 SQLite 文件）。
// 3) 提供常用的数据访问方法（查询、搜索、新增等）。
//
// 说明：当前实现面向单机使用场景，SQLite 简单轻量；若需更高并发或分布式，
package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Song 表示一首歌曲的元数据记录。
// 注意：
// - FilePath 建议在数据库层面设置唯一约束以避免重复导入（当前模型未强制）。
// - 如需区分“未知值”和“明确为 0/空值”，可将部分字段改为指针或使用 sql.NullXxx。
type Song struct {
	ID       uint   `gorm:"primaryKey" json:"id"` // 主键 ID
	Title    string `json:"title"`                // 歌曲标题
	Artist   string `json:"artist"`               // 艺术家/歌手
	Album    string `json:"album"`                // 专辑名
	FilePath string `json:"file_path"`            // 音频文件的绝对/相对路径
	Duration int    `json:"duration"`             // 时长（秒）
	BitRate  int    `json:"bit_rate"`             // 比特率（kbps）
	Format   string `json:"format"`               // 文件格式（如 mp3、flac、wav）
	CoverURL string `json:"cover_url"`            // 封面图片路径或 URL
	TrackNum int    `json:"track_num"`            // 专辑内的曲目序号
	Year     int    `json:"year"`                 // 发行年份
}

// Playlist 表示一个播放列表，Songs 通过 many2many 中间表 playlist_songs 关联。
// 建议：为中间表 (playlist_id, song_id) 添加唯一复合索引以避免重复加入同一歌曲。
type Playlist struct {
	ID    uint   `gorm:"primaryKey" json:"id"`                   // 主键 ID
	Name  string `json:"name"`                                   // 播放列表名称
	Songs []Song `gorm:"many2many:playlist_songs;" json:"songs"` // 列表包含的歌曲，多对多关系（中间表 playlist_songs）
}

// PlayHistory 记录歌曲的播放历史。
// 建议：为 SongID、PlayedAt 建索引，便于按歌曲/时间范围查询。
// 可选：增加外键约束（SQLite 下外键默认关闭，需要 PRAGMA foreign_keys=ON）。
type PlayHistory struct {
	ID       uint  `gorm:"primaryKey" json:"id"` // 主键 ID
	SongID   uint  `json:"song_id"`              // 被播放的歌曲 ID（外键）
	PlayedAt int64 `json:"played_at"`            // 播放时间（Unix 时间戳，秒）
}

// InitDB 初始化数据库连接并进行自动迁移，返回 *gorm.DB。
// 当前使用 sqlite 驱动，数据库保存在给定的文件路径（如 gmusic.db）中。
// 注意：AutoMigrate 只做“无损”变更，复杂 Schema 变更请使用显式迁移工具。
func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移：若表不存在则创建，字段缺失则补齐，不会删除列。
	err = db.AutoMigrate(&Song{}, &Playlist{}, &PlayHistory{})
	if err != nil {
		return nil, err
	}

	// 性能提示：如需更好的读并发，可在应用启动时开启 WAL 模式（SQLite 专有）。
	// _, _ = db.DB() // 获取 *sql.DB 后可执行原生 PRAGMA，例如：
	// db.Exec("PRAGMA journal_mode = WAL;")

	return db, nil
}

// GetAllSongs 返回数据库中的所有歌曲。
// 提示：当数据量较大时，建议在调用方增加分页（LIMIT/OFFSET 或基于主键的游标分页）。
func GetAllSongs(db *gorm.DB) ([]Song, error) {
	var songs []Song
	result := db.Find(&songs)
	return songs, result.Error
}

// SearchSongs 基于关键字在标题、艺术家、专辑上做模糊匹配（%keyword%）。
// 注意：前后皆为 % 的 LIKE 查询无法利用普通 B-Tree 索引，数据增大后性能会下降。
// 可选方案：
// - 使用 SQLite FTS5 建全文索引表以获得更好的搜索性能；
// - 或改为前缀匹配（keyword%）以便部分利用索引（语义会变化）。
func SearchSongs(db *gorm.DB, keyword string) ([]Song, error) {
	var songs []Song
	result := db.Where("title LIKE ? OR artist LIKE ? OR album LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Find(&songs)
	return songs, result.Error
}

// AddSong 插入一条歌曲记录。
// 说明：函数本身不做去重判断，若需避免重复，请在 FilePath 上加唯一约束或在写入前查重。
func AddSong(db *gorm.DB, song *Song) error {
	return db.Create(song).Error
}

// GetSongByPath 根据文件路径查询歌曲，若不存在返回 gorm.ErrRecordNotFound。
// 建议：为 file_path 建索引/唯一约束以提升查询性能与数据一致性。
func GetSongByPath(db *gorm.DB, filePath string) (*Song, error) {
	var song Song
	result := db.Where("file_path = ?", filePath).First(&song)
	return &song, result.Error
}
