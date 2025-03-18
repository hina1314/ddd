package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/lib/pq"
	"regexp"
	"study/db/model"
	"study/internal/app"
	"study/util"
	"time"
)

type UserHandler struct {
	Handler     *Handler
	userService *app.UserService
}

func NewUserHandler(h *Handler, userSvc *app.UserService) *UserHandler {
	return &UserHandler{Handler: h, userService: userSvc}
}

type createUserRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Username string `json:"username" binding:"required,alphanumunicode"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Phone       string    `json:"phone"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
}

func newUserResponse(user model.User) userResponse {
	return userResponse{
		Phone:     user.Phone,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}
func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var req createUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.JSON(errorResponse(err))
	}

	// 手动验证手机号格式
	if !regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(req.Phone) {
		return c.JSON("手机号格式不正确")

	}

	arg := model.CreateUserParams{
		Username: req.Username,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
	}
	user, err := h.userService.RegisterUser(c.Context(), arg.Username, arg.Phone, arg.Email, arg.Password)
	if err != nil {
		//如果是数据库出错
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return c.JSON("手机号或用户名已存在")
			}
		}
		return c.JSON(errorResponse(err))
	}
	resp := userResponse{
		Phone:     user.Phone,
		Username:  user.Username,
		Email:     user.Email.String(),
		CreatedAt: user.CreatedAt,
	}
	return c.JSON(resp)
}

func (h *UserHandler) Login(c fiber.Ctx) error {
	var req createUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.JSON(errorResponse(err))
	}

	// 手动验证手机号格式
	if !regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(req.Phone) {
		return c.JSON("手机号格式不正确")

	}

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return c.JSON(errorResponse(err))
	}
	arg := model.CreateUserParams{
		Username: req.Username,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: hashPassword,
	}
	user, err := h.userService.RegisterUser(c.Context(), arg.Username, arg.Phone, arg.Email, arg.Password)
	if err != nil {
		//如果是数据库出错
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return c.JSON("手机号或用户名已存在")
			}
		}
		return c.JSON(errorResponse(err))
	}
	resp := userResponse{
		Phone:     user.Phone,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
	return c.JSON(resp)
}

func (h *UserHandler) Update(c fiber.Ctx) error {
	var req createUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.JSON(errorResponse(err))
	}

	// 手动验证手机号格式
	if !regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(req.Phone) {
		return c.JSON("手机号格式不正确")

	}

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return c.JSON(errorResponse(err))
	}
	arg := model.CreateUserParams{
		Username: req.Username,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: hashPassword,
	}
	user, err := h.userService.RegisterUser(c.Context(), arg.Username, arg.Phone, arg.Email, arg.Password)
	if err != nil {
		//如果是数据库出错
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return c.JSON("手机号或用户名已存在")
			}
		}
		return c.JSON(errorResponse(err))
	}
	resp := userResponse{
		Phone:     user.Phone,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
	return c.JSON(resp)
}
