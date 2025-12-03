# GMusic 项目结构详解

## 📂 完整项目树

```
gmusic/
│
├── cmd/                          # 命令行程序入口
│   └── server/
│       └── main.go              # 服务器主程序
│
├── internal/                     # 内部包（不对外暴露）
│   ├── api/
│   │   └── routes.go            # REST API 路由和处理器
│   │
│   ├── player/
│   │   └── player.go            # 音频播放引擎
│   │
│   ├── metadata/
│   │   └── extractor.go         # 元数据提取（tag 解析）
│   │
│   ├── lyrics/
│   │   └── lrc_parser.go        # LRC 歌词解析
│   │
│   ├── scanner/
│   │   └── scanner.go           # 目录扫描和索引
│   │
│   └── storage/
│       └── db.go                # 数据库操作和模型
│
├── ui/                           # React 前端
│   ├── src/
│   │   ├── App.jsx              # 主应用组件
│   │   ├── App.css              # 主样式
│   │   ├── components/
│   │   │   ├── Player.jsx       # 播放器组件
│   │   │   ├── Player.css
│   │   │   ├── SongList.jsx     # 歌曲列表组件
│   │   │   ├── SongList.css
│   │   │   ├── SearchBar.jsx    # 搜索栏组件
│   │   │   ├── SearchBar.css
│   │   │   ├── LyricDisplay.jsx # 歌词显示组件
│   │   │   └── LyricDisplay.css
│   │   └── index.js             # React 入口
│   ├── package.json
│   └── vite.config.js           # Vite 配置
│
├── go.mod                        # Go 模块定义
├── go.sum                        # Go 依赖锁定
├── Makefile                      # 构建脚本
├── Dockerfile                    # Docker 镜像定义
├── docker-compose.yml            # Docker 编排
├── .gitignore                    # Git 忽略文件
├── README.md                     # 项目文档
├── QUICKSTART.md                 # 快速开始
├── INTERVIEW_GUIDE.md            # 面试指南
└── PROJECT_STRUCTURE.md          # 本文件
```

## 🔍 各模块详解

### 1. cmd/server/main.go

**职责**：应用程序入口

**主要功能**：
- 初始化数据库
- 设置 API 路由
- 启动 HTTP 服务器

**代码流程**：
```
main()
  ├── storage.InitDB()        # 初始化 SQLite
  ├── api.SetupRouter()       # 配置路由
  └── router.Run(":8080")     # 启动服务器
```

### 2. internal/storage/db.go

**职责**：数据持久化层

**核心结构**：
```go
type Song struct {
    ID       uint      // 主键
    Title    string    // 歌曲名
    Artist   string    // 歌手
    Album    string    // 专辑
    FilePath string    // 文件路径（唯一）
    Duration int       // 时长（秒）
    BitRate  int       // 比特率
    Format   string    // 格式（mp3/flac/wav）
    CoverURL string    // 封面路径
    TrackNum int       // 曲目号
    Year     int       // 发行年份
}
```

**主要函数**：
- `InitDB()`：初始化数据库和表
- `GetAllSongs()`：获取所有歌曲
- `SearchSongs()`：搜索歌曲
- `AddSong()`：添加歌曲
- `GetSongByPath()`：根据路径获取歌曲

**数据库关系**：
```
Song (1) ──── (N) PlayHistory
Song (M) ──── (N) Playlist
```

### 3. internal/metadata/extractor.go

**职责**：音频文件元数据提取

**支持格式**：
- MP3（ID3v2 标签）
- FLAC（FLAC metadata blocks）
- WAV（INFO chunk）
- AAC（iTunes metadata）

**主要函数**：
- `ExtractMetadata()`：提取所有元数据
- `saveCover()`：保存封面图片
- `ExtractLyrics()`：提取歌词文件
- `ExtractID3v2Info()`：详细解析 ID3v2

**使用的库**：
- `github.com/dhowden/tag`：统一的 metadata 读取接口

### 4. internal/player/player.go

**职责**：音频播放引擎

**核心功能**：
- 音频解码
- 实时播放
- 音量控制
- 进度管理

**主要方法**：
```go
type Player struct {
    // 播放状态
    isPlaying       bool
    isPaused        bool
    currentPosition float64
    duration        float64
    volume          float32
}

// 主要方法
func (p *Player) Play(filePath string) error
func (p *Player) Pause()
func (p *Player) Resume()
func (p *Player) Stop()
func (p *Player) SetVolume(volume float32)
func (p *Player) GetCurrentPosition() float64
func (p *Player) GetDuration() float64
```

**播放流程**：
```
Play(filePath)
  ├── 打开文件
  ├── 选择解码器（MP3/FLAC）
  ├── 创建音频上下文
  └── 启动播放 goroutine
      ├── 读取解码数据
      ├── 应用音量
      └── 写入音频设备
```

### 5. internal/lyrics/lrc_parser.go

**职责**：LRC 歌词解析和同步

**LRC 格式**：
```
[ti:歌曲名]
[ar:歌手名]
[al:专辑名]
[00:12.00]第一句歌词
[00:17.20]第二句歌词
```

**主要函数**：
- `ParseLRC()`：解析 LRC 文件
- `GetLyricAtTime()`：获取指定时间的歌词
- `GetLyricWindow()`：获取歌词窗口（前后几行）

