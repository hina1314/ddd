package router

import (
	"github.com/gofiber/fiber/v3"
	"study/internal/api/handler"
	"study/internal/api/middleware"
)

// SetupRouter 配置所有路由
func SetupRouter(app *fiber.App, userHandler *handler.UserHandler) {
	// 全局中间件（如日志）
	app.Use(middleware.Logger())

	// API 版本分组
	v1 := app.Group("/v1")

	// 用户相关路由
	setupUserRoutes(v1, userHandler)
}

// setupUserRoutes 用户模块路由
func setupUserRoutes(group fiber.Router, h *handler.UserHandler) {
	users := group.Group("/user")
	users.Post("/signup", h.CreateUser) // 创建用户
}
