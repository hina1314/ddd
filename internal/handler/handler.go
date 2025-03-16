package handler

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
	"study/internal/svc"
)

type Handler struct {
	svc *svc.ServiceContext
}

func NewHandler(svc *svc.ServiceContext) *Handler {
	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("phone", phoneValidator)
		if err != nil {
			log.Fatal(err)
		}
	}
	return &Handler{
		svc: svc,
	}
}

// 自定义手机号验证器
var phoneValidator validator.Func = func(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	// 中国手机号格式：11 位，以 1 开头，第二位是 3-9
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(phone)
}
