package response

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hina1314/ddd/util/errors"
)

type errorLogKey struct{}

// ErrorLog 保存供请求日志使用的错误详情，不包含请求体或调试堆栈。
type ErrorLog struct {
	Code    errors.ErrorCode
	Message string
	Err     error // 仅为 5xx 保留原始错误。
}

// ErrorLogFromContext 获取统一错误处理器保存的请求错误详情。
func ErrorLogFromContext(c fiber.Ctx) *ErrorLog {
	info, _ := c.Locals(errorLogKey{}).(*ErrorLog)
	return info
}
