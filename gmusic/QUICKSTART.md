# GMusic 快速开始指南

## [object Object] 分钟快速启动

### 前置要求
- Go 1.21+
- Node.js 16+
- 音乐文件（MP3、FLAC、WAV 等）

### 步骤 1：克隆/创建项目

```bash
cd gmusic
```

### 步骤 2：安装依赖

```bash
# 安装所有依赖
make install-deps

# 或者手动安装
go mod download
cd ui && npm install && cd ..
```

### 步骤 3：启动后端（终端 1）

```bash
make dev
# 或
go run cmd/server/main.go
```

你应该看到：
```
🎵 GMusic 服务器启动在 http://localhost:8080
```

### 步骤 4：启动前端（终端 2）

```bash
make frontend
# 或
cd ui && npm run dev
```

你应该看到：
```
VITE v5.0.0  ready in XXX ms

➜  Local:   http://localhost:5173/
```

### 步骤 5：打开浏览器

访问 `http://localhost:5173`，你应该看到 GMusic 界面！

## 📁 添加音乐文件

### 方式 1：通过 API 扫描目录

```bash
curl -X POST http://localhost:8080/api/scan \
  -H "Content-Type: application/json" \
  -d '{
    "dir_path": "/path/to/your/music",
    "workers": 4
  }'
```

### 方式 2：通过前端界面
1. 在前端界面中找到"扫描"按钮
2. 输入音乐目录路径
3. 点击扫描

## 🎵 测试歌曲

如果没有音乐文件，可以：

1. **下载免费音乐**
   - [Free Music Archive](https://freemusicarchive.org/)
   - [ccMixter](http://ccmixter.org/)

2. **创建测试文件**
   ```bash
   # 使用 ffmpeg 创建测试音频
   ffmpeg -f lavfi -i sine=f=440:d=10 -q:a 9 -acodec libmp3lame test.mp3
   ```

## 🎮 基本操作

### 播放歌曲
1. 在右侧列表中点击歌曲
2. 点击播放按钮 ▶️

### 控制播放
- ⏸️ 暂停
- ⏹️ 停止
- 🔊 调节音量

### 搜索歌曲
在顶部搜索框输入歌曲名、歌手或专辑名

### 查看歌词
如果有对应的 `.lrc` 文件，歌词会自动显示

## 📝 添加歌词

### LRC 文件格式

创建 `song_name.lrc` 文件，放在音乐文件同目录：

```
[ti:歌曲名]
[ar:歌手名]
[al:专辑名]
[00:12.00]第一句歌词
[00:17.20]第二句歌词
[00:22.40]第三句歌词
```

### 时间戳格式
- `[MM:SS.ms]` - MM 分钟，SS 秒，ms 毫秒
- 例如：`[01:23.45]` 表示 1 分 23 秒 45 毫秒

### 获取歌词
- [网易云音乐](https://music.163.com/)
- [QQ 音乐](https://y.qq.com/)
- [LRC 歌词库](http://lrclib.net/)

## 🔧 常见问题

### Q: 启动时报错 "音频设备初始化失败"
**A:** 
- Linux：确保 ALSA 已安装 `sudo apt-get install libasound2-dev`
- macOS：检查系统音量是否静音
- Windows：检查音频驱动

### Q: 扫描后看不到歌曲
**A:**
1. 检查目录路径是否正确
2. 检查文件格式是否支持（MP3、FLAC、WAV、AAC）
3. 查看后端日志是否有错误信息

### Q: 歌词不显示
**A:**
1. 确保 `.lrc` 文件与音乐文件同名同目录
2. 检查 LRC 文件编码是 UTF-8
3. 检查时间戳格式是否正确

### Q: 播放卡顿
**A:**
1. 检查系统资源使用情况
2. 尝试减少并发 worker 数量
3. 关闭其他应用

### Q: 前端无法连接后端
**A:**
1. 确保后端已启动（http://localhost:8080）
2. 检查防火墙设置
3. 检查浏览器控制台错误信息

## [object Object] 测试

### 获取所有歌曲
```bash
curl http://localhost:8080/api/songs | jq
```

### 搜索歌曲
```bash
curl "http://localhost:8080/api/songs/search?q=周杰伦" | jq
```

### 播放歌曲
```bash
curl -X POST http://localhost:8080/api/player/play \
  -H "Content-Type: application/json" \
  -d '{"file_path": "/path/to/song.mp3"}'
```

### 获取播放器状态
```bash
curl http://localhost:8080/api/player/status | jq
```

### 获取歌词
```bash
curl http://localhost:8080/api/lyrics/1 | jq
```

## 🐳 Docker 运行

### 构建镜像
```bash
docker-compose build
```

### 启动服务
```bash
docker-compose up
```

### 访问应用
- 后端：http://localhost:8080
- 前端：http://localhost:5173

## 📦 构建发布版本

### 构建后端二进制
```bash
make build
# 输出：bin/gmusic
```

### 构建前端
```bash
cd ui && npm run build
# 输出：ui/dist/
```

## 🧹 清理

### 删除构建文件
```bash
make clean
```

### 重置数据库
```bash
rm gmusic.db
```

### 清理所有缓存
```bash
rm -rf node_modules ui/node_modules
go clean -cache
```

## 📚 下一步

1. **阅读完整文档**：`README.md`
2. **面试准备**：`INTERVIEW_GUIDE.md`
3. **浏览代码**：`internal/` 目录
4. **添加功能**：参考 README 中的扩展功能部分

## 💡 开发提示

### 修改后端代码
```bash
# 后端会自动重新加载（如果使用 air）
# 或手动重启：Ctrl+C，然后 make dev
```

### 修改前端代码
```bash
# 前端会自动热更新
# 只需保存文件即可
```

### 调试
```bash
# 后端调试
dlv debug cmd/server/main.go

# 前端调试
# 在浏览器 F12 打开开发者工具
```

## 🎓 学习路径

1. **第 1 天**：运行项目，熟悉界面
2. **第 2 天**：阅读代码，理解架构
3. **第 3 天**：修改代码，添加功能
4. **第 4 天**：优化性能，完善文档
5. **第 5 天**：准备面试，讲解项目

---

**遇到问题？** 查看 `README.md` 中的故障排除部分或查看后端日志。

**祝你使用愉快！** 🎵

