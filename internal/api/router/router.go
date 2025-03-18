package router

import (
	"github.com/gofiber/fiber/v3"
	"study/internal/api/handler"
	"study/internal/api/middleware"
)

// SetupRouter 配置路由
func SetupRouter(app *fiber.App, userHandler *handler.UserHandler) {

	app.Post("/signup", userHandler.CreateUser)
	app.Post("/login", userHandler.Login)

	user := app.Group("/user").Use(middleware.AuthMiddleware(userHandler.Handler.Svc.TokenMaker))
	{
		user.Post("/update", userHandler.Update)
	}

	// 可以添加更多路由组，例如 admin
	//admin := app.Group("/admin").Use(middleware.AuthMiddleware(svc.TokenMaker)) // 中间件示例
	//{
	//	admin.Get("/dashboard", nil)
	//}
}
