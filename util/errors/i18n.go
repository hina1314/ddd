package errors

import (
	"fmt"
	"strings"
	"sync"
)

// Translator 定义了一个翻译错误消息的接口
type Translator interface {
	Translate(code ErrorCode, locale string, params map[string]interface{}) string
}

// SimpleTranslator 是Translator接口的简单实现
type SimpleTranslator struct {
	mu       sync.RWMutex
	messages map[string]map[ErrorCode]string // locale -> code -> message
}

// NewSimpleTranslator 创建一个新的简单翻译器
func NewSimpleTranslator() *SimpleTranslator {
	return &SimpleTranslator{
		messages: make(map[string]map[ErrorCode]string),
	}
}

// RegisterTranslation 注册一个错误代码的翻译
func (t *SimpleTranslator) RegisterTranslation(code ErrorCode, locale, message string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.messages[locale]; !ok {
		t.messages[locale] = make(map[ErrorCode]string)
	}
	t.messages[locale][code] = message
}

// RegisterTranslations 批量注册错误代码的翻译
func (t *SimpleTranslator) RegisterTranslations(locale string, translations map[ErrorCode]string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.messages[locale]; !ok {
		t.messages[locale] = make(map[ErrorCode]string)
	}

	for code, message := range translations {
		t.messages[locale][code] = message
	}
}

// Translate 根据错误代码和语言环境翻译错误消息
func (t *SimpleTranslator) Translate(code ErrorCode, locale string, params map[string]interface{}) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// 尝试获取指定语言的翻译
	if localeMessages, ok := t.messages[locale]; ok {
		if message, ok := localeMessages[code]; ok {
			return t.interpolate(message, params)
		}
	}

	// 尝试获取默认语言的翻译
	if locale != "en" {
		if defaultMessages, ok := t.messages["en"]; ok {
			if message, ok := defaultMessages[code]; ok {
				return t.interpolate(message, params)
			}
		}
	}

	// 没有找到翻译，返回错误代码
	return string(code)
}

// interpolate 替换消息中的参数占位符
func (t *SimpleTranslator) interpolate(message string, params map[string]interface{}) string {
	if params == nil {
		return message
	}

	result := message
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		result = strings.Replace(result, placeholder, fmt.Sprintf("%v", value), -1)
	}
	return result
}
