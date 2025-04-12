package i18n

import (
	"context"
	"sync"
)

// TranslationService 提供翻译服务
type TranslationService struct {
	translator    Translator
	defaultLocale string
	mu            sync.RWMutex
}

// NewTranslationService 创建一个新的翻译服务
func NewTranslationService(translator Translator, defaultLocale string) *TranslationService {
	return &TranslationService{
		translator:    translator,
		defaultLocale: defaultLocale,
	}
}

// T 翻译一个键，使用上下文中的语言环境
func (s *TranslationService) T(ctx context.Context, key string, params map[string]interface{}) string {
	locale := LocaleFromContext(ctx, s.defaultLocale)
	return s.translator.T(key, locale, params)
}

// SetDefaultLocale 设置默认语言环境
func (s *TranslationService) SetDefaultLocale(locale string) {
	s.mu.Lock()
	s.defaultLocale = locale
	s.mu.Unlock()
}

// 全局翻译函数，类似ThinkPHP的lang()
var defaultService *TranslationService

// InitTranslator 初始化全局翻译器
func InitTranslator(translator Translator, defaultLocale string) {
	defaultService = NewTranslationService(translator, defaultLocale)
}

// T 全局翻译函数，类似ThinkPHP的lang()
func T(ctx context.Context, key string, params ...map[string]interface{}) string {
	if defaultService == nil {
		return key
	}

	var p map[string]interface{}
	if len(params) > 0 {
		p = params[0]
	}

	return defaultService.T(ctx, key, p)
}
