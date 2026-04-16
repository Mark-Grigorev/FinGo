package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestTransactionDelete_BadID(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/abc", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
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

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestTransactionDelete_Success(t *testing.T) {
	env := newTestEnv()
	env.store.deleteTransactionFn = func(_ context.Context, _, _ int64) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/1", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
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

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	ct := w.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/csv") {
		t.Errorf("Content-Type = %q, want text/csv prefix", ct)
	}
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

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}