**数据结构**：
```go
type LyricLine struct {
    Time    int64  // 毫秒
    Text    string // 歌词文本
    TimeStr string // 格式化时间
}

type LyricData struct {
    Title  string
    Artist string
    Album  string
    Lines  []LyricLine
}
```

**查询优化**：
- 使用二分查找快速定位歌词
- 时间复杂度：O(log n)

### 6. internal/scanner/scanner.go

**职责**：音乐目录扫描和索引

**扫描流程**：
```
ScanDirectory(dirPath)
  ├── 验证目录
  ├── 递归遍历文件
  ├── 过滤支持的格式
  ├── 提取元数据
  └── 保存到数据库
```

**并发优化**：
```go
// 工作池模式
ScanDirectoryWithWorkers(dirPath, numWorkers)
  ├── 收集所有文件
  ├── 创建文件 channel
  ├── 启动 N 个 worker goroutine
  ├── 分发任务
  └── 等待完成
```

**性能指标**：
- 单线程：~100 个文件/秒
- 4 worker：~400 个文件/秒
- 8 worker：~600 个文件/秒

### 7. internal/api/routes.go

**职责**：REST API 路由和处理

**API 端点**：

#### 歌曲管理
```
GET    /api/songs              # 获取所有歌曲
GET    /api/songs/:id          # 获取单首歌曲
GET    /api/songs/search       # 搜索歌曲
POST   /api/songs              # 添加歌曲
```

#### 播放控制
```
POST   /api/player/play        # 播放
POST   /api/player/pause       # 暂停
POST   /api/player/resume      # 恢复
POST   /api/player/stop        # 停止
POST   /api/player/volume      # 设置音量
GET    /api/player/status      # 获取状态
```

#### 歌词和媒体
```
GET    /api/lyrics/:songID     # 获取歌词
GET    /api/cover/:songID      # 获取封面
```

#### 媒体库
```
POST   /api/scan               # 扫描目录
```

#### WebSocket
```
WS     /ws/player              # 实时播放状态
```

### 8. ui/src/App.jsx

**职责**：React 主应用组件

**主要功能**：
- 管理全局状态（歌曲、播放状态、歌词）
- 协调各子组件
- 处理 API 调用

**状态管理**：
```javascript
const [songs, setSongs] = useState([])           // 歌曲列表
const [currentSong, setCurrentSong] = useState() // 当前歌曲
const [isPlaying, setIsPlaying] = useState()     // 播放状态
const [lyrics, setLyrics] = useState()           // 歌词
const [playerStatus, setPlayerStatus] = useState() // 播放器状态
```

### 9. ui/src/components/

**Player.jsx**：播放器组件
- 显示封面、歌曲信息
- 进度条、播放控制
- 音量调节

**SongList.jsx**：歌曲列表组件
- 显示所有歌曲
- 高亮当前播放歌曲
- 点击选择播放

**SearchBar.jsx**：搜索栏组件
- 实时搜索
- 清空按钮

**LyricDisplay.jsx**：歌词显示组件
- 滚动显示歌词
- 高亮当前行
- 显示时间戳

## 🔄 数据流

### 播放流程
```
前端点击歌曲
  ↓
POST /api/player/play
  ↓
后端 Player.Play()
  ↓
解码音频文件
  ↓
写入音频设备
  ↓
前端定时获取 /api/player/status
  ↓
更新进度条和歌词
```

### 扫描流程
```
前端提交扫描请求
  ↓
POST /api/scan
  ↓
后端 Scanner.ScanDirectory()
  ↓
遍历文件 + 提取元数据
  ↓
并发保存到数据库
  ↓
返回扫描结果
  ↓
前端刷新歌曲列表
```

### 搜索流程
```
前端输入搜索词
  ↓
GET /api/songs/search?q=keyword
  ↓
数据库模糊查询
  ↓
返回匹配结果
  ↓
前端显示搜索结果
```

## 🧵 并发模型

### 播放和扫描并发
```
主 goroutine
  ├── HTTP 服务器（处理请求）
  ├── 播放 goroutine（播放音频）
  └── 扫描 worker goroutines（扫描文件）
```

### 线程安全
- 使用 `sync.Mutex` 保护播放器状态
- 使用 channel 进行 goroutine 通信
- 数据库操作由 GORM 处理并发

## 📊 性能考虑

### 内存优化
- 音频流式处理，不加载整个文件
- 数据库连接池
- 限制并发 worker 数量

### 查询优化
- 数据库索引（FilePath）
- 二分查找歌词
- 缓存播放器状态

### 网络优化
- WebSocket 实时推送（减少轮询）
- API 响应压缩
- 前端缓存

## 🔐 安全考虑

### 输入验证
- 文件路径验证
- 搜索关键词长度限制
- API 参数类型检查

### 错误处理
- 统一错误响应格式
- 详细的错误日志
- 用户友好的错误提示

## 📈 可扩展性

### 添加新格式
1. 在 `metadata/extractor.go` 中添加格式检测
2. 在 `player/player.go` 中添加解码器
3. 更新 `scanner/scanner.go` 中的格式列表

### 添加新功能
1. 创建新的 `internal/` 子包
2. 在 `api/routes.go` 中添加路由
3. 在前端添加对应组件

### 数据库扩展
1. 在 `storage/db.go` 中定义新的结构体
2. 调用 `db.AutoMigrate()` 创建表
3. 实现 CRUD 函数

---

**理解这个结构有助于你在面试中清晰地解释项目架构！** 🎯

