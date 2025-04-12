package handler

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"net/http"
	er "study/util/errors"
)

// BaseHandler 基础处理器，定义通用方法
type BaseHandler struct {
	*er.ErrorHandler
}

func NewBaseHandler(errHandler *er.ErrorHandler) *BaseHandler {
	return &BaseHandler{
		errHandler,
	}
}

func (h *BaseHandler) handleError(ctx fiber.Ctx, err error) error {
	var statusCode int
	var domainErr *er.DomainError

	// 判断是否为自定义领域错误
	if errors.As(err, &domainErr) {
		switch domainErr.Code {
		case er.ErrUserAlreadyExists:
			statusCode = http.StatusConflict
		case er.ErrInvalidInput:
			statusCode = http.StatusBadRequest
		case er.ErrTxError, er.ErrDatabaseError:
			statusCode = http.StatusInternalServerError
		default:
			statusCode = http.StatusInternalServerError
		}
	} else {
		// 非预期错误，统一返回 500
		statusCode = http.StatusInternalServerError
	}

	// 获取本地化信息和调试追踪
	message, debugTrace := h.Handle(ctx.Context(), err)

	// 构建错误响应
	errorBody := map[string]interface{}{
		"message": message,
	}

	if domainErr != nil {
		errorBody["code"] = domainErr.Code
	} else {
		errorBody["code"] = "INTERNAL_ERROR"
	}

	// 添加调试信息（仅调试模式）
	response := map[string]interface{}{
		"error": errorBody,
	}
	if debugTrace != nil {
		response["debug"] = debugTrace
	}

	return ctx.Status(statusCode).JSON(response)
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
