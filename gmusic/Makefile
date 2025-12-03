.PHONY: help build run dev clean test install-deps

help:
	@echo "GMusic - Golang 本地音乐播放器"
	@echo ""
	@echo "可用命令:"
	@echo "  make install-deps    - 安装所有依赖"
	@echo "  make build          - 构建后端"
	@echo "  make run            - 运行后端服务器"
	@echo "  make dev            - 开发模式运行"
	@echo "  make frontend       - 启动前端开发服务器"
	@echo "  make clean          - 清理构建文件"
	@echo "  make test           - 运行测试"
	@echo "  make lint           - 代码检查"

install-deps:
	@echo "正在安装 Go 依赖..."
	go mod download
	go mod tidy
	@echo "正在安装前端依赖..."
	cd ui && npm install

build:
	@echo "正在构建后端..."
	go build -o bin/gmusic cmd/server/main.go
	@echo "构建完成: bin/gmusic"

run: build
	@echo "启动 GMusic 服务器..."
	./bin/gmusic

dev:
	@echo "开发模式启动..."
	go run cmd/server/main.go

frontend:
	@echo "启动前端开发服务器..."
	cd ui && npm run dev

clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -f gmusic.db
	rm -rf .covers/
	@echo "清理完成"

test:
	@echo "运行测试..."
	go test ./...

lint:
	@echo "代码检查..."
	go vet ./...
	golint ./...

fmt:
	@echo "格式化代码..."
	go fmt ./...
	cd ui && npm run format 2>/dev/null || true

# 一键启动（需要两个终端）
all: install-deps
	@echo "GMusic 已准备就绪！"
	@echo ""
	@echo "请在两个终端分别运行:"
	@echo "  终端 1: make dev"
	@echo "  终端 2: make frontend"

