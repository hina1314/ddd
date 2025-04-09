package errors

import "context"

// ErrorHandler 负责处理和转换领域错误
type ErrorHandler struct {
	translator Translator
	locale     string
}

// NewErrorHandler 创建一个新的错误处理器
func NewErrorHandler(translator Translator, defaultLocale string) *ErrorHandler {
	return &ErrorHandler{
		translator: translator,
		locale:     defaultLocale,
	}
}

// Handle 处理领域错误，包括翻译
func (h *ErrorHandler) Handle(ctx context.Context, err error) string {
	if err == nil {
		return ""
	}

	// 获取上下文中的语言设置
	locale := h.localeFromContext(ctx)

	// 检查是否是领域错误
	var domainErr *DomainError
	if de, ok := err.(*DomainError); ok {
		domainErr = de
	} else {
		// 包装为通用错误
		domainErr = &DomainError{
			Code:    "UNKNOWN_ERROR",
			Message: err.Error(),
		}
	}

	// 翻译错误
	message := h.translator.Translate(domainErr.Code, locale, domainErr.Params)
	if message == string(domainErr.Code) {
		// 没有找到翻译，使用默认消息
		message = domainErr.Message
	}

	return message
}

// localeFromContext 从上下文中获取语言设置
func (h *ErrorHandler) localeFromContext(ctx context.Context) string {
	if ctx == nil {
		return h.locale
	}

	// 这里可以从上下文中获取语言设置
	// 例如，你可能有一个专门的键用于存储语言设置
	if locale, ok := ctx.Value("locale").(string); ok && locale != "" {
		return locale
	}

	return h.locale
}

// SetDefaultLocale 设置默认语言
func (h *ErrorHandler) SetDefaultLocale(locale string) {
	h.locale = locale
}
