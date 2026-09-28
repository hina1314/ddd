package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const tokenKeySize = 32

// Config contains the runtime settings shared by the template infrastructure.
type Config struct {
	AppName             string        `mapstructure:"APP_NAME"`
	Environment         string        `mapstructure:"ENVIRONMENT"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	DBMaxOpenConns      int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns      int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBConnMaxLifetime   time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`
	DBConnMaxIdleTime   time.Duration `mapstructure:"DB_CONN_MAX_IDLE_TIME"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	ServerReadTimeout   time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout  time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
	ServerIdleTimeout   time.Duration `mapstructure:"SERVER_IDLE_TIMEOUT"`
	ServerBodyLimit     int           `mapstructure:"SERVER_BODY_LIMIT"`
	AllowedOrigins      []string      `mapstructure:"ALLOWED_ORIGINS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenIssuer         string        `mapstructure:"TOKEN_ISSUER"`
	TokenAudience       string        `mapstructure:"TOKEN_AUDIENCE"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	DefaultLocale       string        `mapstructure:"DEFAULT_LOCALE"`
	SupportedLocales    []string      `mapstructure:"SUPPORTED_LOCALES"`
	AuthRateLimitMax    int           `mapstructure:"AUTH_RATE_LIMIT_MAX"`
	AuthRateLimitWindow time.Duration `mapstructure:"AUTH_RATE_LIMIT_WINDOW"`
	MetricsToken        string        `mapstructure:"METRICS_TOKEN"`
	Debug               bool          `mapstructure:"DEBUG"`
	TrafficDrainDelay   time.Duration `mapstructure:"TRAFFIC_DRAIN_DELAY"`
	ShutdownTimeout     time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
	HealthCheckInterval time.Duration `mapstructure:"HEALTH_CHECK_INTERVAL"`
	HealthCheckTimeout  time.Duration `mapstructure:"HEALTH_CHECK_TIMEOUT"`
}

func LoadConfig(path string) (Config, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("app")
	v.SetConfigType("env")
	v.AutomaticEnv()
	setDefaults(v)
	for _, key := range configKeys() {
		if err := v.BindEnv(key); err != nil {
			return Config{}, fmt.Errorf("bind environment variable %s: %w", key, err)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return Config{}, fmt.Errorf("read configuration: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode configuration: %w", err)
	}
	cfg.AllowedOrigins = splitCSV(v.GetString("ALLOWED_ORIGINS"))
	cfg.SupportedLocales = splitCSV(v.GetString("SUPPORTED_LOCALES"))
	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("invalid configuration: %w", err)
	}
	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "go-fiber-template")
	v.SetDefault("ENVIRONMENT", "development")
	v.SetDefault("DB_DRIVER", "pgx")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 10)
	v.SetDefault("DB_CONN_MAX_LIFETIME", "30m")
	v.SetDefault("DB_CONN_MAX_IDLE_TIME", "5m")
	v.SetDefault("SERVER_ADDRESS", "0.0.0.0:3000")
	v.SetDefault("SERVER_READ_TIMEOUT", "10s")
	v.SetDefault("SERVER_WRITE_TIMEOUT", "15s")
	v.SetDefault("SERVER_IDLE_TIMEOUT", "60s")
	v.SetDefault("SERVER_BODY_LIMIT", 1<<20)
	v.SetDefault("TOKEN_ISSUER", "go-fiber-template")
	v.SetDefault("TOKEN_AUDIENCE", "go-fiber-template-api")
	v.SetDefault("ACCESS_TOKEN_DURATION", "30m")
	v.SetDefault("DEFAULT_LOCALE", "en")
	v.SetDefault("SUPPORTED_LOCALES", "en,zh")
	v.SetDefault("AUTH_RATE_LIMIT_MAX", 10)
	v.SetDefault("AUTH_RATE_LIMIT_WINDOW", "1m")
	v.SetDefault("TRAFFIC_DRAIN_DELAY", "5s")
	v.SetDefault("SHUTDOWN_TIMEOUT", "15s")
	v.SetDefault("HEALTH_CHECK_INTERVAL", "10s")
	v.SetDefault("HEALTH_CHECK_TIMEOUT", "1s")
}

func configKeys() []string {
	return []string{
		"APP_NAME", "ENVIRONMENT", "DB_DRIVER", "DB_SOURCE",
		"DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS", "DB_CONN_MAX_LIFETIME", "DB_CONN_MAX_IDLE_TIME",
		"SERVER_ADDRESS", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT", "SERVER_IDLE_TIMEOUT", "SERVER_BODY_LIMIT",
		"ALLOWED_ORIGINS", "TOKEN_SYMMETRIC_KEY", "TOKEN_ISSUER", "TOKEN_AUDIENCE", "ACCESS_TOKEN_DURATION",
		"DEFAULT_LOCALE", "SUPPORTED_LOCALES", "AUTH_RATE_LIMIT_MAX", "AUTH_RATE_LIMIT_WINDOW",
		"METRICS_TOKEN", "DEBUG", "TRAFFIC_DRAIN_DELAY", "SHUTDOWN_TIMEOUT",
		"HEALTH_CHECK_INTERVAL", "HEALTH_CHECK_TIMEOUT",
	}
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if value := strings.TrimSpace(part); value != "" {
			result = append(result, value)
		}
	}
	return result
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.AppName) == "" {
		return errors.New("APP_NAME must not be empty")
	}
	if c.DBDriver != "pgx" {
		return errors.New("DB_DRIVER must be pgx")
	}
	if strings.TrimSpace(c.DBSource) == "" || strings.TrimSpace(c.ServerAddress) == "" {
		return errors.New("DB_SOURCE and SERVER_ADDRESS must not be empty")
	}
	if c.DBMaxOpenConns <= 0 || c.DBMaxIdleConns < 0 || c.DBMaxIdleConns > c.DBMaxOpenConns {
		return errors.New("database pool sizes are invalid")
	}
	if c.DBConnMaxLifetime <= 0 || c.DBConnMaxIdleTime <= 0 {
		return errors.New("database connection lifetimes must be greater than zero")
	}
	if c.ServerReadTimeout <= 0 || c.ServerWriteTimeout <= 0 || c.ServerIdleTimeout <= 0 {
		return errors.New("server timeouts must be greater than zero")
	}
	if c.ServerBodyLimit < 1024 {
		return errors.New("SERVER_BODY_LIMIT must be at least 1024 bytes")
	}
	if len(c.TokenSymmetricKey) != tokenKeySize {
		return fmt.Errorf("TOKEN_SYMMETRIC_KEY must be exactly %d bytes", tokenKeySize)
	}
	if c.TokenIssuer == "" || c.TokenAudience == "" || c.AccessTokenDuration <= 0 {
		return errors.New("token issuer, audience and duration must be configured")
	}
	if len(c.SupportedLocales) == 0 || !contains(c.SupportedLocales, c.DefaultLocale) {
		return errors.New("DEFAULT_LOCALE must be included in SUPPORTED_LOCALES")
	}
	if c.AuthRateLimitMax <= 0 || c.AuthRateLimitWindow <= 0 {
		return errors.New("authentication rate limit must be greater than zero")
	}
	if c.TrafficDrainDelay < 0 || c.ShutdownTimeout <= 0 {
		return errors.New("shutdown durations are invalid")
	}
	if c.HealthCheckInterval <= 0 || c.HealthCheckTimeout <= 0 || c.HealthCheckTimeout > c.HealthCheckInterval {
		return errors.New("health check durations are invalid")
	}
	for _, origin := range c.AllowedOrigins {
		parsed, err := url.ParseRequestURI(origin)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			return fmt.Errorf("invalid allowed origin %q", origin)
		}
	}
	return nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
