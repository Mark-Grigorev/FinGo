package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestAuthLogin_BadJSON(t *testing.T) {
	env := newTestEnv()
	body := strings.NewReader(`{"email": "not-an-email"}`) // missing password, invalid email
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
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

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
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

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	cookie := w.Header().Get("Set-Cookie")
	if !strings.Contains(cookie, "session_token") {
		t.Errorf("expected session_token cookie, got: %q", cookie)
	}
}

func TestAuthRegister_BadJSON(t *testing.T) {
	env := newTestEnv()
	body := strings.NewReader(`{"email":"not-an-email","name":"","password":"short"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
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

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
	cookie := w.Header().Get("Set-Cookie")
	if !strings.Contains(cookie, "session_token") {
		t.Errorf("expected session_token cookie, got: %q", cookie)
	}
}

func TestAuthLogout(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
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

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}
