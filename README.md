# GMusic — Golang 本地音乐播放器（Go + Gin + SQLite + Oto + Vue 3）

GMusic 是一套面向实习/校招展示的本地音乐播放器项目：后端使用 Golang + Gin，实现本地媒体扫描、元数据与封面提取、歌词解析、播放控制（Oto 播放）；前端使用 Vue 3 + Vite + Pinia + Vue Router，提供歌曲列表、搜索、播放控制、歌词同步展示的界面。

---

## 功能特性
- 音频播放
  - MP3 解码（go-mp3）
  - 播放 / 暂停 / 恢复 / 停止
  - 进度与时长（近似，按已写入字节推算）
  - 音量控制（对 16-bit PCM 缩放）
- 媒体元数据
  - 歌名、歌手、专辑、年份、Track（dhowden/tag）
  - 专辑封面提取（保存至同目录 .covers/）
- 歌词
  - 读取同名 .lrc
  - LRC 解析、时间轴同步、窗口滚动显示
- 媒体库
  - 目录扫描（并发工作池）
  - SQLite + GORM 存储
  - 搜索（歌名、歌手、专辑）
- API 与前端
  - REST API（Gin）
  - WebSocket 播放状态（备用）
  - 前端 Vue 3 + Vite + Pinia + Router

---

## 技术栈
- 后端
  - Web 框架：Gin
  - 数据库：SQLite + GORM
  - 音频：Oto(v1) 输出、go-mp3 解码（后续可扩展 FLAC）
  - 元数据：dhowden/tag
  - 实时通信：gorilla/websocket（可选）
  - CORS：gin-contrib/cors
- 前端
  - Vue 3、Vite、Vue Router、Pinia
  - Axios（统一 HTTP 客户端）

---

## 目录结构（关键部分）
```
gmusic/
├── cmd/server/main.go          # 服务器入口
├── internal/
│   ├── api/routes.go           # REST 路由 & 控制器（含 CORS、WS）
│   ├── lyrics/lrc_parser.go    # LRC 解析
│   ├── metadata/extractor.go   # 元数据与封面提取、读取 .lrc
│   ├── player/player.go        # 播放引擎（Oto v1 + go-mp3）
│   ├── scanner/scanner.go      # 目录扫描、并发导入
│   └── storage/db.go           # SQLite 模型与 DAO
├── ui/                         # 前端（Vue 3 + Vite）
│   ├── public/
│   └── src/
│       ├── api/music.js        # API 封装
│       ├── service/http.js     # Axios 实例
│       ├── router/index.js     # 路由
│       ├── stores/player.js    # Pinia 播放状态
│       ├── views/              # Library / NowPlaying 页面
│       ├── components/         # Player / SongList / SearchBar / LyricDisplay
│       ├── App.vue
│       └── main.js
├── build.ps1 / build.bat       # Windows 脚本（install-deps/dev/frontend）
├── Makefile                    # *nix 环境常用命令
├── Dockerfile / docker-compose.yml
├── go.mod
└── ... 其他文档
```

---

## 快速开始（Windows）
1. 安装依赖
- 后端与前端依赖
```
cd gmusic
# 如遇到 PowerShell 执行策略限制，先运行：
# Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
.\build.ps1 install-deps
```

2. 启动后端（终端 1）
```
.\build.ps1 dev
# 输出：🎵 GMusic 服务器启动在 http://localhost:8080
```

3. 启动前端（终端 2）
```
.\build.ps1 frontend
# 打开 http://localhost:5173
```

4. 扫描媒体目录（可选）
```
# 使用 curl 示例（替换你的音乐目录）
curl -X POST http://localhost:8080/api/scan \
  -H "Content-Type: application/json" \
  -d '{"dir_path":"D:/Music","workers":4}'
```

---

## 快速开始（通用）
- 后端
```
cd gmusic
# 推荐设置 Go 代理
# go env -w GOPROXY=https://goproxy.cn,direct

go mod tidy

go run cmd/server/main.go
```
- 前端
```
cd gmusic/ui
npm install
npm run dev
```

---

## API 速查
- 歌曲
  - GET /api/songs
  - GET /api/songs/:id
  - GET /api/songs/search?q=keyword
  - POST /api/songs { file_path }
- 播放控制
  - POST /api/player/play { file_path }
  - POST /api/player/pause
  - POST /api/player/resume
  - POST /api/player/stop
  - POST /api/player/volume { volume: 0..1 }
  - GET  /api/player/status
- 歌词与封面
  - GET /api/lyrics/:songID
  - GET /api/cover/:songID
- 扫描
  - POST /api/scan { dir_path, workers }
- WebSocket（可选）
  - GET /ws/player

说明：后端为纯 API，直接访问 http://localhost:8080 会返回 404；前端开发服务器在 http://localhost:5173。

---

## 常见问题（FAQ）
- 访问 8080 返回 404？
  - 正常。后端仅提供 API，请访问前端 http://localhost:5173 或直接调用 /api/*
- go.sum 缺失 / missing go.sum entry？
  - 运行 `.\build.ps1 install-deps` 或 `go mod tidy && go mod download`
- Oto API 编译报错（NewContext 参数不匹配 / Player.Write 未定义）？
  - 当前使用 Oto v1（go.mod: github.com/hajimehoshi/oto v0.7.x），player.go 已按 v1 API 实现
- 跨域（CORS）失败？
  - routes.go 已启用 `router.Use(cors.Default())`；重启后端生效
- 歌词不显示？
  - 确认与音频同目录且同名的 .lrc 存在，编码建议 UTF-8
- 时长显示不准确？
  - 目前按字节推算时长（MP3），后续可通过帧级解析或引入更精确的解码器计算

---

## 推送到 GitHub（重要）
当前 module 名称默认为 `github.com/yourusername/gmusic`。若你要开源到 GitHub，请：
1. 修改 go.mod 顶部：
```
module github.com/<你的GitHub用户名>/gmusic
```
2. 全局替换 import 前缀（将所有 `github.com/yourusername/gmusic` 替换为你的真实路径）
3. 整理依赖并验证：`go mod tidy`、`go build ./...` 或 `./build.ps1 dev`
4. Git 推送：
```
git init
git add .
git commit -m "feat: init GMusic (Go+Vue)"
git branch -M main
git remote add origin https://github.com/<你的GitHub用户名>/gmusic.git
git push -u origin main
```

---

## Roadmap（可选优化）
- 播放格式扩展：FLAC/WAV/AAC 播放
- 精确时长与进度：基于帧或解码器时间戳
- 进度拖动 / 快进快退
- 播放模式：顺序/随机/单曲循环
- 播放列表与收藏/喜欢
- 歌词逐字高亮与网络歌词搜索
- 桌面端封装（Wails）
- Docker 化前端（Nginx 托管 dist）

---

## 许可证
MIT

---

## 致谢
- Oto: 跨平台音频输出
- go-mp3 / mewkiz/flac: 音频解码
- dhowden/tag: 元数据读取
- Gin / GORM / SQLite / Vue / Vite / Pinia / Vue Router

