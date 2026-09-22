package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/google/uuid"

	"study/internal/api/handler"
	"study/internal/api/middleware"
	"study/internal/di"
)

func SetupMiddleware(app *fiber.App, deps *di.Dependencies) {
	// 必须最先生成 request_id，后面的日志才能读取它。
	app.Use(requestid.New(requestid.Config{
		Generator: uuid.NewString,
	}))

	app.Use(middleware.Logger())
	app.Use(middleware.Metrics())
	app.Get("/metrics", middleware.MetricsHandler())

	app.Use(middleware.Cors(deps.Config.AllowedOrigins))
	app.Use(middleware.Locale(deps.Config.DefaultLocale))
}

// Setup 只负责注册业务路由。
func Setup(app *fiber.App, deps *di.Dependencies) {
	v1 := app.Group("/v1")
	v1.Post("/signup", deps.UserHandler.CreateUser)
	v1.Post("/login", deps.UserHandler.Login)

	user := v1.Group("user",
		middleware.Auth(deps.ResponseHandler, deps.TokenMaker),
	)
	order := v1.Group("order",
		middleware.Auth(deps.ResponseHandler, deps.TokenMaker),
	)

	userRoutes(user, deps.UserHandler)
	orderRoutes(order, deps.OrderHandler)
}

// userRoutes 配置用户相关的路由。
func userRoutes(user fiber.Router, h *handler.UserHandler) {
	user.Post("/info", h.Info)
	user.Post("/update", h.Update)
}

func orderRoutes(order fiber.Router, h *handler.OrderHandler) {
	order.Post("/create", h.CreateOrder)
}
