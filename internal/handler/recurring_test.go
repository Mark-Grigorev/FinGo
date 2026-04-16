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

func TestRecurringList_Success(t *testing.T) {
	env := newTestEnv()
	env.store.listRecurringFn = func(_ context.Context, _ int64) ([]domain.RecurringPayment, error) {
		return []domain.RecurringPayment{{ID: 1, Name: "Netflix", Amount: 500}}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/recurring", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecurringCreate_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPost, "/api/recurring", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecurringCreate_Success(t *testing.T) {
	env := newTestEnv()
	env.store.createRecurringFn = func(_ context.Context, r *domain.RecurringPayment) (*domain.RecurringPayment, error) {
		r.ID = 7
		return r, nil
	}

	body := strings.NewReader(`{"account_id":1,"name":"Netflix","amount":500,"frequency":"monthly"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/recurring", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRecurringUpdate_BadID(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/recurring/abc", strings.NewReader(`{"account_id":1,"name":"Test","amount":100}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecurringUpdate_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/recurring/1", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecurringUpdate_Success(t *testing.T) {
	env := newTestEnv()
	env.store.updateRecurringFn = func(_ context.Context, id, _ int64, _ string, _ float64, _ string, _ time.Time, _ int64, _ *int64) (*domain.RecurringPayment, error) {
		return &domain.RecurringPayment{ID: id, Name: "Gym"}, nil
	}

	body := strings.NewReader(`{"account_id":1,"name":"Gym","amount":1500,"frequency":"monthly"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/recurring/3", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecurringDelete_BadID(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodDelete, "/api/recurring/abc", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecurringDelete_NotFound(t *testing.T) {
	env := newTestEnv()
	env.store.deleteRecurringFn = func(_ context.Context, _, _ int64) error {
		return domain.ErrNotFound
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/recurring/99", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRecurringDelete_Success(t *testing.T) {
	env := newTestEnv()
	env.store.deleteRecurringFn = func(_ context.Context, _, _ int64) error { return nil }

	req := httptest.NewRequest(http.MethodDelete, "/api/recurring/1", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}
