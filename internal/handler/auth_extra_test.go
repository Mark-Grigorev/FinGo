package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestAuthUpdateProfile_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/user/profile", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthUpdateProfile_Success(t *testing.T) {
	env := newTestEnv()
	env.store.updateUserFn = func(_ context.Context, _ int64, name, email string) (*domain.User, error) {
		return &domain.User{ID: 1, Name: name, Email: email}, nil
	}

	body := strings.NewReader(`{"name":"Alice","email":"alice@example.com"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/user/profile", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthChangePassword_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/user/password", strings.NewReader(`{"old_password":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthChangePassword_Success(t *testing.T) {
	env := newTestEnv()
	env.store.getUserByIDFn = func(_ context.Context, _ int64) (*domain.User, error) {
		return &domain.User{ID: 1, PasswordHash: makeHash("oldpassword")}, nil
	}
	env.store.updatePasswordFn = func(_ context.Context, _ int64, _ string) error { return nil }

	body := strings.NewReader(`{"old_password":"oldpassword","new_password":"newpassword"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/user/password", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestAuthForgotPassword_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthForgotPassword_Success(t *testing.T) {
	env := newTestEnv()
	env.store.getUserByEmailFn = func(_ context.Context, _ string) (*domain.User, error) {
		return &domain.User{ID: 1, Email: "user@example.com"}, nil
	}
	env.store.createPasswordResetFn = func(_ context.Context, _ string, _ int64, _ time.Time) error {
		return nil
	}

	body := strings.NewReader(`{"email":"user@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthResetPassword_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/reset-password", strings.NewReader(`{"token":"abc"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthResetPassword_Success(t *testing.T) {
	env := newTestEnv()
	env.store.getPasswordResetFn = func(_ context.Context, _ string) (*domain.PasswordReset, error) {
		return &domain.PasswordReset{
			Token:     "validtoken",
			UserID:    1,
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}, nil
	}
	env.store.updatePasswordFn = func(_ context.Context, _ int64, _ string) error { return nil }
	env.store.markPasswordResetUsedFn = func(_ context.Context, _ string) error { return nil }

	body := strings.NewReader(`{"token":"validtoken","new_password":"newpassword"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/reset-password", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}
