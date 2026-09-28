package response

import (
	stdErr "errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/hina1314/ddd/util/errors"
	"github.com/hina1314/ddd/util/i18n"
)

// Response 定义响应的结构。
type Response struct {
	Code      interface{}        `json:"code"`
	Msg       string             `json:"msg"`
	Data      interface{}        `json:"data,omitempty"`
	RequestID string             `json:"request_id,omitempty"`
	Debug     *errors.ErrorTrace `json:"debug,omitempty"`
}

func (h *ResponseHandler) Success(c fiber.Ctx, msg string, data interface{}) error {
	message := h.TranslationService.T(c.Context(), msg, nil)
	response := Response{
		Code:      0,
		Msg:       message,
		Data:      data,
		RequestID: requestid.FromContext(c),
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
		fiberErr       *fiber.Error
	)

	if stdErr.As(err, &validationErrs) && len(validationErrs) > 0 {
		statusCode = http.StatusBadRequest
		domainErr = errors.ValidationErrorToDomainError(validationErrs[0])
	} else if stdErr.As(err, &domainErr) {
		switch domainErr.Code {
		case errors.ErrUnauthorized:
			statusCode = http.StatusUnauthorized
		case errors.ErrRateLimited:
			statusCode = http.StatusTooManyRequests
		case errors.ErrUserAlreadyExists:
			statusCode = http.StatusConflict
		case errors.ErrUserNotFound:
			statusCode = http.StatusNotFound
		case errors.ErrInvalidInput, errors.ErrRequired, errors.ErrInteger, errors.ErrDateFormat,
			errors.ErrUserInfoIncorrect, errors.ErrMinLength, errors.ErrPhoneEmpty,
			errors.ErrPhoneFormat, errors.ErrEmailEmpty, errors.ErrEmailFormat,
			errors.ErrAlphaNumUnicode:
			statusCode = http.StatusBadRequest
		case errors.ErrTxError, errors.ErrDatabaseError:
			statusCode = http.StatusInternalServerError
		default:
			statusCode = http.StatusInternalServerError
		}
	} else if stdErr.As(err, &fiberErr) {
		statusCode = fiberErr.Code

		errorCode := errors.ErrorCode(
			fmt.Sprintf("HTTP_%d", statusCode),
		)

		// 5xx 对外统一显示内部错误，不泄露具体实现细节。
		if statusCode >= http.StatusInternalServerError {
			errorCode = errors.ErrInternalError
		}

		domainErr = &errors.DomainError{
			Code:    errorCode,
			Message: fiberErr.Message,
			Cause:   err,
			Stack:   errors.CaptureStack(2),
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
	if statusCode >= http.StatusBadRequest {
		info := &ErrorLog{
			Code:    domainErr.Code,
			Message: domainErr.Message,
		}
		if statusCode >= http.StatusInternalServerError {
			info.Err = err
		}

		// 4xx 记录业务错误原因，不记录请求体或未经脱敏的底层 Cause。
		ctx.Locals(errorLogKey{}, info)
	}

	response := Response{
		Code:      domainErr.Code,
		Msg:       message,
		Data:      nil,
		RequestID: requestid.FromContext(ctx),
		Debug:     debugTrace,
	}

	return ctx.Status(statusCode).JSON(response)
}
