package service

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, 1, total)
}

func TestTransactionDelete_Success(t *testing.T) {
	called := false
	store := &mockStore{
		deleteTransactionFn: func(_ context.Context, _, _ int64) error {
			called = true
			return nil
		},
	}
	svc := NewTransaction(store, slog.Default())
	require.NoError(t, svc.Delete(context.Background(), 1, 1))
	assert.True(t, called, "DeleteTransaction was not called")
}

func TestTransactionDelete_NotFound(t *testing.T) {
	store := &mockStore{
		deleteTransactionFn: func(_ context.Context, _, _ int64) error {
			return domain.ErrNotFound
		},
	}
	svc := NewTransaction(store, slog.Default())
	err := svc.Delete(context.Background(), 99, 1)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTransactionExport_Success(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	want := []domain.Transaction{
		{ID: 1, Amount: 1000, Type: "income"},
		{ID: 2, Amount: 500, Type: "expense"},
	}
	store := &mockStore{
		exportTransactionsFn: func(_ context.Context, _ int64, _, _ time.Time) ([]domain.Transaction, error) {
			return want, nil
		},
	}
	svc := NewTransaction(store, slog.Default())
	got, err := svc.Export(context.Background(), 1, from, to)
	require.NoError(t, err)
	assert.Len(t, got, 2)
}
