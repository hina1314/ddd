package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"study/internal/api/handler/dto"
	"study/internal/api/response"
	"study/internal/app"
	"study/util/context"
	"study/util/errors"
)

// UserHandler 处理用户相关的 HTTP 请求。
type UserHandler struct {
	res         *response.ResponseHandler
	userService *app.UserService
	validator   *validator.Validate
}

// NewUserHandler 创建一个新的 UserHandler。
func NewUserHandler(userService *app.UserService, base *response.ResponseHandler, v *validator.Validate) *UserHandler {
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

	return h.res.SuccessResponse(c, "user.create", newUser)
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

	return h.res.SuccessResponse(c, "user.login", user)
}

func (h *UserHandler) Info(c fiber.Ctx) error {

	payload, err := context.GetAuthPayload(c)
	if err != nil {
		return h.res.HandleError(c, err)
	}

	return h.res.SuccessResponse(c, "user.info", payload)
}

//func (h *UserHandler) Update(c fiber.Ctx) error {
//	var req createUserRequest
//	if err := c.Bind().JSON(&req); err != nil {
//		return c.JSON(errorResponse(err))
//	}
//
//	// 手动验证手机号格式
//	if !regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(req.Phone) {
//		return c.JSON("手机号格式不正确")
//
//	}
//
//	hashPassword, err := util.HashPassword(req.Password)
//	if err != nil {
//		return c.JSON(errorResponse(err))
//	}
//	arg := model.CreateUserParams{
//		Username: req.Username,
//		Phone:    req.Phone,
//		Email:    req.Email,
//		Password: hashPassword,
//	}
//	user, err := h.userService.RegisterUser(c.Context(), arg.Username, arg.Phone, arg.Email, arg.Password)
//	if err != nil {
//		//如果是数据库出错
//		if pqErr, ok := err.(*pq.Error); ok {
//			switch pqErr.Code.Name() {
//			case "unique_violation":
//				return c.JSON("手机号或用户名已存在")
//			}
//		}
//		return c.JSON(errorResponse(err))
//	}
//	resp := userResponse{
//		Phone:     user.Phone,
//		Username:  user.Username,
//		CreatedAt: user.CreatedAt,
//	}
//	return c.JSON(resp)
//}
