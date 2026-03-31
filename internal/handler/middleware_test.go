package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestAuthMiddleware_NoToken(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/accounts", nil)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/accounts", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.value")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddleware_ValidBearerHeader(t *testing.T) {
	env := newTestEnv()
	env.store.listAccountsFn = func(_ context.Context, _ int64) ([]domain.Account, error) {
		return []domain.Account{}, nil
	}

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/accounts", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestAuthMiddleware_ValidCookie(t *testing.T) {
	env := newTestEnv()
	env.store.listAccountsFn = func(_ context.Context, _ int64) ([]domain.Account, error) {
		return []domain.Account{}, nil
	}

	tok, _, _ := env.maker.CreateToken(1, "test@example.com")
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/accounts", nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: tok})
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}
