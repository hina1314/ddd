package response

import (
	stdErr "errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"study/util/errors"
	"study/util/i18n"
)

// Response 定义响应的结构。
type Response struct {
	Code  interface{}        `json:"code"` // 成功是 int 0，失败是 string errorCode
	Msg   string             `json:"msg"`  // 已翻译好的消息
	Data  interface{}        `json:"data,omitempty"`
	Debug *errors.ErrorTrace `json:"debug,omitempty"`
}

func (h *ResponseHandler) Success(c fiber.Ctx, msg string, data interface{}) error {
	message := h.TranslationService.T(c.Context(), msg, nil)
	response := Response{
		Code: 0,
		Msg:  message,
		Data: data,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// ResponseHandler 提供 HTTP 处理程序的通用方法。
type ResponseHandler struct {
	ErrorHandler       *errors.ErrorHandler
	TranslationService *i18n.TranslationService
}

// NewResponseHandler 创建一个新的 ResponseHandler。
func NewResponseHandler(errHandler *errors.ErrorHandler, translationService *i18n.TranslationService) *ResponseHandler {
	return &ResponseHandler{
		ErrorHandler:       errHandler,
		TranslationService: translationService,
	}
}

func (h *ResponseHandler) HandleError(ctx fiber.Ctx, err error) error {
	var (
		statusCode     int
		domainErr      *errors.DomainError
		validationErrs validator.ValidationErrors
	)

	if stdErr.As(err, &validationErrs) && len(validationErrs) > 0 {
		statusCode = http.StatusBadRequest
		domainErr = errors.ValidationErrorToDomainError(validationErrs[0])
	} else if stdErr.As(err, &domainErr) {
		switch domainErr.Code {
		case errors.ErrUnauthorized:
			statusCode = http.StatusUnauthorized
		case errors.ErrUserAlreadyExists:
			statusCode = http.StatusConflict
		case errors.ErrInvalidInput, errors.ErrUserInfoIncorrect:
			statusCode = http.StatusBadRequest
		case errors.ErrTxError, errors.ErrDatabaseError:
			statusCode = http.StatusInternalServerError
		default:
			statusCode = http.StatusInternalServerError
		}
	} else {
		statusCode = http.StatusInternalServerError
		domainErr = &errors.DomainError{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
			Stack:   errors.CaptureStack(2),
		}
	}

	// 翻译错误
	translationKey := domainErr.TranslationKey()
	message := h.TranslationService.T(ctx.Context(), translationKey, domainErr.Params)
	if message == translationKey {
		message = domainErr.Message
	}

	// 获取调试追踪
	debugTrace := h.ErrorHandler.GetErrorTrace(domainErr)

	response := Response{
		Code:  domainErr.Code,
		Msg:   message,
		Data:  nil,
		Debug: debugTrace,
	}

	return ctx.Status(statusCode).JSON(response)
}
