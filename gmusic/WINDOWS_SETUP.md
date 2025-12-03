# Windows 系统快速启动指南

如果你在 Windows 上使用，这份指南会帮助你快速启动项目。

## ⚠️ 问题排查

### 问题：`make` 命令不存在

**原因**：Windows 原生不支持 `make` 命令

**解决方案**：使用提供的 PowerShell 或 Batch 脚本

## 🚀 方案 1：使用 PowerShell 脚本（推荐）

### 前置要求
- PowerShell 5.0+（Windows 10/11 默认自带）
- Go 1.21+
- Node.js 16+

### 启动步骤

**步骤 1**：打开 PowerShell（以管理员身份）

```powershell
# 如果遇到执行策略错误，运行此命令
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

**步骤 2**：进入项目目录

```powershell
cd D:\GMusic\gmusic
```

**步骤 3**：安装依赖

```powershell
.\build.ps1 install-deps
```

**步骤 4**：启动后端（PowerShell 窗口 1）

```powershell
.\build.ps1 dev
```

你应该看到：
```
🎵 GMusic 服务器启动在 http://localhost:8080
```

**步骤 5**：启动前端（PowerShell 窗口 2）

```powershell
.\build.ps1 frontend
```

你应该看到：
```
VITE v5.0.0  ready in XXX ms
➜  Local:   http://localhost:5173/
```

**步骤 6**：打开浏览器

访问 `http://localhost:5173`

✅ 完成！

## 🚀 方案 2：使用 Batch 脚本

### 启动步骤

**步骤 1**：打开 CMD（命令提示符）

**步骤 2**：进入项目目录

```cmd
cd D:\GMusic\gmusic
```

**步骤 3**：安装依赖

```cmd
build.bat install-deps
```

**步骤 4**：启动后端（CMD 窗口 1）

```cmd
build.bat dev
```

**步骤 5**：启动前端（CMD 窗口 2）

```cmd
build.bat frontend
```

**步骤 6**：打开浏览器

访问 `http://localhost:5173`

## 🚀 方案 3：手动启动（不使用脚本）

### 安装依赖

```powershell
# Go 依赖
go mod download
go mod tidy

# 前端依赖
cd ui
npm install
cd ..
```

### 启动后端（PowerShell 窗口 1）

```powershell
go run cmd/server/main.go
```

### 启动前端（PowerShell 窗口 2）

```powershell
cd ui
npm run dev
cd ..
```

### 打开浏览器

访问 `http://localhost:5173`

## 📝 PowerShell 脚本命令

```powershell
# 安装依赖
.\build.ps1 install-deps

# 开发模式运行后端
.\build.ps1 dev

# 启动前端开发服务器
.\build.ps1 frontend

# 构建后端
.\build.ps1 build

# 清理构建文件
.\build.ps1 clean

# 显示帮助
.\build.ps1 help
```

## 📝 Batch 脚本命令

```cmd
REM 安装依赖
build.bat install-deps

REM 开发模式运行后端
build.bat dev

REM 启动前端开发服务器
build.bat frontend

REM 构建后端
build.bat build

REM 清理构建文件
build.bat clean

REM 显示帮助
build.bat help
```

## ❓ 常见问题

### Q: PowerShell 执行策略错误

**错误信息**：
```
build.ps1 cannot be loaded because running scripts is disabled on this system.
```

**解决方案**：
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

然后选择 `Y` 确认。

### Q: Go 命令不存在

**解决方案**：
1. 确保 Go 已安装
2. 检查 Go 是否在 PATH 中
3. 重启 PowerShell/CMD

### Q: npm 命令不存在

**解决方案**：
1. 确保 Node.js 已安装
2. 检查 npm 是否在 PATH 中
3. 重启 PowerShell/CMD

### Q: 端口被占用

**错误信息**：
```
listen tcp :8080: bind: An attempt was made to use a socket in a way forbidden by its access rules.
```

**解决方案**：
```powershell
# 查找占用 8080 端口的进程
netstat -ano | findstr :8080

# 杀死进程（替换 PID）
taskkill /PID <PID> /F
```

### Q: 前端无法连接后端

**解决方案**：
1. 确保后端已启动（http://localhost:8080）
2. 检查防火墙设置
3. 查看浏览器控制台错误信息

## 🐳 Docker 方案（无需 Go 和 Node.js）

如果你不想安装 Go 和 Node.js，可以使用 Docker：

### 前置要求
- Docker Desktop（Windows 版本）

### 启动步骤

```powershell
# 进入项目目录
cd D:\GMusic\gmusic

# 构建镜像
docker-compose build

# 启动服务
docker-compose up
```

### 访问应用
- 后端：http://localhost:8080
- 前端：http://localhost:5173

## 📂 项目结构

```
gmusic/
├── build.ps1              # PowerShell 脚本
├── build.bat              # Batch 脚本
├── cmd/
├── internal/
├── ui/
├── go.mod
├── Makefile               # Linux/macOS 用
├── Dockerfile
├── docker-compose.yml
└── ...
```

## 💡 推荐工作流

### 开发流程

1. **打开 PowerShell 窗口 1**（后端）
   ```powershell
   cd D:\GMusic\gmusic
   .\build.ps1 dev
   ```

2. **打开 PowerShell 窗口 2**（前端）
   ```powershell
   cd D:\GMusic\gmusic
   .\build.ps1 frontend
   ```

3. **打开浏览器**
   访问 `http://localhost:5173`

4. **修改代码**
   - 后端代码改动后，按 Ctrl+C 重启后端
   - 前端代码改动后，会自动热更新

## 🔧 环境变量配置

如果需要自定义配置，可以创建 `.env` 文件：

```
MUSIC_DIR=/path/to/music
DB_PATH=gmusic.db
PORT=8080
```

## 📚 更多帮助

- **快速开始**：查看 QUICKSTART.md
- **完整文档**：查看 README.md
- **面试准备**：查看 INTERVIEW_GUIDE.md
- **项目结构**：查看 PROJECT_STRUCTURE.md

## ✅ 检查清单

- [ ] 已安装 Go 1.21+
- [ ] 已安装 Node.js 16+
- [ ] 已进入项目目录
- [ ] 已运行 `.\build.ps1 install-deps`
- [ ] 后端已启动（http://localhost:8080）
- [ ] 前端已启动（http://localhost:5173）
- [ ] 浏览器可以访问前端
- [ ] 可以看到播放器界面

## 🎉 成功标志

当你看到以下内容时，说明项目启动成功：

**后端**：
```
🎵 GMusic 服务器启动在 http://localhost:8080
```

**前端**：
```
VITE v5.0.0  ready in XXX ms
➜  Local:   http://localhost:5173/
```

**浏览器**：
看到漂亮的播放器界面，包含搜索栏、歌曲列表、播放器等

---

**祝你在 Windows 上使用愉快！** 🎵✨

有任何问题，查看相应的文档或重新阅读本指南。

