package i18n

import (
	"context"
	"sync"
)

// TranslationService provides a high-level interface for translation operations.
type TranslationService struct {
	translator    Translator
	defaultLocale string
	mu            sync.RWMutex
}

// NewTranslationService creates a new TranslationService with the given translator and default locale.
func NewTranslationService(translator Translator, defaultLocale string) *TranslationService {
	return &TranslationService{
		translator:    translator,
		defaultLocale: defaultLocale,
	}
}

// T translates a key using the locale from the context, applying parameters if provided.
func (s *TranslationService) T(ctx context.Context, key string, params map[string]interface{}) string {
	s.mu.RLock()
	defaultLocale := s.defaultLocale
	s.mu.RUnlock()
	locale := LocaleFromContext(ctx, defaultLocale)
	return s.translator.T(key, locale, params)
}

// SetDefaultLocale sets the default locale for translations.
func (s *TranslationService) SetDefaultLocale(locale string) {
	s.mu.Lock()
	s.defaultLocale = locale
	s.mu.Unlock()
}

// LoadTranslations loads translation files using the underlying translator.
func (s *TranslationService) LoadTranslations(dir string) error {
	return s.translator.LoadTranslations(dir)
}
