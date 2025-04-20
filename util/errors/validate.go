package errors

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

// PhoneValidator 自定义校验函数
func PhoneValidator(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}

// ValidationErrorToDomainError 将单个验证错误转换为 DomainError。
func ValidationErrorToDomainError(ve validator.FieldError) *DomainError {
	field := strings.ToLower(ve.Field())
	params := map[string]interface{}{
		"field": field,
	}

	switch ve.Tag() {
	case "required":
		return New(ErrRequired, fmt.Sprintf("field %s is required", field)).
			WithParams(params)
	case "phone":
		return New(ErrPhone, fmt.Sprintf("field %s is incorrect", field))
	case "email":
		return New(ErrEmail, fmt.Sprintf("field %s must be a valid email address", field)).
			WithParams(params)
	case "min":
		params["min"] = ve.Param()
		return New(ErrMinLength, fmt.Sprintf("field %s must be at least %s characters", field, ve.Param())).
			WithParams(params)
	case "alphanumunicode":
		return New(ErrAlphaNumUnicode, fmt.Sprintf("field %s must contain only alphanumeric characters", field)).
			WithParams(params)
	default:
		return New(ErrInvalidInput, fmt.Sprintf("field %s is invalid: %s", field, ve.Tag())).
			WithParams(params)
	}
}
