package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func Cors(allowOrigin []string) fiber.Handler {
	return cors.New(cors.Config{
		Next:                nil,
		AllowOriginsFunc:    nil,
		AllowOrigins:        allowOrigin,
		AllowMethods:        nil,
		AllowHeaders:        []string{"Origin, Content-Type, Accept, Authorization"},
		ExposeHeaders:       nil,
		MaxAge:              3600,
		AllowCredentials:    true,
		AllowPrivateNetwork: false,
	})
}
