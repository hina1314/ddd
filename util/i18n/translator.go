package i18n

// Translator defines the interface for translation services.
type Translator interface {
	// T translates a key for a given locale, applying parameters if provided.
	T(key string, locale string, params map[string]interface{}) string

	// LoadTranslations loads translation data from a directory.
	LoadTranslations(dir string) error

	// SetFallbackLocale sets the fallback locale for translations.
	SetFallbackLocale(locale string)
}
