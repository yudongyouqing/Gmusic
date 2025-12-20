# GMusic — Golang 本地音乐播放器（Go + Gin + SQLite + Oto + Vue 3）

GMusic 是一套面向实习/校招展示的本地音乐播放器项目：后端使用 Golang + Gin，实现本地媒体扫描、元数据与封面提取、歌词解析、播放控制（Oto 播放）；前端使用 Vue 3 + Vite + Pinia + Vue Router，提供歌曲列表、搜索、播放控制、歌词同步展示的界面。

---

## 功能特性
- **音频播放**
  - MP3 解码（go-mp3）
  - 播放 / 暂停 / 恢复 / 停止
  - 进度与时长（精确，基于 MP3 帧解析）
  - 音量控制
  - 播放模式：列表循环、随机播放、单曲循环
- **媒体元数据**
  - 歌名、歌手、专辑、年份、Track（dhowden/tag）
  - 专辑封面提取（保存至同目录 .covers/）
- **歌词**
  - 读取同名 .lrc
  - LRC 解析、时间轴同步、窗口滚动显示
  - 播放页歌词设置：字体大小、粗细、模糊非当前行、显示/隐藏翻译
  - 播放页背景模糊度调节
- **媒体库**
  - 目录扫描（并发工作池）
  - SQLite + GORM 存储
  - 搜索（歌名、歌手、专辑）
  - 手动排序（拖拽）与按标题/歌手/专辑排序
- **API 与前端**
  - REST API（Gin）
  - 前端 Vue 3 + Vite + Pinia + Router
  - 主题设置：毛玻璃/当前风格、透明度、饱和度
  - 播放页自定义背景

---

## 技术栈
- **后端**
  - Web 框架：Gin
  - 数据库：SQLite + GORM
  - 音频：Oto(v1) 输出、go-mp3 解码
  - 元数据：dhowden/tag
- **前端**
  - Vue 3、Vite、Vue Router、Pinia
  - Axios（统一 HTTP 客户端）

---

## 目录结构（关键部分）
```
gmusic/
├── cmd/server/main.go          # 服务器入口
├── internal/
│   ├── api/routes.go           # REST 路由 & 控制器
│   ├── lyrics/lrc_parser.go    # LRC 解析
│   ├── metadata/extractor.go   # 元数据与封面提取
│   ├── player/player.go        # 播放引擎
│   ├── scanner/scanner.go      # 目录扫描
│   └── storage/db.go           # SQLite 模型
├── ui/                         # 前端（Vue 3 + Vite）
│   └── src/
│       ├── api/music.js        # API 封装
│       ├── stores/             # Pinia 状态管理 (player, ui, lyric, settings)
│       ├── views/              # 页面 (Library, NowPlaying, Queue, Settings)
│       └── components/         # 组件 (Player, SongList, LyricDisplay, LyricControls)
├── build.ps1 / build.bat       # Windows 构建脚本
└── go.mod
```

---

## 快速开始（Windows）
1. **安装依赖**
```powershell
cd gmusic
# 如遇 PowerShell 执行策略限制，先运行：
# Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
.\build.ps1 install-deps
```

2. **启动后端**（终端 1）
```powershell
.\build.ps1 dev
# 输出：🎵 GMusic 服务器启动在 http://localhost:8080
```

3. **启动前端**（终端 2）
```powershell
.\build.ps1 frontend
# 打开 http://localhost:5173
```

4. **扫描媒体目录**（可选）
```powershell
# 使用 curl 示例（替换你的音乐目录）
curl -X POST http://localhost:8080/api/scan \
  -H "Content-Type: application/json" \
  -d '{"dir_path":"D:/Music","workers":4}'
```

---

## API 速查
- 歌曲：`GET /api/songs`, `GET /api/songs/:id`, `GET /api/songs/search?q=keyword`
- 播放控制：`POST /api/player/play`, `POST /api/player/pause`, `POST /api/player/resume`, `POST /api/player/stop`, `POST /api/player/volume`, `GET /api/player/status`
- 歌词与封面：`GET /api/lyrics/:songID`, `GET /api/cover/:songID`
- 扫描：`POST /api/scan`

---

## 许可证
MIT
