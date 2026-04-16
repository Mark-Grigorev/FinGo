package service

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/internal/repository"
)

func TestTransactionList_Success(t *testing.T) {
	want := []domain.Transaction{{ID: 1, Amount: 500, Type: "income"}}
	store := &mockStore{
		listTransactionsFn: func(_ context.Context, _ int64, _ repository.TransactionFilter) ([]domain.Transaction, int, error) {
			return want, 1, nil
		},
	}
	svc := NewTransaction(store, slog.Default())
	got, total, err := svc.List(context.Background(), 1, repository.TransactionFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || total != 1 {
		t.Errorf("got len=%d total=%d, want len=1 total=1", len(got), total)
	}
}

func TestTransactionDelete_Success(t *testing.T) {
	called := false
	store := &mockStore{
		deleteTransactionFn: func(_ context.Context, id, userID int64) error {
			called = true
			return nil
		},
	}
	svc := NewTransaction(store, slog.Default())
	if err := svc.Delete(context.Background(), 1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("DeleteTransaction was not called")
	}
}

func TestTransactionDelete_NotFound(t *testing.T) {
	store := &mockStore{
		deleteTransactionFn: func(_ context.Context, id, userID int64) error {
			return domain.ErrNotFound
		},
	}
	svc := NewTransaction(store, slog.Default())
	err := svc.Delete(context.Background(), 99, 1)
	if err != domain.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestTransactionExport_Success(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	want := []domain.Transaction{
		{ID: 1, Amount: 1000, Type: "income"},
		{ID: 2, Amount: 500, Type: "expense"},
	}
	store := &mockStore{
		exportTransactionsFn: func(_ context.Context, _ int64, f, t time.Time) ([]domain.Transaction, error) {
			return want, nil
		},
	}
	svc := NewTransaction(store, slog.Default())
	got, err := svc.Export(context.Background(), 1, from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("len = %d, want 2", len(got))
	}
}
