package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hina1314/ddd/util/i18n"
	"strings"
)

func Locale(defaultLocale string, supportedLocales []string) fiber.Handler {
	supported := make(map[string]string, len(supportedLocales))
	for _, locale := range supportedLocales {
		normalized := strings.ToLower(strings.TrimSpace(locale))
		supported[normalized] = locale
	}

	return func(c fiber.Ctx) error {
		locale := selectLocale(c.Get("Accept-Language"), defaultLocale, supported)
		ctx := i18n.WithLocale(c.Context(), locale)
		c.SetContext(ctx)
		return c.Next()
	}
}

func selectLocale(header, fallback string, supported map[string]string) string {
	for _, candidate := range strings.Split(header, ",") {
		candidate = strings.TrimSpace(strings.SplitN(candidate, ";", 2)[0])
		candidate = strings.ToLower(strings.ReplaceAll(candidate, "_", "-"))
		if locale, ok := supported[candidate]; ok {
			return locale
		}
		if base := strings.SplitN(candidate, "-", 2)[0]; base != candidate {
			if locale, ok := supported[base]; ok {
				return locale
			}
		}
	}
	return fallback
}
