package handler

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/hina1314/ddd/internal/api/handler/dto"
	"github.com/hina1314/ddd/internal/app/assemble"
	"github.com/hina1314/ddd/internal/app/user"
	"github.com/hina1314/ddd/internal/domain/user/usererrors"
	"github.com/hina1314/ddd/util/context"
	"github.com/hina1314/kit/errors"
	"github.com/hina1314/kit/response"
)

// UserHandler 处理用户相关的 HTTP 请求。
type UserHandler struct {
	res         *response.ResponseHandler
	userService *user.UserService
	validator   *validator.Validate
}

// NewUserHandler 创建一个新的 UserHandler。
func NewUserHandler(userService *user.UserService, base *response.ResponseHandler, v *validator.Validate) *UserHandler {
	return &UserHandler{
		res:         base,
		userService: userService,
		validator:   v,
	}
}

// CreateUser 处理用户注册请求。
func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return h.res.HandleError(c, errors.Wrap(err, errors.ErrInvalidInput, "invalid request body"))
	}
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	switch req.Type {
	case 1: // phone
		if req.Phone == "" {
			return h.res.HandleError(c, errors.New(usererrors.ErrPhoneEmpty, "phone is empty"))
		}
		req.Email = ""
	case 2: // email
		if req.Email == "" {
			return h.res.HandleError(c, errors.New(usererrors.ErrEmailEmpty, "email is empty"))
		}
		req.Phone = ""
	default:
		return h.res.HandleError(c, errors.New(errors.ErrInvalidInput, "invalid login type"))
	}

	if err := h.validator.Struct(req); err != nil {
		return h.res.HandleError(c, err)
	}

	newUser, err := h.userService.RegisterUser(c.Context(), req.Phone, req.Email, req.Password)
	if err != nil {
		return h.res.HandleError(c, err)
	}

	return h.res.Success(c, "user.create", newUser)
}

func (h *UserHandler) Login(c fiber.Ctx) error {
	var req dto.LoginUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return h.res.HandleError(c, errors.Wrap(err, errors.ErrInvalidInput, "invalid request body"))
	}
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	switch req.Type {
	case 1: // phone
		if req.Phone == "" {
			return h.res.HandleError(c, errors.New(usererrors.ErrPhoneEmpty, "phone is empty"))
		}
		req.Email = ""
	case 2: // email
		if req.Email == "" {
			return h.res.HandleError(c, errors.New(usererrors.ErrEmailEmpty, "email is empty"))
		}
		req.Phone = ""
	default:
		return h.res.HandleError(c, errors.New(errors.ErrInvalidInput, "invalid login type"))
	}

	if err := h.validator.Struct(req); err != nil {
		return h.res.HandleError(c, err)
	}

	user, err := h.userService.LoginUser(c.Context(), req.Phone, req.Email, req.Password)
	if err != nil {
		return h.res.HandleError(c, err)
	}

	return h.res.Success(c, "user.login", user)
}

func (h *UserHandler) Info(c fiber.Ctx) error {
	payload, err := context.GetAuthPayload(c)
	if err != nil {
		return h.res.HandleError(c, err)
	}

	user, err := h.userService.GetUserByID(c.Context(), payload.UserID)
	if err != nil {
		return h.res.HandleError(c, err)
	}

	return h.res.Success(c, "user.info", user)
}

func (h *UserHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return h.res.HandleError(c, errors.Wrap(err, errors.ErrInvalidInput, "invalid request body"))
	}
	if req.Phone != nil {
		value := strings.TrimSpace(*req.Phone)
		req.Phone = &value
	}
	if req.Email != nil {
		value := strings.ToLower(strings.TrimSpace(*req.Email))
		req.Email = &value
	}
	if req.Username != nil {
		value := strings.TrimSpace(*req.Username)
		req.Username = &value
	}

	if err := h.validator.Struct(req); err != nil {
		return h.res.HandleError(c, err)
	}

	var payload, err = context.GetAuthPayloadFromContext(c.Context())
	if err != nil {
		return err
	}

	cmd := assemble.NewUpdateUserCommand(req, payload)
	user, err := h.userService.UpdateUser(c.Context(), cmd)
	if err != nil {
		return h.res.HandleError(c, err)
	}
	return h.res.Success(c, "user.update", user)
}
