# Windows 用户必读

## 🎯 你的问题

```
make : The term 'make' is not recognized...
```

**原因**：Windows 不支持 `make` 命令

**解决**：使用 PowerShell 或 Batch 脚本

## ✅ 立即修复（3 步）

### 步骤 1：打开 PowerShell

右键点击开始菜单 → 选择 "Windows PowerShell (管理员)"

### 步骤 2：进入项目目录

```powershell
cd D:\GMusic\gmusic
```

### 步骤 3：运行启动脚本

```powershell
# 如果遇到执行策略错误，先运行这个
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# 安装依赖
.\build.ps1 install-deps

# 启动后端（在这个窗口运行）
.\build.ps1 dev
```

### 步骤 4：打开另一个 PowerShell 窗口

```powershell
cd D:\GMusic\gmusic

# 启动前端
.\build.ps1 frontend
```

### 步骤 5：打开浏览器

访问 `http://localhost:5173`

✅ **完成！** 你应该看到播放器界面了

## 📝 可用命令

```powershell
# 安装依赖
.\build.ps1 install-deps

# 启动后端
.\build.ps1 dev

# 启动前端
.\build.ps1 frontend

# 构建项目
.\build.ps1 build

# 清理文件
.\build.ps1 clean

# 显示帮助
.\build.ps1 help
```

## 🆘 常见问题

### Q: 执行策略错误

**错误**：
```
build.ps1 cannot be loaded because running scripts is disabled on this system.
```

**解决**：
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Q: Go 或 npm 命令不存在

**解决**：
1. 确保已安装 Go 和 Node.js
2. 重启 PowerShell
3. 检查 PATH 环境变量

### Q: 端口被占用

**解决**：
```powershell
# 查找占用 8080 端口的进程
netstat -ano | findstr :8080

# 杀死进程
taskkill /PID <PID> /F
```

## 🔄 替代方案

### 使用 Batch 脚本

```cmd
cd D:\GMusic\gmusic
build.bat install-deps
build.bat dev
```

### 手动启动

```powershell
cd D:\GMusic\gmusic

# 安装依赖
go mod download
cd ui && npm install && cd ..

# 启动后端
go run cmd/server/main.go

# 启动前端（新窗口）
cd ui && npm run dev
```

### 使用 Docker

```powershell
cd D:\GMusic\gmusic
docker-compose up
```

## 📚 更多帮助

- **完整 Windows 指南**：WINDOWS_SETUP.md
- **快速开始**：QUICKSTART.md
- **项目文档**：README.md

## 🎉 成功标志

### 后端
```
🎵 GMusic 服务器启动在 http://localhost:8080
```

### 前端
```
VITE v5.0.0  ready in XXX ms
➜  Local:   http://localhost:5173/
```

### 浏览器
看到漂亮的播放器界面

---

**现在就试试吧！** 🚀

```powershell
cd D:\GMusic\gmusic
.\build.ps1 install-deps
.\build.ps1 dev
```

