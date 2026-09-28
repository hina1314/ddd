package response_test

import (
	"bytes"
	"encoding/json"
	stdErr "errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/hina1314/ddd/internal/api/middleware"
	"github.com/hina1314/ddd/internal/api/response"
	"github.com/hina1314/ddd/util/errors"
	"github.com/hina1314/ddd/util/i18n"
	"github.com/stretchr/testify/require"
)

func TestHandleErrorLogsRequestFailure(t *testing.T) {
	translator := i18n.NewFileTranslator("zh")
	require.NoError(t, translator.LoadTranslations("../../../config/i18n"))
	h := response.NewResponseHandler(
		errors.NewErrorHandler(false, false),
		i18n.NewTranslationService(translator, "zh"),
	)

	tests := []struct {
		name     string
		err      error
		status   int
		level    string
		code     errors.ErrorCode
		returned bool
		panics   bool
	}{
		{
			name:   "invalid request body",
			err:    errors.Wrap(stdErr.New("sensitive-password-in-cause"), errors.ErrInvalidInput, "invalid request body"),
			status: http.StatusBadRequest,
			level:  "WARN",
			code:   errors.ErrInvalidInput,
		},
		{
			name:   "invalid login type",
			err:    errors.New(errors.ErrInvalidInput, "invalid login type"),
			status: http.StatusBadRequest,
			level:  "WARN",
			code:   errors.ErrInvalidInput,
		},
		{
			name:   "internal failure",
			err:    stdErr.New("internal failure"),
			status: http.StatusInternalServerError,
			level:  "ERROR",
			code:   errors.ErrInternalError,
		},
		{
			name:     "returned domain error",
			err:      errors.New(errors.ErrInvalidInput, "invalid login type"),
			status:   http.StatusBadRequest,
			level:    "WARN",
			code:     errors.ErrInvalidInput,
			returned: true,
		},
		{
			name:     "returned internal error",
			err:      stdErr.New("internal failure"),
			status:   http.StatusInternalServerError,
			level:    "ERROR",
			code:     errors.ErrInternalError,
			returned: true,
		},
		{
			name:   "recovered panic",
			err:    stdErr.New("panic failure"),
			status: http.StatusInternalServerError,
			level:  "ERROR",
			code:   errors.ErrInternalError,
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			previousLogger := slog.Default()
			slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, nil)))
			t.Cleanup(func() { slog.SetDefault(previousLogger) })

			errorHandlerCalls := 0
			app := fiber.New(fiber.Config{
				ErrorHandler: func(c fiber.Ctx, err error) error {
					errorHandlerCalls++
					return h.HandleError(c, err)
				},
			})
			app.Use(requestid.New())
			app.Use(middleware.Logger())
			app.Use(recover.New())
			app.Post("/v1/login", func(c fiber.Ctx) error {
				if tt.panics {
					panic(tt.err)
				}
				if tt.returned {
					return tt.err
				}
				return h.HandleError(c, tt.err)
			})
			app.Get("/ok", func(c fiber.Ctx) error {
				return h.Success(c, "user.login", nil)
			})

			res, err := app.Test(httptest.NewRequest(http.MethodPost, "/v1/login", nil))
			require.NoError(t, err)
			defer res.Body.Close()
			require.Equal(t, tt.status, res.StatusCode)
			if tt.returned || tt.panics {
				require.Equal(t, 1, errorHandlerCalls)
			} else {
				require.Zero(t, errorHandlerCalls)
			}

			var body response.Response
			require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
			require.Equal(t, string(tt.code), body.Code)
			require.NotEmpty(t, body.RequestID)
			require.Nil(t, body.Debug)
			if tt.status == http.StatusBadRequest {
				require.Equal(t, translator.T("errors.INVALID_INPUT", "zh", nil), body.Msg)
			}

			require.NotContains(t, logs.String(), "sensitive-password-in-cause")
			decoder := json.NewDecoder(&logs)
			var failure map[string]interface{}
			require.NoError(t, decoder.Decode(&failure))
			require.Equal(t, "http_request", failure["msg"])
			require.Equal(t, tt.level, failure["level"])
			require.Equal(t, string(tt.code), failure["error_code"])
			require.Equal(t, tt.err.Error(), failure["error_message"])
			require.Equal(t, float64(tt.status), failure["status"])
			require.Equal(t, body.RequestID, failure["request_id"])
			require.Equal(t, http.MethodPost, failure["method"])
			require.Equal(t, "/v1/login", failure["path"])
			require.Contains(t, failure, "duration_ms")
			if tt.status == http.StatusBadRequest {
				require.NotContains(t, failure, "error")
				require.NotContains(t, failure, "cause")
			} else {
				require.Equal(t, tt.err.Error(), failure["error"])
			}
			require.ErrorIs(t, decoder.Decode(&failure), io.EOF)

			// 后续成功请求也只有一条日志，并且不会继承前一个请求的错误详情。
			logs.Reset()
			okRes, err := app.Test(httptest.NewRequest(http.MethodGet, "/ok", nil))
			require.NoError(t, err)
			defer okRes.Body.Close()
			require.Equal(t, http.StatusOK, okRes.StatusCode)
			decoder = json.NewDecoder(&logs)
			var success map[string]interface{}
			require.NoError(t, decoder.Decode(&success))
			require.Equal(t, "http_request", success["msg"])
			require.Equal(t, "INFO", success["level"])
			require.NotContains(t, success, "error_code")
			require.NotContains(t, success, "error_message")
			require.NotContains(t, success, "error")
			require.NotEqual(t, body.RequestID, success["request_id"])
			require.ErrorIs(t, decoder.Decode(&success), io.EOF)
		})
	}
}
