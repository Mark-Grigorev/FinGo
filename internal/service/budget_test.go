package service

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestBudgetList_Success(t *testing.T) {
	want := []domain.Budget{{ID: 1, UserID: 1, CategoryID: 2, Limit: 5000}}
	store := &mockStore{
		listBudgetsFn: func(_ context.Context, _ int64, _ time.Time) ([]domain.Budget, error) {
			return want, nil
		},
	}
	svc := NewBudget(store, slog.Default())
	got, err := svc.List(context.Background(), 1, time.Now())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, int64(1), got[0].ID)
}

func TestBudgetCreate_InvalidLimit(t *testing.T) {
	svc := NewBudget(&mockStore{}, slog.Default())
	for _, limit := range []float64{0, -1, -100} {
		_, err := svc.Create(context.Background(), 1, 1, time.Now(), limit)
		assert.ErrorIsf(t, err, domain.ErrInvalidInput, "limit=%v", limit)
	}
}

func TestBudgetCreate_Success(t *testing.T) {
	want := &domain.Budget{ID: 10, UserID: 1, CategoryID: 2, Limit: 3000}
	store := &mockStore{
		createBudgetFn: func(_ context.Context, _ *domain.Budget) (*domain.Budget, error) {
			return want, nil
		},
	}
	svc := NewBudget(store, slog.Default())
	got, err := svc.Create(context.Background(), 1, 2, time.Now(), 3000)
	require.NoError(t, err)
	assert.Equal(t, int64(10), got.ID)
}

func TestBudgetUpdate_InvalidLimit(t *testing.T) {
	svc := NewBudget(&mockStore{}, slog.Default())
	_, err := svc.Update(context.Background(), 1, 1, 0)
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestBudgetUpdate_Success(t *testing.T) {
	want := &domain.Budget{ID: 5, Limit: 8000}
	store := &mockStore{
		updateBudgetFn: func(_ context.Context, _, _ int64, _ float64) (*domain.Budget, error) {
			return want, nil
		},
	}
	svc := NewBudget(store, slog.Default())
	got, err := svc.Update(context.Background(), 5, 1, 8000)
	require.NoError(t, err)
	assert.Equal(t, int64(5), got.ID)
}

func TestBudgetDelete_Success(t *testing.T) {
	called := false
	store := &mockStore{
		deleteBudgetFn: func(_ context.Context, _, _ int64) error {
			called = true
			return nil
		},
	}
	svc := NewBudget(store, slog.Default())
	require.NoError(t, svc.Delete(context.Background(), 1, 1))
	assert.True(t, called, "DeleteBudget was not called")
}
