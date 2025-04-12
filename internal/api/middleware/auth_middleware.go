package middleware

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"strings"
	"study/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func Auth(tokenMaker token.Maker) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authorizationHeader := ctx.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			return ctx.Status(http.StatusUnauthorized).JSON(err)
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			return ctx.Status(http.StatusUnauthorized).JSON(err)
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			return ctx.Status(http.StatusUnauthorized).JSON(err)
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(err)
		}
		ctx.Locals(AuthorizationPayloadKey, payload)
		return ctx.Next()
	}
}
