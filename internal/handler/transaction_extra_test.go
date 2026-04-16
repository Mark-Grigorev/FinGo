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

func TestTransactionDelete_BadID(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/abc", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionDelete_NotFound(t *testing.T) {
	env := newTestEnv()
	env.store.deleteTransactionFn = func(_ context.Context, _, _ int64) error {
		return domain.ErrNotFound
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/99", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTransactionDelete_Success(t *testing.T) {
	env := newTestEnv()
	env.store.deleteTransactionFn = func(_ context.Context, _, _ int64) error { return nil }

	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/1", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestTransactionExport_Success(t *testing.T) {
	env := newTestEnv()
	env.store.exportTransactionsFn = func(_ context.Context, _ int64, _, _ time.Time) ([]domain.Transaction, error) {
		return []domain.Transaction{
			{ID: 1, Type: "income", Amount: 1000, Name: "Salary", Date: time.Now()},
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/export?from=2024-01-01&to=2024-01-31", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.HasPrefix(w.Header().Get("Content-Type"), "text/csv"))
}

func TestTransactionExport_DefaultDates(t *testing.T) {
	env := newTestEnv()
	env.store.exportTransactionsFn = func(_ context.Context, _ int64, _, _ time.Time) ([]domain.Transaction, error) {
		return []domain.Transaction{}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/export", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
