package errors

// common
const (
	ErrInvalidInput  ErrorCode = "Invalid_Input"
	ErrInternalError ErrorCode = "Internal_Error"
	ErrTxError       ErrorCode = "Transaction_Failure"
	ErrDatabaseError ErrorCode = "Database_Error"
)

// 用户领域错误代码
const (
	ErrUserNotFound      ErrorCode = "USER_NOT_FOUND"
	ErrUserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"
	ErrInvalidEmail      ErrorCode = "USER_INVALID_EMAIL"
	ErrWeakPassword      ErrorCode = "USER_WEAK_PASSWORD"
)

// 订单领域错误代码
const (
	ErrOrderNotFound     ErrorCode = "ORDER_NOT_FOUND"
	ErrInsufficientStock ErrorCode = "ORDER_INSUFFICIENT_STOCK"
	ErrPaymentFailed     ErrorCode = "ORDER_PAYMENT_FAILED"
)
