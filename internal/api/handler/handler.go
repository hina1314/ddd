package handler

import (
	stdErr "errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"study/util/errors"
	"study/util/i18n"
)

// BaseHandler 提供 HTTP 处理程序的通用方法。
type BaseHandler struct {
	ErrorHandler       *errors.ErrorHandler
	TranslationService *i18n.TranslationService
}

// NewBaseHandler 创建一个新的 BaseHandler。
func NewBaseHandler(errHandler *errors.ErrorHandler, translationService *i18n.TranslationService) *BaseHandler {
	return &BaseHandler{
		ErrorHandler:       errHandler,
		TranslationService: translationService,
	}
}

func (h *BaseHandler) handleError(ctx fiber.Ctx, err error) error {
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
		case errors.ErrUserAlreadyExists:
			statusCode = http.StatusConflict
		case errors.ErrInvalidInput:
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

	response := errorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    string(domainErr.Code),
			Message: message,
		},
		Debug: debugTrace,
	}

	return ctx.Status(statusCode).JSON(response)
}
