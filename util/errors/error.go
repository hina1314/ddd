package errors

// ErrorCode 表示特定类型的错误代码
type ErrorCode string

// DomainError 表示领域错误
type DomainError struct {
	Code    ErrorCode
	Message string
	Params  map[string]interface{}
	Cause   error
}

func (e *DomainError) Error() string {
	return e.Message
}

// WithParams 添加参数到错误
func (e *DomainError) WithParams(params map[string]interface{}) *DomainError {
	if e.Params == nil {
		e.Params = make(map[string]interface{})
	}
	for k, v := range params {
		e.Params[k] = v
	}
	return e
}

// WithCause 添加原因错误
func (e *DomainError) WithCause(cause error) *DomainError {
	e.Cause = cause
	return e
}

// Unwrap 实现errors.Unwrap接口
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// Is 实现errors.Is接口
func (e *DomainError) Is(target error) bool {
	t, ok := target.(*DomainError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// New 创建新的领域错误
func New(code ErrorCode, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Params:  make(map[string]interface{}),
	}
}

// Wrap 包装一个已有错误
func Wrap(err error, code ErrorCode, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Params:  make(map[string]interface{}),
		Cause:   err,
	}
}
