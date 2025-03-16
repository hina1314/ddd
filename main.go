package main

import (
	"github.com/gofiber/fiber/v3"
	"log"
	"study/internal/router"
	"study/internal/svc"
	"study/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config files")
	}

	svcCtx := svc.NewServiceContext(config)
	app := fiber.New()

	// 配置路由
	router.SetupRouter(app, svcCtx)

	err = app.Listen(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
