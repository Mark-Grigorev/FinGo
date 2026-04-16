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

func TestAccountList_Success(t *testing.T) {
	env := newTestEnv()
	env.store.listAccountsFn = func(_ context.Context, _ int64) ([]domain.Account, error) {
		return []domain.Account{{ID: 1, Name: "Cash", Type: "cash", Currency: "RUB"}}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/accounts", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAccountCreate_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPost, "/api/accounts", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAccountCreate_EmptyName(t *testing.T) {
	env := newTestEnv()
	body := strings.NewReader(`{"name":"   "}`)
	req := httptest.NewRequest(http.MethodPost, "/api/accounts", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAccountCreate_Success(t *testing.T) {
	env := newTestEnv()
	env.store.createAccountFn = func(_ context.Context, a *domain.Account) (*domain.Account, error) {
		a.ID = 42
		return a, nil
	}

	body := strings.NewReader(`{"name":"Salary","type":"card","currency":"USD"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/accounts", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAccountDelete_NotFound(t *testing.T) {
	env := newTestEnv()
	env.store.deleteAccountFn = func(_ context.Context, _, _ int64) error {
		return domain.ErrNotFound
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/accounts/99", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAccountDelete_Success(t *testing.T) {
	env := newTestEnv()
	env.store.deleteAccountFn = func(_ context.Context, _, _ int64) error { return nil }

	req := httptest.NewRequest(http.MethodDelete, "/api/accounts/1", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}
