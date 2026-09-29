// Package presentation configures the framework for this application's API contract.
package presentation

import (
	"net/http"

	"github.com/hina1314/ddd/internal/api/validation"
	"github.com/hina1314/ddd/internal/domain/user/usererrors"
	"github.com/hina1314/kit/errors"
	"github.com/hina1314/kit/i18n"
	"github.com/hina1314/kit/response"
)

func NewResponseHandler(handler *errors.ErrorHandler, translator *i18n.TranslationService) *response.ResponseHandler {
	return response.NewResponseHandler(handler, translator, response.Options{
		StatusCodes: map[errors.ErrorCode]int{
			usererrors.ErrUnauthorized:      http.StatusUnauthorized,
			usererrors.ErrUserNotFound:      http.StatusNotFound,
			usererrors.ErrUserAlreadyExists: http.StatusConflict,
			usererrors.ErrUserInfoIncorrect: http.StatusBadRequest,
			usererrors.ErrMinLength:         http.StatusBadRequest,
			usererrors.ErrPasswordTooLong:   http.StatusBadRequest,
			usererrors.ErrPhoneEmpty:        http.StatusBadRequest,
			usererrors.ErrPhoneFormat:       http.StatusBadRequest,
			usererrors.ErrEmailEmpty:        http.StatusBadRequest,
			usererrors.ErrEmailFormat:       http.StatusBadRequest,
			usererrors.ErrAlphaNumUnicode:   http.StatusBadRequest,
		},
		ValidationMapper: validation.ValidationErrorToDomainError,
	})
}
