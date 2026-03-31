package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestAccountCreate_EmptyName(t *testing.T) {
	svc := NewAccount(&mockStore{}, slog.Default())
	_, err := svc.Create(context.Background(), 1, CreateAccountInput{Name: "   "})
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Type != "card" {
		t.Errorf("Type = %q, want %q", got.Type, "card")
	}
	if got.Currency != "RUB" {
		t.Errorf("Currency = %q, want %q", got.Currency, "RUB")
	}
}

func TestAccountCreate_Success(t *testing.T) {
	want := &domain.Account{ID: 10, UserID: 1, Name: "Salary", Type: "card", Currency: "USD"}
	store := &mockStore{
		createAccountFn: func(_ context.Context, a *domain.Account) (*domain.Account, error) {
			return want, nil
		},
	}
	svc := NewAccount(store, slog.Default())
	acc, err := svc.Create(context.Background(), 1, CreateAccountInput{Name: "Salary", Type: "card", Currency: "USD"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acc.ID != 10 {
		t.Errorf("acc.ID = %d, want 10", acc.ID)
	}
}

func TestAccountUpdate_EmptyName(t *testing.T) {
	svc := NewAccount(&mockStore{}, slog.Default())
	_, err := svc.Update(context.Background(), 1, 1, UpdateAccountInput{Name: ""})
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotName != "Wallet" {
		t.Errorf("name = %q, want %q", gotName, "Wallet")
	}
	if gotType != "card" {
		t.Errorf("type = %q, want %q", gotType, "card")
	}
	if gotCurrency != "RUB" {
		t.Errorf("currency = %q, want %q", gotCurrency, "RUB")
	}
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
	if err := svc.Delete(context.Background(), 5, 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("store.DeleteAccount was not called")
	}
}
