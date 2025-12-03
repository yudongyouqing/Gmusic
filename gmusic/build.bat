@echo off
REM GMusic 构建脚本 (Batch)
REM 用法: build.bat install-deps | build.bat dev | build.bat frontend

setlocal enabledelayedexpansion

if "%1"=="" (
    call :show_help
    exit /b 0
)

if /i "%1"=="install-deps" (
    call :install_deps
    exit /b 0
)

if /i "%1"=="dev" (
    call :start_dev
    exit /b 0
)

if /i "%1"=="frontend" (
    call :start_frontend
    exit /b 0
)

if /i "%1"=="build" (
    call :build_backend
    exit /b 0
)

if /i "%1"=="clean" (
    call :clean
    exit /b 0
)

if /i "%1"=="help" (
    call :show_help
    exit /b 0
)

echo 未知命令: %1
call :show_help
exit /b 1

:show_help
echo.
echo GMusic - Golang 本地音乐播放器
echo.
echo 可用命令:
echo   build.bat install-deps    - 安装所有依赖
echo   build.bat dev             - 开发模式运行后端
echo   build.bat frontend        - 启动前端开发服务器
echo   build.bat build           - 构建后端
echo   build.bat clean           - 清理构建文件
echo   build.bat help            - 显示此帮助信息
echo.
exit /b 0

:install_deps
echo 正在安装 Go 依赖...
go mod download
go mod tidy

echo 正在安装前端依赖...
cd ui
call npm install
cd ..

echo ✅ 依赖安装完成！
exit /b 0

:start_dev
echo 🎵 开发模式启动后端...
go run cmd/server/main.go
exit /b 0

:start_frontend
echo 🎨 启动前端开发服务器...
cd ui
call npm run dev
cd ..
exit /b 0

:build_backend
echo 正在构建后端...
if not exist "bin" mkdir bin
go build -o bin/gmusic.exe cmd/server/main.go
echo ✅ 构建完成: bin/gmusic.exe
exit /b 0

:clean
echo 清理构建文件...
if exist "bin" rmdir /s /q bin
if exist "gmusic.db" del gmusic.db
if exist ".covers" rmdir /s /q .covers
echo ✅ 清理完成
exit /b 0

