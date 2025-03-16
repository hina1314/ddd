package handler

import (
	"github.com/gofiber/fiber/v3"
)

func errorResponse(err error) fiber.Map {
	return fiber.Map{"error": err.Error()}
}

func successResponse(msg string) fiber.Map {
	return fiber.Map{"msg": msg}
}
