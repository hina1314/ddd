package context

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"study/token"
	"study/util/errors"
)

// contextKey 定义上下文键类型，防止冲突。
type contextKey string

const (
	AuthPayloadKey          contextKey = "auth_payload"
	AuthorizationPayloadKey            = "authorization_payload"
)

// WithAuthPayload 将 token.Payload 存储到 context.Context。
func WithAuthPayload(ctx context.Context, payload *token.Payload) context.Context {
	return context.WithValue(ctx, AuthPayloadKey, payload)
}

// GetAuthPayloadFromContext 从 context.Context 获取 token.Payload。
func GetAuthPayloadFromContext(ctx context.Context) (*token.Payload, error) {
	payload, ok := ctx.Value(AuthPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return nil, errors.New(errors.ErrUnauthorized, "authentication payload not found in context")
	}
	return payload, nil
}

// GetAuthPayload 从 fiber.Ctx 获取 token.Payload。
func GetAuthPayload(ctx fiber.Ctx) (*token.Payload, error) {
	payload, ok := ctx.Locals(AuthorizationPayloadKey).(*token.Payload)
	if !ok || payload == nil {
		return nil, errors.New(errors.ErrUnauthorized, "authentication payload not found")
	}
	return payload, nil
}
