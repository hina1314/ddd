package handler

import (
	"study/internal/svc"
)

type Handler struct {
	Svc *svc.ServiceContext
}

func NewHandler(svc *svc.ServiceContext) *Handler {
	return &Handler{
		Svc: svc,
	}
}

// 自定义手机号验证器
//var phoneValidator validator.Func = func(fl validator.FieldLevel) bool {
//	phone, ok := fl.Field().Interface().(string)
//	if !ok {
//		return false
//	}
//	// 中国手机号格式：11 位，以 1 开头，第二位是 3-9
//	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
//	return re.MatchString(phone)
//}
