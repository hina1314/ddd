package main

import (
	"github.com/gofiber/fiber/v3"
	"log"
	"study/internal/api/router"
	"study/internal/di"
	"study/util"
	"study/util/i18n"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config files")
	}

	translator := i18n.NewFileTranslator("en")
	// 加载翻译文件
	if err = translator.LoadTranslations("./config/i18n"); err != nil {
		log.Fatal(err)
	}

	// 初始化全局翻译服务
	i18n.InitTranslator(translator, "zh")

	deps, err := di.InitializeDependencies(config)
	if err != nil {
		log.Fatal("cannot initialize dependencies:", err)
	}

	server := fiber.New()
	router.SetupRouter(server, deps.UserHandler)

	err = server.Listen(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
