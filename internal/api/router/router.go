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
	app.Use(middleware.Cors(deps.Config.AllowedOrigins))
	app.Use(middleware.Logger())
	app.Use(middleware.Locale(deps.Config.DefaultLocale)) // 使用配置文件中的语言偏好
	// API 版本分组
	v1 := app.Group("/v1")
	v1.Post("/signup", deps.UserHandler.CreateUser)
	v1.Post("/login", deps.UserHandler.Login)

	user := v1.Group("user", middleware.Auth(deps.ResponseHandler, deps.TokenMaker)) // 仅作用于需要认证的接口
	// 用户路由
	userRoutes(user, deps.UserHandler)

}

// userRoutes 配置用户相关的路由。
func userRoutes(user fiber.Router, h *handler.UserHandler) {
	user.Post("/info", h.Info)
}
