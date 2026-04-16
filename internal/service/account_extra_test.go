package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestAccountList_Success(t *testing.T) {
	want := []domain.Account{{ID: 1, Name: "Cash"}, {ID: 2, Name: "Card"}}
	store := &mockStore{
		listAccountsFn: func(_ context.Context, _ int64) ([]domain.Account, error) {
			return want, nil
		},
	}
	svc := NewAccount(store, slog.Default())
	got, err := svc.List(context.Background(), 1)
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestAccountGet_Success(t *testing.T) {
	want := &domain.Account{ID: 5, Name: "Savings"}
	store := &mockStore{
		getAccountFn: func(_ context.Context, id, _ int64) (*domain.Account, error) {
			if id != 5 {
				return nil, domain.ErrNotFound
			}
			return want, nil
		},
	}
	svc := NewAccount(store, slog.Default())
	got, err := svc.Get(context.Background(), 5, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(5), got.ID)
}
