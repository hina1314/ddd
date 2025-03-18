package main

import (
	"github.com/gofiber/fiber/v3"
	"log"
	"study/internal/api/handler"
	"study/internal/api/router"
	"study/internal/app"
	"study/internal/infra/repository"
	"study/internal/svc"
	"study/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config files")
	}

	// 依赖注入
	svcCtx := svc.NewServiceContext(config)

	userRepo := repository.NewUserRepository(svcCtx)
	userAccountRepo := repository.NewUserAccountRepository(svcCtx)

	userService := app.NewUserService(userRepo, userAccountRepo)

	h := handler.NewHandler(svcCtx)
	userHandler := handler.NewUserHandler(h, userService)
	server := fiber.New()

	// 配置路由
	router.SetupRouter(server, userHandler)

	err = server.Listen(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
