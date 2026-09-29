package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hina1314/ddd/internal/domain/user/passwordpolicy"
	"github.com/hina1314/ddd/internal/domain/user/usererrors"
	"github.com/hina1314/kit/errors"
)

var mainlandChinaPhonePattern = regexp.MustCompile(`^1[3-9]\d{9}$`)

// PhoneValidator 自定义校验函数
func PhoneValidator(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return mainlandChinaPhonePattern.MatchString(phone)
}

// PasswordValidator applies the same password policy to registration, login,
// and profile updates.
func PasswordValidator(fl validator.FieldLevel) bool {
	return passwordpolicy.Valid(fl.Field().String())
}

// ValidationErrorToDomainError 将单个验证错误转换为 DomainError。
func ValidationErrorToDomainError(ve validator.FieldError) *errors.DomainError {
	field := strings.ToLower(ve.Field())
	params := map[string]interface{}{
		"field": field,
	}

	switch ve.Tag() {
	case "required":
		return errors.New(errors.ErrRequired, fmt.Sprintf("field %s is required", field)).
			WithParams(params)
	case "phone":
		return errors.New(usererrors.ErrPhoneFormat, fmt.Sprintf("field %s is incorrect", field))
	case "email":
		return errors.New(usererrors.ErrEmailFormat, fmt.Sprintf("field %s must be a valid email address", field)).
			WithParams(params)
	case "min":
		params["min"] = ve.Param()
		return errors.New(usererrors.ErrMinLength, fmt.Sprintf("field %s must be at least %s characters", field, ve.Param())).
			WithParams(params)
	case "password":
		password, _ := ve.Value().(string)
		if passwordpolicy.TooLong(password) {
			params["max"] = passwordpolicy.MaxBytes
			return errors.New(usererrors.ErrPasswordTooLong, fmt.Sprintf("field %s must contain at most %d bytes", field, passwordpolicy.MaxBytes)).
				WithParams(params)
		}
		params["min"] = passwordpolicy.MinCharacters
		return errors.New(usererrors.ErrMinLength, fmt.Sprintf("field %s must be at least %d characters", field, passwordpolicy.MinCharacters)).
			WithParams(params)
	case "alphanumunicode":
		return errors.New(usererrors.ErrAlphaNumUnicode, fmt.Sprintf("field %s must contain only alphanumeric characters", field)).
			WithParams(params)
	case "numeric":
		return errors.New(errors.ErrInteger, fmt.Sprintf("field %s must be a valid integer", field)).WithParams(params)
	default:
		return errors.New(errors.ErrInvalidInput, fmt.Sprintf("field %s is invalid: %s", field, ve.Tag())).
			WithParams(params)
	}
}
