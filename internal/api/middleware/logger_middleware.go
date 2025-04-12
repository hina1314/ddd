package middleware

import (
	"github.com/gofiber/fiber/v3"
	"log"
	"time"
)

func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		log.Printf("%s %s %v in %v", c.Method(), c.Path(), c.Response().StatusCode(), time.Since(start))
		return err
	}
}
