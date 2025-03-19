package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"study/internal/app"
	"time"
)

type UserHandler struct {
	base        *BaseHandler
	userService *app.UserService
}

func NewUserHandler(userService *app.UserService) *UserHandler {
	return &UserHandler{
		base:        NewBaseHandler(),
		userService: userService,
	}
}

type CreateUserRequest struct {
	Phone    string `json:"phone" validate:"required,phone"`
	Username string `json:"username" validate:"required,alphanumunicode"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	Phone       string    `json:"phone"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	AccessToken string    `json:"access_token,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return h.base.ErrorResponse(c, err)
	}

	// 使用 validator 库验证（需引入 "github.com/go-playground/validator/v10"）
	// 这里假设添加了 phone 自定义验证规则
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return h.base.ErrorResponse(c, err)
	}

	newUser, err := h.userService.RegisterUser(c.Context(), req.Username, req.Phone, req.Email, req.Password)
	if err != nil {
		return h.base.ErrorResponse(c, err)
	}

	resp := UserResponse{
		Phone:     newUser.Phone,
		Username:  newUser.Username,
		Email:     newUser.Email.String(),
		CreatedAt: newUser.CreatedAt,
	}
	return c.JSON(resp)
}

//func (h *UserHandler) Login(c fiber.Ctx) error {
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
