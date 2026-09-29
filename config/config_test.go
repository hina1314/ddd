package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testTokenKey = "01234567890123456789012345678901"

func setRequiredEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("DB_SOURCE", "postgres://postgres:password@localhost:5432/app?sslmode=disable")
	t.Setenv("TOKEN_SYMMETRIC_KEY", testTokenKey)
	t.Setenv("ALLOWED_ORIGINS", "https://admin.example.com, http://localhost:5173")
}

func TestLoadConfigFromEnvironment(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("DB_MAX_OPEN_CONNS", "40")
	t.Setenv("SERVER_TRUSTED_PROXIES", "192.0.2.10, 172.18.0.1/32, ::1")

	cfg, err := LoadConfig(t.TempDir())
	require.NoError(t, err)
	require.Equal(t, "pgx", cfg.DBDriver)
	require.Equal(t, 40, cfg.DBMaxOpenConns)
	require.Equal(t, 10*time.Second, cfg.ServerReadTimeout)
	require.Equal(t, []string{"https://admin.example.com", "http://localhost:5173"}, cfg.AllowedOrigins)
	require.Equal(t, []string{"en", "zh"}, cfg.SupportedLocales)
	require.Equal(t, []string{"192.0.2.10", "172.18.0.1/32", "::1"}, cfg.ServerTrustedProxies)
}

func TestLoadConfigRejectsInvalidTrustedProxy(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("SERVER_TRUSTED_PROXIES", "nginx.example.com")
	_, err := LoadConfig(t.TempDir())
	require.ErrorContains(t, err, "SERVER_TRUSTED_PROXIES")
}

func TestLoadConfigRejectsInvalidSecret(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("TOKEN_SYMMETRIC_KEY", "too-short")

	_, err := LoadConfig(t.TempDir())
	require.ErrorContains(t, err, "exactly 32 bytes")
}

func TestValidateRejectsUnknownDefaultLocale(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("DEFAULT_LOCALE", "fr")

	_, err := LoadConfig(t.TempDir())
	require.ErrorContains(t, err, "DEFAULT_LOCALE")
}
