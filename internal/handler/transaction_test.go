package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/internal/repository"
)

func TestTransactionList_Success(t *testing.T) {
	env := newTestEnv()
	env.store.listTransactionsFn = func(_ context.Context, _ int64, _ repository.TransactionFilter) ([]domain.Transaction, int, error) {
		return []domain.Transaction{}, 0, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/transactions?page=1&limit=10", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestTransactionCreate_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestTransactionCreate_InvalidType(t *testing.T) {
	env := newTestEnv()
	body := strings.NewReader(`{"account_id":1,"type":"transfer","amount":100}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestTransactionCreate_Success(t *testing.T) {
	env := newTestEnv()
	env.store.createTransactionFn = func(_ context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
		tx.ID = 10
		return tx, nil
	}

	body := strings.NewReader(`{"account_id":1,"type":"income","amount":500,"name":"Salary","category_id":1}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
}
