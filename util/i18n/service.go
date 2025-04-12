package i18n

import (
	"context"
	"fmt"
	"sync"
)

// TranslationService provides a high-level interface for translation operations.
type TranslationService struct {
	translator    Translator
	defaultLocale string
	cache         *sync.Map // Optional cache for translated messages
	mu            sync.RWMutex
}

// NewTranslationService creates a new TranslationService with the given translator and default locale.
func NewTranslationService(translator Translator, defaultLocale string) *TranslationService {
	return &TranslationService{
		translator:    translator,
		defaultLocale: defaultLocale,
		cache:         &sync.Map{},
	}
}

// T translates a key using the locale from the context, applying parameters if provided.
func (s *TranslationService) T(ctx context.Context, key string, params map[string]interface{}) string {
	locale := LocaleFromContext(ctx, s.defaultLocale)

	// Check cache
	cacheKey := fmt.Sprintf("%s:%s:%v", locale, key, params)
	if cached, ok := s.cache.Load(cacheKey); ok {
		return cached.(string)
	}

	// Translate
	translated := s.translator.T(key, locale, params)

	// Store in cache
	s.cache.Store(cacheKey, translated)
	return translated
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
