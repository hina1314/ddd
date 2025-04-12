package handler

import "github.com/gofiber/fiber/v3"

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

//type successResponse struct {
//	Code    string
//	Message string
//	Data    interface{}
//}

//	func errorResponse(err error) fiber.Map {
//		return fiber.Map{"error": err.Error()}
//	}
func successResponse(c fiber.Ctx, msg string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": fiber.StatusOK, "msg": msg, "data": data})
}
