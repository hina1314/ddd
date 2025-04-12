package i18n

import "context"

// contextKey defines a key type for storing locale in context.
type contextKey string

const (
	localeKey contextKey = "locale"
)

// WithLocale adds a locale to the context.
func WithLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, localeKey, locale)
}

// LocaleFromContext retrieves the locale from the context, falling back to defaultLocale if not set.
func LocaleFromContext(ctx context.Context, defaultLocale string) string {
	if ctx == nil {
		return defaultLocale
	}
	if locale, ok := ctx.Value(localeKey).(string); ok && locale != "" {
		return locale
	}
	return defaultLocale
}
