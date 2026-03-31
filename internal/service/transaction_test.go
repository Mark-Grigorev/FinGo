package service

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestTransactionCreate_InvalidType(t *testing.T) {
	svc := NewTransaction(&mockStore{}, slog.Default())
	cases := []string{"", "transfer", "INCOME"}
	for _, typ := range cases {
		_, err := svc.Create(context.Background(), 1, CreateTransactionInput{
			Type: typ, Amount: 100, AccountID: 1,
		})
		if err != domain.ErrInvalidInput {
			t.Errorf("type=%q: expected ErrInvalidInput, got %v", typ, err)
		}
	}
}

func TestTransactionCreate_InvalidAmount(t *testing.T) {
	svc := NewTransaction(&mockStore{}, slog.Default())
	cases := []float64{0, -1, -100}
	for _, amount := range cases {
		_, err := svc.Create(context.Background(), 1, CreateTransactionInput{
			Type: "income", Amount: amount, AccountID: 1,
		})
		if err != domain.ErrInvalidInput {
			t.Errorf("amount=%v: expected ErrInvalidInput, got %v", amount, err)
		}
	}
}

func TestTransactionCreate_ZeroAccountID(t *testing.T) {
	svc := NewTransaction(&mockStore{}, slog.Default())
	_, err := svc.Create(context.Background(), 1, CreateTransactionInput{
		Type: "income", Amount: 50, AccountID: 0,
	})
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestTransactionCreate_WithDate(t *testing.T) {
	var gotDate time.Time
	store := &mockStore{
		getAccountFn: func(_ context.Context, id, userID int64) (*domain.Account, error) {
			return &domain.Account{ID: id, UserID: userID, Balance: 1000, Type: "card"}, nil
		},
		createTransactionFn: func(_ context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
			gotDate = tx.Date
			return tx, nil
		},
	}
	svc := NewTransaction(store, slog.Default())
	var categoriID int64 = 1
	_, err := svc.Create(context.Background(), 1, CreateTransactionInput{
		Type: "expense", Amount: 200, AccountID: 3, CategoryID: &categoriID,
		Date: "2024-06-15",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	if !gotDate.Equal(want) {
		t.Errorf("date = %v, want %v", gotDate, want)
	}
}

func TestTransactionCreate_NilCategoryID(t *testing.T) {
	svc := NewTransaction(&mockStore{}, slog.Default())
	_, err := svc.Create(context.Background(), 1, CreateTransactionInput{
		Type: "income", Amount: 50, AccountID: 1, CategoryID: nil,
	})
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput for nil CategoryID, got %v", err)
	}
}

func TestTransactionCreate_WithoutDate_UsesNow(t *testing.T) {
	before := time.Now().Truncate(time.Second)
	var gotDate time.Time
	store := &mockStore{
		createTransactionFn: func(_ context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
			gotDate = tx.Date
			return tx, nil
		},
	}
	svc := NewTransaction(store, slog.Default())
	catID := int64(1)
	_, err := svc.Create(context.Background(), 1, CreateTransactionInput{
		Type: "income", Amount: 500, AccountID: 2, CategoryID: &catID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after := time.Now().Add(time.Second)
	if gotDate.Before(before) || gotDate.After(after) {
		t.Errorf("expected date near now, got %v", gotDate)
	}
}

func TestTransactionCreate_InvalidDateFallsBackToNow(t *testing.T) {
	before := time.Now().Truncate(time.Second)
	var gotDate time.Time
	store := &mockStore{
		createTransactionFn: func(_ context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
			gotDate = tx.Date
			return tx, nil
		},
	}
	svc := NewTransaction(store, slog.Default())
	catID := int64(1)
	_, err := svc.Create(context.Background(), 1, CreateTransactionInput{
		Type: "income", Amount: 100, AccountID: 1, CategoryID: &catID,
		Date: "not-a-date",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after := time.Now().Add(time.Second)
	if gotDate.Before(before) || gotDate.After(after) {
		t.Errorf("expected date near now for invalid input, got %v", gotDate)
	}
}
