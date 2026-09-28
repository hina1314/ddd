package presentation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/hina1314/ddd/internal/domain/user/usererrors"
	"github.com/hina1314/kit/errors"
	"github.com/hina1314/kit/i18n"
	"github.com/stretchr/testify/require"
)

func TestUserErrorStatusMappingsRemainCompatible(t *testing.T) {
	h := NewResponseHandler(errors.NewErrorHandler(false, false), i18n.NewTranslationService(i18n.NewFileTranslator("en"), "en"))
	for code, status := range map[errors.ErrorCode]int{
		usererrors.ErrUnauthorized:      http.StatusUnauthorized,
		usererrors.ErrUserNotFound:      http.StatusNotFound,
		usererrors.ErrUserAlreadyExists: http.StatusConflict,
		usererrors.ErrUserInfoIncorrect: http.StatusBadRequest,
		usererrors.ErrMinLength:         http.StatusBadRequest,
		usererrors.ErrPhoneEmpty:        http.StatusBadRequest,
		usererrors.ErrPhoneFormat:       http.StatusBadRequest,
		usererrors.ErrEmailEmpty:        http.StatusBadRequest,
		usererrors.ErrEmailFormat:       http.StatusBadRequest,
		usererrors.ErrAlphaNumUnicode:   http.StatusBadRequest,
	} {
		t.Run(string(code), func(t *testing.T) {
			app := fiber.New(fiber.Config{ErrorHandler: h.HandleError})
			app.Get("/", func(c fiber.Ctx) error { return errors.New(code, "business failure") })
			res, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
			require.NoError(t, err)
			defer res.Body.Close()
			require.Equal(t, status, res.StatusCode)
		})
	}
}
