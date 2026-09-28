package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/hina1314/ddd/internal/domain/user/usererrors"
	"github.com/hina1314/ddd/token"
	"github.com/hina1314/ddd/util/context"
	"github.com/hina1314/kit/errors"
	"github.com/hina1314/kit/response"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

func Auth(res *response.ResponseHandler, tokenMaker token.Maker) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authorizationHeader := ctx.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New(usererrors.ErrUnauthorized, "authorization header is not provided")
			return res.HandleError(ctx, err)
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			err := errors.New(usererrors.ErrUnauthorized, "invalid authorization header format")
			return res.HandleError(ctx, err)
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := errors.New(usererrors.ErrUnauthorized, fmt.Sprintf("unsupported authorization type %s", authorizationType))
			return res.HandleError(ctx, err)
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			err = errors.Wrap(err, usererrors.ErrUnauthorized, "invalid token")
			return res.HandleError(ctx, err)
		}

		ctx.Locals(context.AuthorizationPayloadKey, payload)
		newCtx := context.WithAuthPayload(ctx.Context(), payload)
		ctx.SetContext(newCtx)

		return ctx.Next()
	}
}
