package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	AllowedOrigins      []string      `mapstructure:"ALLOWED_ORIGINS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	DefaultLocale       string        `mapstructure:"DEFAULT_LOCALE"`
	Debug               bool          `mapstructure:"DEBUG"`
	TrafficDrainDelay   time.Duration `mapstructure:"TRAFFIC_DRAIN_DELAY"`
	ShutdownTimeout     time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") //jso,xml

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	if err = config.Validate(); err != nil {
		err = fmt.Errorf("invalid configuration: %w", err)
		return
	}
	return
}

func (c Config) Validate() error {
	if c.TrafficDrainDelay < 0 {
		return fmt.Errorf(
			"TRAFFIC_DRAIN_DELAY must not be negative",
		)
	}

	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf(
			"SHUTDOWN_TIMEOUT must be greater than zero",
		)
	}

	return nil
}
