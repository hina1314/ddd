package middleware

import (
	"github.com/gofiber/fiber/v3"
	"study/util/i18n"
)

func Locale(defaultLocale string) fiber.Handler {
	return func(c fiber.Ctx) error {
		locale := c.Get("Accept-Language", defaultLocale)
		ctx := i18n.WithLocale(c.Context(), locale)
		c.SetContext(ctx)
		return c.Next()
	}
}
