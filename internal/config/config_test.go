package config_test

import (
	"testing"

	"github.com/Mark-Grigorev/FinGo/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigOK(t *testing.T) {

	appEnv := "local"
	tokenDur := "24h"
	dbConStr := "postgres://fingo:fingo_secret@localhost:5433/fingo?sslmode=disable"
	tokenSymKey := "ssssWWW"

	t.Setenv("APP_PORT", "8000")
	t.Setenv("APP_DEBUG", "true")
	t.Setenv("APP_ENV", appEnv)
	t.Setenv("TOKEN_DURATION", tokenDur)
	t.Setenv("DB_CONN_STRING", dbConStr)
	t.Setenv("TOKEN_SYMMETRIC_KEY", tokenSymKey)

	cfg, err := config.Load()
	assert.NoError(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, cfg.App.Debug, true)
	require.Equal(t, cfg.App.Port, 8000)
	require.Equal(t, cfg.App.Env, appEnv)
	require.Equal(t, cfg.DB, dbConStr)
	require.Equal(t, cfg.Token.Duration, tokenDur)
	require.Equal(t, cfg.Token.SymmetricKey, tokenSymKey)
}

func TestLoadConfigErr(t *testing.T) {

	cfg, err := config.Load()
	require.Nil(t, cfg)
	assert.EqualError(t, err, "DB_CONN_STRING is required")

}
