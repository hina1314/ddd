package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"study/internal/api/handler/dto"
	"study/internal/app"
	"study/util/errors"
	"study/util/i18n"
)

type UserHandler struct {
	base        *BaseHandler
	userService *app.UserService
}

func NewUserHandler(userService *app.UserService, errHandler *errors.ErrorHandler, translationService *i18n.TranslationService) *UserHandler {
	return &UserHandler{
		base:        NewBaseHandler(errHandler, translationService),
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return h.base.handleError(c, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return h.base.handleError(c, err)
	}

	newUser, err := h.userService.RegisterUser(c.Context(), req.Username, req.Phone, req.Email, req.Password)
	if err != nil {
		return h.base.handleError(c, err)
	}

	return h.base.successResponse(c, "user.created", newUser)
}

//func (h *UserHandler) Login(c fiber.Ctx) error {
//	var req dto.LoginUserRequest
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
//	phoneOrEmail := req.Phone
//	if phoneOrEmail == "" {
//		phoneOrEmail = req.Email
//	}
//	user, err := h.userService.LoginUser(c.Context(), phoneOrEmail, req.Password)
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
//
//	return c.JSON(user)
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
