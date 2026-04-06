package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig
	DB    string
	Token TokenConfig
	Email EmailConfig
}

type EmailConfig struct {
	ResendAPIKey string
	BaseURL      string
}

type AppConfig struct {
	Port  int
	Debug bool
	Env   string
}

type TokenConfig struct {
	SymmetricKey string
	Duration     string
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("APP_PORT", 8008)
	v.SetDefault("APP_DEBUG", true)
	v.SetDefault("APP_ENV", "local")
	v.SetDefault("TOKEN_DURATION", "24h")

	cfg := &Config{
		App: AppConfig{
			Port:  v.GetInt("APP_PORT"),
			Debug: v.GetBool("APP_DEBUG"),
			Env:   v.GetString("APP_ENV"),
		},
		DB: v.GetString("DB_CONN_STRING"),
		Token: TokenConfig{
			SymmetricKey: v.GetString("TOKEN_SYMMETRIC_KEY"),
			Duration:     v.GetString("TOKEN_DURATION"),
		},
		Email: EmailConfig{
			ResendAPIKey: v.GetString("RESEND_API_KEY"),
			BaseURL:      v.GetString("APP_BASE_URL"),
		},
	}

	if cfg.DB == "" {
		return nil, fmt.Errorf("DB_CONN_STRING is required")
	}
	if cfg.Email.ResendAPIKey == "" {
		return nil, fmt.Errorf("RESEND_API_KEY is required")
	}

	return cfg, nil
}
