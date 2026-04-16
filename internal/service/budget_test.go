package service

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestBudgetList_Success(t *testing.T) {
	want := []domain.Budget{{ID: 1, UserID: 1, CategoryID: 2, Limit: 5000}}
	store := &mockStore{
		listBudgetsFn: func(_ context.Context, userID int64, _ time.Time) ([]domain.Budget, error) {
			return want, nil
		},
	}
	svc := NewBudget(store, slog.Default())
	got, err := svc.List(context.Background(), 1, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != 1 {
		t.Errorf("unexpected list: %v", got)
	}
}

func TestBudgetCreate_InvalidLimit(t *testing.T) {
	svc := NewBudget(&mockStore{}, slog.Default())
	for _, limit := range []float64{0, -1, -100} {
		_, err := svc.Create(context.Background(), 1, 1, time.Now(), limit)
		if err != domain.ErrInvalidInput {
			t.Errorf("limit=%v: expected ErrInvalidInput, got %v", limit, err)
		}
	}
}

func TestBudgetCreate_Success(t *testing.T) {
	want := &domain.Budget{ID: 10, UserID: 1, CategoryID: 2, Limit: 3000}
	store := &mockStore{
		createBudgetFn: func(_ context.Context, b *domain.Budget) (*domain.Budget, error) {
			return want, nil
		},
	}
	svc := NewBudget(store, slog.Default())
	got, err := svc.Create(context.Background(), 1, 2, time.Now(), 3000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 10 {
		t.Errorf("got.ID = %d, want 10", got.ID)
	}
}

func TestBudgetUpdate_InvalidLimit(t *testing.T) {
	svc := NewBudget(&mockStore{}, slog.Default())
	_, err := svc.Update(context.Background(), 1, 1, 0)
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestBudgetUpdate_Success(t *testing.T) {
	want := &domain.Budget{ID: 5, Limit: 8000}
	store := &mockStore{
		updateBudgetFn: func(_ context.Context, id, _ int64, limit float64) (*domain.Budget, error) {
			return want, nil
		},
	}
	svc := NewBudget(store, slog.Default())
	got, err := svc.Update(context.Background(), 5, 1, 8000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 5 {
		t.Errorf("got.ID = %d, want 5", got.ID)
	}
}

func TestBudgetDelete_Success(t *testing.T) {
	called := false
	store := &mockStore{
		deleteBudgetFn: func(_ context.Context, id, userID int64) error {
			called = true
			return nil
		},
	}
	svc := NewBudget(store, slog.Default())
	if err := svc.Delete(context.Background(), 1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("DeleteBudget was not called")
	}
}
