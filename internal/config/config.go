package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config содержит всю конфигурацию приложения, загружаемую из env.
type Config struct {
	App AppConfig
	DB  DBConfig
}

type AppConfig struct {
	Port int    // APP_PORT (default: 8008)
	Env  string // APP_ENV: "local" | "dev" | "prod" (default: "local")
}

type DBConfig struct {
	Host     string // DB_HOST (default: "localhost")
	Port     int    // DB_PORT (default: 5432)
	User     string // DB_USER (required)
	Password string // DB_PASSWORD (required)
	Name     string // DB_NAME (required)
	SSLMode  string // DB_SSLMODE (default: "disable")
}

// DSN возвращает строку подключения к PostgreSQL.
func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// Load читает конфигурацию из переменных окружения.
// Возвращает ошибку если обязательные переменные не заданы.
func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Port: getInt("APP_PORT", 8008),
			Env:  getStr("APP_ENV", "local"),
		},
		DB: DBConfig{
			Host:     getStr("DB_HOST", "localhost"),
			Port:     getInt("DB_PORT", 5432),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  getStr("DB_SSLMODE", "disable"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	var missing []string

	if c.DB.User == "" {
		missing = append(missing, "DB_USER")
	}
	if c.DB.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.DB.Name == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return fmt.Errorf("required env vars not set: %s", strings.Join(missing, ", "))
	}

	validEnvs := map[string]bool{"local": true, "dev": true, "prod": true}
	if !validEnvs[c.App.Env] {
		return errors.New("APP_ENV must be one of: local, dev, prod")
	}

	if c.App.Port < 1 || c.App.Port > 65535 {
		return errors.New("APP_PORT must be between 1 and 65535")
	}

	return nil
}

func getStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
