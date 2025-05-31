package errors

// common
const (
	ErrInvalidInput  ErrorCode = "INVALID_INPUT"
	ErrRequired      ErrorCode = "REQUIRED"
	ErrInternalError ErrorCode = "INTERNAL_ERROR"
	ErrTxError       ErrorCode = "Transaction_Failure"
	ErrDatabaseError ErrorCode = "Database_Error"
	ErrInteger       ErrorCode = "ERROR_INTEGER"
	ErrDateFormat    ErrorCode = "ERROR_DATE_FORMAT"
)

// 用户领域错误代码
const (
	ErrUnauthorized      ErrorCode = "USER_UNAUTHORIZED"
	ErrUserNotFound      ErrorCode = "USER_NOT_FOUND"
	ErrUserInfoIncorrect ErrorCode = "USER_INFO_INCORRECT"
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
	ErrStartDatePast     ErrorCode = "ORDER_START_DATE_PAST"
	ErrStartDateDisorder ErrorCode = "ORDER_START_DATE_DISORDER"
	ErrBookingConflict   ErrorCode = "ORDER_BOOKING_CONFLICT"
)

// 酒店领域错误

const (
	ErrHotelSkuNotFound ErrorCode = "HOTEL_SKU_NOT_FOUND"
	ErrNoSkuPrice       ErrorCode = "HOTEL_NO_SKU_PRICE"
	ErrNoStock          ErrorCode = "HOTEL_NO_STOCK"
	ErrTicketNotSupport ErrorCode = "HOTEL_TICKET_NOT_SUPPORT"
)
