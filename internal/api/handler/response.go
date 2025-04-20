package handler

import (
	"github.com/gofiber/fiber/v3"
	"study/util/errors"
)

// SuccessResponse 定义成功响应的结构。
type successResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (h *BaseHandler) successResponse(c fiber.Ctx, msg string, data interface{}) error {
	message := h.TranslationService.T(c.Context(), msg, nil)
	response := successResponse{
		Code: fiber.StatusOK,
		Msg:  message,
		Data: data,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

type errorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Debug *errors.ErrorTrace `json:"debug,omitempty"`
}
