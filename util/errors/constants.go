package errors

// common
const (
	ErrInvalidInput  ErrorCode = "INVALID_INPUT"
	ErrRequired      ErrorCode = "REQUIRED"
	ErrInternalError ErrorCode = "INTERNAL_ERROR"
	ErrTxError       ErrorCode = "Transaction_Failure"
	ErrDatabaseError ErrorCode = "Database_Error"
)

// 用户领域错误代码
const (
	ErrUnauthorized      ErrorCode = "USER_UNAUTHORIZED"
	ErrUserNotFound      ErrorCode = "USER_NOT_FOUND"
	ErrUserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"
	ErrMinLength         ErrorCode = "USER_MIN_LENGTH"
	ErrPhoneEmpty        ErrorCode = "USER_PHONE_EMPTY"
	ErrPhoneFormat       ErrorCode = "USER_PHONE_FORMAT"
	ErrEmailEmpty        ErrorCode = "USER_EMAIL_EMPTY"
	ErrEmailFormat       ErrorCode = "USER_EMAIL_FORMAT"
	ErrAlphaNumUnicode   ErrorCode = "USER_ALPHA_NUM_UNICODE"
)

// 订单领域错误代码
const (
	ErrOrderNotFound     ErrorCode = "ORDER_NOT_FOUND"
	ErrInsufficientStock ErrorCode = "ORDER_INSUFFICIENT_STOCK"
	ErrPaymentFailed     ErrorCode = "ORDER_PAYMENT_FAILED"
)
