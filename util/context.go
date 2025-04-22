package util

import (
	"github.com/gofiber/fiber/v3"
	"study/internal/api/middleware"
	"study/token"
)

func GetAuthPayload(c fiber.Ctx) *token.Payload {
	payload, ok := c.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)
	if !ok {
		return nil // 或者 panic/log
	}
	return payload
}
