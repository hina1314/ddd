package handler

import (
	"github.com/gofiber/fiber/v3"
	"study/internal/app"
)

type UserHandler struct {
	userService *app.UserService
}

func NewUserHandler(userService *app.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	name := c.Params("name")
	user, err := h.userService.RegisterUser(c.Context(), name, "", "", "")
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}
