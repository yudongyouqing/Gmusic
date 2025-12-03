# GMusic 构建脚本 (PowerShell)
# 用法: .\build.ps1 install-deps | .\build.ps1 dev | .\build.ps1 frontend

param(
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "GMusic - Golang 本地音乐播放器" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "可用命令:" -ForegroundColor Green
    Write-Host "  .\build.ps1 install-deps    - 安装所有依赖"
    Write-Host "  .\build.ps1 dev             - 开发模式运行后端"
    Write-Host "  .\build.ps1 frontend        - 启动前端开发服务器"
    Write-Host "  .\build.ps1 build           - 构建后端"
    Write-Host "  .\build.ps1 clean           - 清理构建文件"
    Write-Host "  .\build.ps1 help            - 显示此帮助信息"
    Write-Host ""
}

function Install-Deps {
    Write-Host "正在安装 Go 依赖..." -ForegroundColor Yellow
    go mod download
    go mod tidy
    
    Write-Host "正在安装前端依赖..." -ForegroundColor Yellow
    Set-Location ui
    npm install
    Set-Location ..
    
    Write-Host "✅ 依赖安装完成！" -ForegroundColor Green
}

function Start-Dev {
    Write-Host "🎵 开发模式启动后端..." -ForegroundColor Cyan
    go run cmd/server/main.go
}

function Start-Frontend {
    Write-Host "🎨 启动前端开发服务器..." -ForegroundColor Cyan
    Set-Location ui
    npm run dev
    Set-Location ..
}

function Build-Backend {
    Write-Host "正在构建后端..." -ForegroundColor Yellow
    if (-not (Test-Path "bin")) {
        New-Item -ItemType Directory -Path "bin" | Out-Null
    }
    go build -o bin/gmusic.exe cmd/server/main.go
    Write-Host "✅ 构建完成: bin/gmusic.exe" -ForegroundColor Green
}

function Clean {
    Write-Host "清理构建文件..." -ForegroundColor Yellow
    if (Test-Path "bin") {
        Remove-Item -Recurse -Force "bin"
    }
    if (Test-Path "gmusic.db") {
        Remove-Item "gmusic.db"
    }
    if (Test-Path ".covers") {
        Remove-Item -Recurse -Force ".covers"
    }
    Write-Host "✅ 清理完成" -ForegroundColor Green
}

# 执行命令
switch ($Command.ToLower()) {
    "install-deps" { Install-Deps }
    "dev" { Start-Dev }
    "frontend" { Start-Frontend }
    "build" { Build-Backend }
    "clean" { Clean }
    "help" { Show-Help }
    default { 
        Write-Host "❌ 未知命令: $Command" -ForegroundColor Red
        Show-Help
    }
}

