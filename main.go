package main

import (
	"log"
	"study/config"
	"study/internal/api/router"
	"study/internal/di"
)

// main 是应用程序的入口点。
func main() {
	// 加载配置
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化依赖
	deps, err := di.NewDependencies(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// 创建 Fiber 服务器
	server := deps.NewServer()

	// 设置路由
	router.Setup(server, deps)

	// 启动服务器
	if err := server.Listen(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
