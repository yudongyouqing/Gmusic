package main

import (
	"fmt"
	"log"

	"github.com/yudongyouqing/GMusic/internal/api"
	"github.com/yudongyouqing/GMusic/internal/storage"
)

func main() {
	// 初始化数据库
	db, err := storage.InitDB("gmusic.db")
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化 API 服务器
	router := api.SetupRouter(db)

	fmt.Println("🎵 GMusic 服务器启动在 http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

