package util

import (
	"fmt"
	"log"

	"github.com/zeromicro/x/errors"
)

var Debug = false

// 定义错误模块代码
const (
	UserModule    = 10000
	OrderModule   = 20000
	PaymentModule = 30000
)

// 预定义常见错误
var (
	ErrUserAlreadyExists = NewError(UserModule, 1, "phone is already registered")
	ErrUsernameTaken     = NewError(UserModule, 2, "username is taken")
	ErrUserNotFound      = NewError(UserModule, 3, "user not found")
)

// NewError 创建新错误
func NewError(module, code int, msg string) error {
	fullCode := module + code
	return errors.New(fullCode, msg)
}

// Errorf 统一错误处理
func Errorf(module, code int, msg string, err error) error {
	fullCode := module + code

	if !Debug {
		// 生产环境：只返回简洁的用户友好消息
		return errors.New(fullCode, msg)
	}

	// 开发环境：包含详细错误信息并记录日志
	detailedMsg := fmt.Sprintf("%s: %v", msg, err)
	log.Println("ERROR:", detailedMsg)
	return errors.New(fullCode, detailedMsg)
}
