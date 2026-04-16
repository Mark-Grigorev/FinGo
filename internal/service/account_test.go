package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestAccountCreate_EmptyName(t *testing.T) {
	svc := NewAccount(&mockStore{}, slog.Default())
	_, err := svc.Create(context.Background(), 1, CreateAccountInput{Name: "   "})
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestAccountCreate_Defaults(t *testing.T) {
	var got *domain.Account
	store := &mockStore{
		createAccountFn: func(_ context.Context, a *domain.Account) (*domain.Account, error) {
			got = a
			return a, nil
		},
	}
	svc := NewAccount(store, slog.Default())
	_, err := svc.Create(context.Background(), 1, CreateAccountInput{Name: "Cash"})
	require.NoError(t, err)
	assert.Equal(t, "card", got.Type)
	assert.Equal(t, "RUB", got.Currency)
}

func TestAccountCreate_Success(t *testing.T) {
	want := &domain.Account{ID: 10, UserID: 1, Name: "Salary", Type: "card", Currency: "USD"}
	store := &mockStore{
		createAccountFn: func(_ context.Context, _ *domain.Account) (*domain.Account, error) {
			return want, nil
		},
	}
	svc := NewAccount(store, slog.Default())
	acc, err := svc.Create(context.Background(), 1, CreateAccountInput{Name: "Salary", Type: "card", Currency: "USD"})
	require.NoError(t, err)
	assert.Equal(t, int64(10), acc.ID)
}

func TestAccountUpdate_EmptyName(t *testing.T) {
	svc := NewAccount(&mockStore{}, slog.Default())
	_, err := svc.Update(context.Background(), 1, 1, UpdateAccountInput{Name: ""})
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestAccountUpdate_Defaults(t *testing.T) {
	var gotName, gotType, gotCurrency string
	store := &mockStore{
		updateAccountFn: func(_ context.Context, _, _ int64, name, typ, currency string) (*domain.Account, error) {
			gotName, gotType, gotCurrency = name, typ, currency
			return &domain.Account{}, nil
		},
	}
	svc := NewAccount(store, slog.Default())
	_, err := svc.Update(context.Background(), 1, 1, UpdateAccountInput{Name: "Wallet"})
	require.NoError(t, err)
	assert.Equal(t, "Wallet", gotName)
	assert.Equal(t, "card", gotType)
	assert.Equal(t, "RUB", gotCurrency)
}

func TestAccountDelete_DelegatesToStore(t *testing.T) {
	called := false
	store := &mockStore{
		deleteAccountFn: func(_ context.Context, id, userID int64) error {
			called = true
			if id != 5 || userID != 2 {
				return domain.ErrNotFound
			}
			return nil
		},
	}
	svc := NewAccount(store, slog.Default())
	require.NoError(t, svc.Delete(context.Background(), 5, 2))
	assert.True(t, called, "store.DeleteAccount was not called")
}
