package middleware

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// Logger 这里有几个生产日志原则：
// - duration_ms 使用数字，而不是 "1.2ms" 字符串，方便统计和排序。
// - 不记录 URL 查询参数，避免密码、Token 等进入日志。
// - 2xx/3xx 使用 INFO。
// - 4xx 使用 WARN。
// - 5xx 使用 ERROR。
func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		err := c.Next()
		status := responseStatus(c, err)
		duration := time.Since(start)

		// 成功的内部探测请求不写访问日志；失败时仍然记录。
		if isQuietPath(c.Path()) && status < fiber.StatusBadRequest {
			return err
		}

		level := slog.LevelInfo
		switch {
		case status >= fiber.StatusInternalServerError:
			level = slog.LevelError
		case status >= fiber.StatusBadRequest:
			level = slog.LevelWarn
		}

		attributes := []slog.Attr{
			slog.String("request_id", requestid.FromContext(c)),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.Float64(
				"duration_ms",
				float64(duration.Microseconds())/1000,
			),
		}

		if err != nil {
			attributes = append(
				attributes,
				slog.String("error", err.Error()),
			)
		}

		slog.LogAttrs(
			c.Context(),
			level,
			"http_request",
			attributes...,
		)

		return err
	}
}

func responseStatus(c fiber.Ctx, err error) int {
	if err == nil {
		return c.Response().StatusCode()
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return fiberErr.Code
	}

	return fiber.StatusInternalServerError
}

func isQuietPath(path string) bool {
	switch path {
	case "/metrics", "/livez", "/readyz":
		return true
	default:
		return false
	}
}
