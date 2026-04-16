package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestAuthLogin_BadJSON(t *testing.T) {
	env := newTestEnv()
	body := strings.NewReader(`{"email": "not-an-email"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthLogin_WrongCredentials(t *testing.T) {
	env := newTestEnv()
	env.store.getUserByEmailFn = func(_ context.Context, _ string) (*domain.User, error) {
		return nil, domain.ErrNotFound
	}

	body := strings.NewReader(`{"email":"user@example.com","password":"wrong"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthLogin_Success(t *testing.T) {
	env := newTestEnv()
	env.store.getUserByEmailFn = func(_ context.Context, _ string) (*domain.User, error) {
		return &domain.User{
			ID:           1,
			Email:        "user@example.com",
			PasswordHash: makeHash("secret123"),
		}, nil
	}

	body := strings.NewReader(`{"email":"user@example.com","password":"secret123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Set-Cookie"), "session_token")
}

func TestAuthRegister_BadJSON(t *testing.T) {
	env := newTestEnv()
	body := strings.NewReader(`{"email":"not-an-email","name":"","password":"short"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthRegister_Success(t *testing.T) {
	env := newTestEnv()
	env.store.createUserFn = func(_ context.Context, email, name, _ string) (*domain.User, error) {
		return &domain.User{ID: 2, Email: email, Name: name}, nil
	}

	body := strings.NewReader(`{"email":"new@example.com","name":"Alice","password":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Header().Get("Set-Cookie"), "session_token")
}

func TestAuthLogout(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestAuthMe_Success(t *testing.T) {
	env := newTestEnv()
	env.store.getUserByIDFn = func(_ context.Context, id int64) (*domain.User, error) {
		return &domain.User{ID: id, Email: "test@example.com", Name: "Test"}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
