package main

import (
	"github.com/gofiber/fiber/v3"
	"log"
	"study/internal/api/router"
	"study/internal/di"
	"study/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config files")
	}

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
