package router

import (
	"github.com/gofiber/fiber/v3"
	"study/internal/api/handler"
	"study/internal/api/middleware"
	"study/internal/di"
)

// Setup 配置应用程序的所有路由。
func Setup(app *fiber.App, deps *di.Dependencies) {
	// 全局中间件
	app.Use(middleware.Logger())
	app.Use(middleware.Locale(deps.Config.DefaultLocale)) // 使用配置文件中的语言偏好

	// API 版本分组
	v1 := app.Group("/v1")

	// 用户路由
	userRoutes(v1, deps.UserHandler)
}

// userRoutes 配置用户相关的路由。
func userRoutes(group fiber.Router, h *handler.UserHandler) {
	users := group.Group("/user")
	users.Post("/signup", h.CreateUser)
	// users.Post("/login", h.Login) // 未来启用
}
