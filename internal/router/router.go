package router

import (
	"study/internal/handler"
	"study/internal/middleware"
	"study/internal/svc"

	"github.com/gofiber/fiber/v3"
)

// SetupRouter 配置路由
func SetupRouter(app *fiber.App, svc *svc.ServiceContext) {
	// 示例：用户相关的路由组
	h := handler.NewHandler(svc)
	app.Post("/signup", h.CreateUser)
	app.Post("/login", h.Login)

	user := app.Group("/user").Use(middleware.AuthMiddleware(svc.TokenMaker))
	{
		user.Post("/update", h.Update)
	}

	// 可以添加更多路由组，例如 admin
	//admin := app.Group("/admin").Use(middleware.AuthMiddleware(svc.TokenMaker)) // 中间件示例
	//{
	//	admin.Get("/dashboard", nil)
	//}
}
