package errors

import (
	"fmt"
)

// ErrorCode represents a specific type of error code.
type ErrorCode string

// DomainError represents a domain-specific error with stack trace information.
type DomainError struct {
	Code           ErrorCode
	Message        string
	Params         map[string]interface{}
	Cause          error
	Stack          *StackTrace // Captures stack at creation
	translationKey string
}

// New creates a new DomainError with the given code and message.
func New(code ErrorCode, message string) *DomainError {
	return &DomainError{
		Code:           code,
		Message:        message,
		Params:         make(map[string]interface{}),
		Stack:          CaptureStack(2), // Skip New and its caller
		translationKey: fmt.Sprintf("errors.%s", code),
	}
}

// Wrap wraps an existing error into a DomainError.
func Wrap(err error, code ErrorCode, message string) *DomainError {
	return &DomainError{
		Code:           code,
		Message:        message,
		Params:         make(map[string]interface{}),
		Cause:          err,
		Stack:          CaptureStack(2), // Skip Wrap and its caller
		translationKey: fmt.Sprintf("errors.%s", code),
	}
}

func (e *DomainError) Error() string {
	return e.Message
}

func (e *DomainError) TranslationKey() string {
	return fmt.Sprintf("errors.%s", e.Code)
}

func (e *DomainError) WithParams(params map[string]interface{}) *DomainError {
	if e.Params == nil {
		e.Params = make(map[string]interface{})
	}
	for k, v := range params {
		e.Params[k] = v
	}
	return e
}

func (e *DomainError) WithCause(cause error) *DomainError {
	e.Cause = cause
	return e
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

func (e *DomainError) Is(target error) bool {
	t, ok := target.(*DomainError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}
