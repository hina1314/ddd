package handler

import (
	"github.com/gofiber/fiber/v3"
)

// BaseHandler 基础处理器，定义通用方法
type BaseHandler struct{}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

// ErrorResponse 统一错误响应
func (h *BaseHandler) ErrorResponse(c fiber.Ctx, err error) error {
	switch err.Error() {
	case "unique_violation":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "手机号或用户名已存在"})
	case "invalid_phone":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "手机号格式不正确"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
}

// 自定义手机号验证器
//var phoneValidator validator.Func = func(fl validator.FieldLevel) bool {
//	phone, ok := fl.Field().Interface().(string)
//	if !ok {
//		return false
//	}
//	// 中国手机号格式：11 位，以 1 开头，第二位是 3-9
//	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
//	return re.MatchString(phone)
//}
