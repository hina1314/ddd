package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/google/uuid"

	"github.com/hina1314/ddd/internal/api/handler"
	"github.com/hina1314/ddd/internal/api/middleware"
	"github.com/hina1314/ddd/internal/di"
	"github.com/hina1314/ddd/util/errors"
)

func SetupMiddleware(app *fiber.App, deps *di.Dependencies) {
	// 必须最先生成 request_id，后面的日志才能读取它。
	app.Use(requestid.New(requestid.Config{
		Generator: uuid.NewString,
	}))
	app.Use(middleware.Logger())
	// recover 在 Logger 内层运行，确保 panic 恢复后也能输出请求日志。
	app.Use(recover.New())
	app.Use(middleware.Metrics())
	app.Use(helmet.New())
	app.Use(middleware.Cors(deps.Config.AllowedOrigins))
	app.Use(middleware.Locale(deps.Config.DefaultLocale, deps.Config.SupportedLocales))
	app.Get("/metrics", middleware.MetricsAuth(deps.Config.MetricsToken), middleware.MetricsHandler())
}

// Setup 只负责注册业务路由。
func Setup(app *fiber.App, deps *di.Dependencies) {
	v1 := app.Group("/v1")
	authLimiter := limiter.New(limiter.Config{
		Max:        deps.Config.AuthRateLimitMax,
		Expiration: deps.Config.AuthRateLimitWindow,
		LimitReached: func(c fiber.Ctx) error {
			return deps.ResponseHandler.HandleError(
				c,
				errors.New(errors.ErrRateLimited, "too many authentication attempts"),
			)
		},
	})
	v1.Post("/signup", authLimiter, deps.UserHandler.CreateUser)
	v1.Post("/login", authLimiter, deps.UserHandler.Login)

	user := v1.Group("user",
		middleware.Auth(deps.ResponseHandler, deps.TokenMaker),
	)
	userRoutes(user, deps.UserHandler)
}

// userRoutes 配置用户相关的路由。
func userRoutes(user fiber.Router, h *handler.UserHandler) {
	user.Get("/info", h.Info)
	user.Patch("/profile", h.Update)
}
