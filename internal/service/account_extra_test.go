package service

import (
	"context"
	"log/slog"
	"testing"

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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("len = %d, want 2", len(got))
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 5 {
		t.Errorf("got.ID = %d, want 5", got.ID)
	}
}
