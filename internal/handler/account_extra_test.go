package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestAccountUpdate_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/accounts/1", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAccountUpdate_EmptyName(t *testing.T) {
	env := newTestEnv()
	body := strings.NewReader(`{"name":"  ","type":"card","currency":"RUB"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/accounts/1", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAccountUpdate_Success(t *testing.T) {
	env := newTestEnv()
	env.store.updateAccountFn = func(_ context.Context, id, _ int64, name, typ, currency string) (*domain.Account, error) {
		return &domain.Account{ID: id, Name: name, Type: typ, Currency: currency}, nil
	}

	body := strings.NewReader(`{"name":"Wallet","type":"cash","currency":"RUB"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/accounts/1", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
