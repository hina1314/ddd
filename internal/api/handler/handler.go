package handler

import (
	stdErr "errors"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"study/util/errors"
	"study/util/i18n"
)

// BaseHandler provides common methods for HTTP handlers.
type BaseHandler struct {
	ErrorHandler       *errors.ErrorHandler
	TranslationService *i18n.TranslationService
}

// NewBaseHandler creates a new BaseHandler with the given dependencies.
func NewBaseHandler(errHandler *errors.ErrorHandler, translationService *i18n.TranslationService) *BaseHandler {
	return &BaseHandler{
		ErrorHandler:       errHandler,
		TranslationService: translationService,
	}
}

func (h *BaseHandler) handleError(ctx fiber.Ctx, err error) error {
	var statusCode int
	var domainErr *errors.DomainError
	if stdErr.As(err, &domainErr) {
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
	
	// Perform translation
	translationKey := domainErr.TranslationKey()

	message := h.TranslationService.T(ctx.Context(), translationKey, domainErr.Params)
	if message == translationKey {
		message = domainErr.Message
	}

	// Get debug trace
	debugTrace := h.ErrorHandler.GetErrorTrace(err)

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
