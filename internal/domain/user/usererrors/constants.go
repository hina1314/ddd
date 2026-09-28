// Package usererrors owns the user API's error codes independently of the framework.
package usererrors

import "github.com/hina1314/kit/errors"

// Keep the existing wire values to preserve compatibility with API clients.
const (
	ErrUnauthorized      errors.ErrorCode = "USER_UNAUTHORIZED"
	ErrUserNotFound      errors.ErrorCode = "USER_NOT_FOUND"
	ErrUserInfoIncorrect errors.ErrorCode = "USER_INFO_INCORRECT"
	ErrUserAlreadyExists errors.ErrorCode = "USER_ALREADY_EXISTS"
	ErrMinLength         errors.ErrorCode = "USER_MIN_LENGTH"
	ErrPhoneEmpty        errors.ErrorCode = "USER_PHONE_EMPTY"
	ErrPhoneFormat       errors.ErrorCode = "USER_PHONE_FORMAT"
	ErrEmailEmpty        errors.ErrorCode = "USER_EMAIL_EMPTY"
	ErrEmailFormat       errors.ErrorCode = "USER_EMAIL_FORMAT"
	ErrAlphaNumUnicode   errors.ErrorCode = "USER_ALPHA_NUM_UNICODE"
)
