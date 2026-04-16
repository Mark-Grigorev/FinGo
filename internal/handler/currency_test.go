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

func TestCurrencyListRates_Success(t *testing.T) {
	env := newTestEnv()
	env.store.listExchangeRatesFn = func(_ context.Context, _ int64) ([]domain.ExchangeRate, error) {
		return []domain.ExchangeRate{{ID: 1, Currency: "USD", Rate: 90.5}}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/currencies/rates", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCurrencyUpsertRate_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/currencies/rates/USD", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCurrencyUpsertRate_InvalidCurrency(t *testing.T) {
	env := newTestEnv()
	body := strings.NewReader(`{"rate":1.5}`)
	req := httptest.NewRequest(http.MethodPut, "/api/currencies/rates/XYZ", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCurrencyUpsertRate_Success(t *testing.T) {
	env := newTestEnv()
	env.store.upsertExchangeRateFn = func(_ context.Context, _ int64, currency string, rate float64) (*domain.ExchangeRate, error) {
		return &domain.ExchangeRate{Currency: currency, Rate: rate}, nil
	}

	body := strings.NewReader(`{"rate":90.5}`)
	req := httptest.NewRequest(http.MethodPut, "/api/currencies/rates/USD", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCurrencyDeleteRate_Success(t *testing.T) {
	env := newTestEnv()
	env.store.deleteExchangeRateFn = func(_ context.Context, _ int64, _ string) error { return nil }

	req := httptest.NewRequest(http.MethodDelete, "/api/currencies/rates/USD", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestCurrencyGetBase_Success(t *testing.T) {
	env := newTestEnv()
	env.store.getBaseCurrencyFn = func(_ context.Context, _ int64) (string, error) {
		return "USD", nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/currencies/base", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCurrencySetBase_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/currencies/base", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCurrencySetBase_InvalidCurrency(t *testing.T) {
	env := newTestEnv()
	body := strings.NewReader(`{"currency":"XYZ"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/currencies/base", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCurrencySetBase_Success(t *testing.T) {
	env := newTestEnv()
	env.store.setBaseCurrencyFn = func(_ context.Context, _ int64, _ string) error { return nil }

	body := strings.NewReader(`{"currency":"EUR"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/currencies/base", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}
