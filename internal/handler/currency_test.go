package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCurrencyUpsertRate_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/currencies/rates/USD", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCurrencyUpsertRate_InvalidCurrency(t *testing.T) {
	env := newTestEnv()
	body := strings.NewReader(`{"rate":1.5}`)
	req := httptest.NewRequest(http.MethodPut, "/api/currencies/rates/XYZ", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
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

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestCurrencyDeleteRate_Success(t *testing.T) {
	env := newTestEnv()
	env.store.deleteExchangeRateFn = func(_ context.Context, _ int64, _ string) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/currencies/rates/USD", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
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

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCurrencySetBase_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/currencies/base", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCurrencySetBase_InvalidCurrency(t *testing.T) {
	env := newTestEnv()
	body := strings.NewReader(`{"currency":"XYZ"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/currencies/base", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCurrencySetBase_Success(t *testing.T) {
	env := newTestEnv()
	env.store.setBaseCurrencyFn = func(_ context.Context, _ int64, _ string) error {
		return nil
	}

	body := strings.NewReader(`{"currency":"EUR"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/currencies/base", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}
