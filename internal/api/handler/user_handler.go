package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"study/internal/api/handler/dto"
	"study/internal/api/response"
	"study/internal/app/assemble"
	"study/internal/app/user"
	"study/util/context"
	"study/util/errors"
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
		return h.res.HandleError(c, err)
	}

	switch req.Type {
	case 1: // phone
		if req.Phone == "" {
			return h.res.HandleError(c, errors.New(errors.ErrPhoneEmpty, "phone is empty"))
		}
		req.Email = ""
	case 2: // email
		if req.Email == "" {
			return h.res.HandleError(c, errors.New(errors.ErrEmailEmpty, "email is empty"))
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
		return h.res.HandleError(c, err)
	}

	switch req.Type {
	case 1: // phone
		if req.Phone == "" {
			return h.res.HandleError(c, errors.New(errors.ErrPhoneEmpty, "phone is empty"))
		}
		req.Email = ""
	case 2: // email
		if req.Email == "" {
			return h.res.HandleError(c, errors.New(errors.ErrEmailEmpty, "email is empty"))
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

	user, err := h.userService.GetUserByID(c.Context(), payload.UserId)
	if err != nil {
		return err
	}

	return h.res.Success(c, "user.info", user)
}

func (h *UserHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return h.res.HandleError(c, err)
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
