# Windows 快速修复 - 立即启动

你遇到了 `make` 命令不存在的问题。这很正常，因为 Windows 不支持 `make`。

## ⚡ 30 秒快速修复

### 选项 1：使用 PowerShell 脚本（推荐）

```powershell
# 1. 打开 PowerShell（以管理员身份）
# 2. 进入项目目录
cd D:\GMusic\gmusic

# 3. 如果遇到执行策略错误，运行此命令
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# 4. 安装依赖
.\build.ps1 install-deps

# 5. 启动后端（PowerShell 窗口 1）
.\build.ps1 dev

# 6. 启动前端（PowerShell 窗口 2）
.\build.ps1 frontend

# 7. 打开浏览器访问 http://localhost:5173
```

### 选项 2：使用 Batch 脚本

```cmd
REM 1. 打开 CMD（命令提示符）
REM 2. 进入项目目录
cd D:\GMusic\gmusic

REM 3. 安装依赖
build.bat install-deps

REM 4. 启动后端（CMD 窗口 1）
build.bat dev

REM 5. 启动前端（CMD 窗口 2）
build.bat frontend

REM 6. 打开浏览器访问 http://localhost:5173
```

### 选项 3：手动启动（不使用脚本）

```powershell
# 进入项目目录
cd D:\GMusic\gmusic

# 安装 Go 依赖
go mod download
go mod tidy

# 安装前端依赖
cd ui
npm install
cd ..

# 启动后端（PowerShell 窗口 1）
go run cmd/server/main.go

# 启动前端（PowerShell 窗口 2）
cd ui
npm run dev
cd ..

# 打开浏览器访问 http://localhost:5173
```

## ✅ 成功标志

### 后端启动成功
```
🎵 GMusic 服务器启动在 http://localhost:8080
```

### 前端启动成功
```
VITE v5.0.0  ready in XXX ms
➜  Local:   http://localhost:5173/
```

### 浏览器显示
看到漂亮的播放器界面

## ❓ 遇到问题？

### 问题 1：PowerShell 执行策略错误

**错误**：
```
build.ps1 cannot be loaded because running scripts is disabled on this system.
```

**解决**：
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 问题 2：Go 或 npm 命令不存在

**解决**：
1. 确保已安装 Go 1.21+ 和 Node.js 16+
2. 重启 PowerShell/CMD
3. 检查 PATH 环境变量

### 问题 3：端口被占用

**解决**：
```powershell
# 查找占用 8080 端口的进程
netstat -ano | findstr :8080

# 杀死进程（替换 PID）
taskkill /PID <PID> /F
```

## 📚 完整指南

详细的 Windows 启动指南：查看 **WINDOWS_SETUP.md**

## 🎉 现在就开始！

选择上面的任意一个选项，立即启动项目！

**祝你成功！** 🚀

