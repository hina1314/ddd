package i18n

import (
	"context"
)

type contextKey string

const (
	localeKey contextKey = "locale"
)

// WithLocale 向上下文添加语言环境
func WithLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, localeKey, locale)
}

// LocaleFromContext 从上下文获取语言环境
func LocaleFromContext(ctx context.Context, defaultLocale string) string {
	if ctx == nil {
		return defaultLocale
	}

	if locale, ok := ctx.Value(localeKey).(string); ok && locale != "" {
		return locale
	}

	return defaultLocale
}
