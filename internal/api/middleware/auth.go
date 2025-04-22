package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"strings"
	"study/internal/api/response"
	"study/token"
	"study/util/errors"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func Auth(res *response.ResponseHandler, tokenMaker token.Maker) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authorizationHeader := ctx.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New(errors.ErrUnauthorized, "authorization header is not provided")
			return res.HandleError(ctx, err)
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New(errors.ErrUnauthorized, "invalid authorization header format")
			return res.HandleError(ctx, err)
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := errors.New(errors.ErrUnauthorized, fmt.Sprintf("unsupported authorization type %s", authorizationType))
			return res.HandleError(ctx, err)
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			err = errors.Wrap(err, errors.ErrUnauthorized, "invalid token")
			return res.HandleError(ctx, err)
		}

		ctx.Locals(AuthorizationPayloadKey, payload)
		return ctx.Next()
	}
}
