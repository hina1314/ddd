package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/hina1314/ddd/internal/api/handler"
	"github.com/hina1314/ddd/internal/api/presentation"
	"github.com/hina1314/ddd/internal/api/validation"
	"github.com/hina1314/kit/errors"
	"github.com/hina1314/kit/i18n"
	"github.com/hina1314/kit/middleware"
	"github.com/hina1314/kit/response"
	"github.com/stretchr/testify/require"
)

func TestUserInputResponsesRemainCompatible(t *testing.T) {
	translator := i18n.NewFileTranslator("zh")
	require.NoError(t, translator.LoadTranslations("../../../config/i18n"))
	responder := presentation.NewResponseHandler(errors.NewErrorHandler(false, false), i18n.NewTranslationService(translator, "zh"))
	v := validator.New()
	require.NoError(t, v.RegisterValidation("phone", validation.PhoneValidator))
	require.NoError(t, v.RegisterValidation("password", validation.PasswordValidator))
	userHandler := handler.NewUserHandler(nil, responder, v)

	cases := []struct {
		name, body, code, message string
	}{
		{"malformed JSON", `{"type":1,"password":"12345678}`, "INVALID_INPUT", "参数无效！"},
		{"invalid type", `{"type":3}`, "INVALID_INPUT", "参数无效！"},
		{"missing phone", `{"type":1}`, "USER_PHONE_EMPTY", "手机号不能为空！"},
		{"missing email", `{"type":2}`, "USER_EMAIL_EMPTY", "邮箱不能为空！"},
		{"invalid phone", `{"type":1,"phone":"123","password":"12345678"}`, "USER_PHONE_FORMAT", "手机号格式不正确！"},
		{"invalid email", `{"type":2,"email":"invalid","password":"12345678"}`, "USER_EMAIL_FORMAT", "邮箱格式不正确！"},
		{"short password", `{"type":1,"phone":"13458667726","password":"123"}`, "USER_MIN_LENGTH", "password的长度最少为8位！"},
		{"password over bcrypt byte limit", `{"type":1,"phone":"13458667726","password":"密密密密密密密密密密密密密密密密密密密密密密密密密"}`, "USER_PASSWORD_TOO_LONG", "password不能超过72字节！"},
		{"missing password", `{"type":1,"phone":"13458667726"}`, "REQUIRED", "password为必填项！"},
	}
	for _, path := range []string{"/v1/login", "/v1/signup"} {
		for _, tt := range cases {
			t.Run(path+"/"+tt.name, func(t *testing.T) {
				var logs bytes.Buffer
				app := fiber.New(fiber.Config{ErrorHandler: responder.HandleError})
				app.Use(requestid.New())
				app.Use(middleware.Logger(middleware.LoggerConfig{Logger: slog.New(slog.NewJSONHandler(&logs, nil))}))
				app.Use(recover.New())
				app.Post("/v1/login", userHandler.Login)
				app.Post("/v1/signup", userHandler.CreateUser)
				req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(tt.body))
				req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
				res, err := app.Test(req)
				require.NoError(t, err)
				defer res.Body.Close()
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				var body response.Response
				require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
				require.Equal(t, tt.code, body.Code)
				require.Equal(t, tt.message, body.Msg)
				require.Nil(t, body.Debug)
				var record map[string]interface{}
				decoder := json.NewDecoder(&logs)
				require.NoError(t, decoder.Decode(&record))
				require.Equal(t, "http_request", record["msg"])
				require.Equal(t, tt.code, record["error_code"])
				require.Equal(t, body.RequestID, record["request_id"])
				require.ErrorIs(t, decoder.Decode(&record), io.EOF)
			})
		}
	}
}
