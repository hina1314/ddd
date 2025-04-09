package errors

// 用户领域错误代码
const (
	ErrUserNotFound      ErrorCode = "USER_NOT_FOUND"
	ErrUserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"
	ErrInvalidEmail      ErrorCode = "INVALID_EMAIL"
	ErrWeakPassword      ErrorCode = "WEAK_PASSWORD"
)

// 订单领域错误代码
const (
	ErrOrderNotFound     ErrorCode = "ORDER_NOT_FOUND"
	ErrInsufficientStock ErrorCode = "INSUFFICIENT_STOCK"
	ErrPaymentFailed     ErrorCode = "PAYMENT_FAILED"
)
