package errors

import (
	"context"
	"study/util/i18n"
)

// ErrorHandler 负责处理和转换领域错误
type ErrorHandler struct {
	debugMode bool
}

// NewErrorHandler 创建一个新的错误处理器
func NewErrorHandler(debugMode bool) *ErrorHandler {
	return &ErrorHandler{
		debugMode: debugMode,
	}
}

// SetDebugMode 设置调试模式
func (h *ErrorHandler) SetDebugMode(debug bool) {
	h.debugMode = debug
}

// Handle 处理领域错误，包括翻译
func (h *ErrorHandler) Handle(ctx context.Context, err error) (string, *ErrorTrace) {
	if err == nil {
		return "", nil
	}

	// 提取领域错误
	var domainErr *DomainError
	if de, ok := err.(*DomainError); ok {
		domainErr = de
	} else {
		// 包装为通用错误
		domainErr = &DomainError{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}
	}

	// 获取错误追踪（仅在调试模式下）
	var trace *ErrorTrace
	if h.debugMode {
		trace = GetErrorTrace(err)
	}

	// 翻译错误消息
	translationKey := domainErr.TranslationKey()
	message := i18n.T(ctx, translationKey, domainErr.Params)

	// 如果翻译失败（返回键本身），使用默认消息
	if message == translationKey {
		message = domainErr.Message
	}

	return message, trace
}
