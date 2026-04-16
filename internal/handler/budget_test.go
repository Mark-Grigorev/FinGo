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

func TestBudgetList_Success(t *testing.T) {
	env := newTestEnv()
	env.store.listBudgetsFn = func(_ context.Context, _ int64, _ time.Time) ([]domain.Budget, error) {
		return []domain.Budget{{ID: 1, CategoryID: 1, Limit: 5000}}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/budgets", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestBudgetList_WithMonth(t *testing.T) {
	env := newTestEnv()
	env.store.listBudgetsFn = func(_ context.Context, _ int64, _ time.Time) ([]domain.Budget, error) {
		return []domain.Budget{}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/budgets?month=2024-01", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestBudgetCreate_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPost, "/api/budgets", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestBudgetCreate_Success(t *testing.T) {
	env := newTestEnv()
	env.store.createBudgetFn = func(_ context.Context, b *domain.Budget) (*domain.Budget, error) {
		b.ID = 1
		return b, nil
	}

	body := strings.NewReader(`{"category_id":1,"limit":5000}`)
	req := httptest.NewRequest(http.MethodPost, "/api/budgets", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestBudgetUpdate_BadID(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/budgets/abc", strings.NewReader(`{"limit":1000}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestBudgetUpdate_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/budgets/1", strings.NewReader(`{"limit":"not a number"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestBudgetUpdate_Success(t *testing.T) {
	env := newTestEnv()
	env.store.updateBudgetFn = func(_ context.Context, id, _ int64, limit float64) (*domain.Budget, error) {
		return &domain.Budget{ID: id, Limit: limit}, nil
	}

	body := strings.NewReader(`{"limit":8000}`)
	req := httptest.NewRequest(http.MethodPut, "/api/budgets/1", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestBudgetDelete_BadID(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodDelete, "/api/budgets/abc", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestBudgetDelete_NotFound(t *testing.T) {
	env := newTestEnv()
	env.store.deleteBudgetFn = func(_ context.Context, _, _ int64) error {
		return domain.ErrNotFound
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/budgets/99", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestBudgetDelete_Success(t *testing.T) {
	env := newTestEnv()
	env.store.deleteBudgetFn = func(_ context.Context, _, _ int64) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/budgets/1", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}
