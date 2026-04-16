package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_RandomKey(t *testing.T) {
	m, err := New("", time.Hour)
	require.NoError(t, err)
	require.NotNil(t, m)
}

func TestNew_ValidHexKey(t *testing.T) {
	key := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	m, err := New(key, time.Hour)
	require.NoError(t, err)
	require.NotNil(t, m)
}

func TestNew_InvalidHexKey(t *testing.T) {
	_, err := New("not-valid-hex", time.Hour)
	require.Error(t, err)
}

func TestNew_ZeroDuration_UsesDefault(t *testing.T) {
	m, err := New("", 0)
	require.NoError(t, err)
	assert.Equal(t, defaultDuration, m.duration)
}

func TestCreateAndVerifyToken(t *testing.T) {
	m, _ := New("", time.Hour)

	tokenStr, payload, err := m.CreateToken(42, "user@example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)
	assert.Equal(t, int64(42), payload.UserID)
	assert.Equal(t, "user@example.com", payload.Email)

	verified, err := m.VerifyToken(tokenStr)
	require.NoError(t, err)
	assert.Equal(t, int64(42), verified.UserID)
	assert.Equal(t, "user@example.com", verified.Email)
}

func TestVerifyToken_Expired(t *testing.T) {
	m, _ := New("", -time.Second)
	tokenStr, _, err := m.CreateToken(1, "user@example.com")
	require.NoError(t, err)

	_, err = m.VerifyToken(tokenStr)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerifyToken_InvalidString(t *testing.T) {
	m, _ := New("", time.Hour)
	_, err := m.VerifyToken("not-a-valid-token")
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerifyToken_WrongKey(t *testing.T) {
	m1, _ := New("", time.Hour)
	m2, _ := New("", time.Hour)
	tokenStr, _, _ := m1.CreateToken(1, "user@example.com")

	_, err := m2.VerifyToken(tokenStr)
	require.ErrorIs(t, err, ErrInvalidToken)
}
