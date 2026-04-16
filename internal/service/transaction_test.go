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

func TestTransactionCreate_InvalidType(t *testing.T) {
	svc := NewTransaction(&mockStore{}, slog.Default())
	for _, typ := range []string{"", "transfer", "INCOME"} {
		_, err := svc.Create(context.Background(), 1, CreateTransactionInput{
			Type: typ, Amount: 100, AccountID: 1,
		})
		assert.ErrorIsf(t, err, domain.ErrInvalidInput, "type=%q", typ)
	}
}

func TestTransactionCreate_InvalidAmount(t *testing.T) {
	svc := NewTransaction(&mockStore{}, slog.Default())
	for _, amount := range []float64{0, -1, -100} {
		_, err := svc.Create(context.Background(), 1, CreateTransactionInput{
			Type: "income", Amount: amount, AccountID: 1,
		})
		assert.ErrorIsf(t, err, domain.ErrInvalidInput, "amount=%v", amount)
	}
}

func TestTransactionCreate_ZeroAccountID(t *testing.T) {
	svc := NewTransaction(&mockStore{}, slog.Default())
	_, err := svc.Create(context.Background(), 1, CreateTransactionInput{
		Type: "income", Amount: 50, AccountID: 0,
	})
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestTransactionCreate_NilCategoryID(t *testing.T) {
	svc := NewTransaction(&mockStore{}, slog.Default())
	_, err := svc.Create(context.Background(), 1, CreateTransactionInput{
		Type: "income", Amount: 50, AccountID: 1, CategoryID: nil,
	})
	require.ErrorIs(t, err, domain.ErrInvalidInput)
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
	catID := int64(1)
	_, err := svc.Create(context.Background(), 1, CreateTransactionInput{
		Type: "expense", Amount: 200, AccountID: 3, CategoryID: &catID,
		Date: "2024-06-15",
	})
	require.NoError(t, err)
	assert.Equal(t, time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), gotDate)
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
	require.NoError(t, err)
	assert.False(t, gotDate.Before(before), "date should not be before test start")
	assert.False(t, gotDate.After(time.Now().Add(time.Second)), "date should not be in the future")
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
	require.NoError(t, err)
	assert.False(t, gotDate.Before(before), "date should not be before test start")
	assert.False(t, gotDate.After(time.Now().Add(time.Second)), "date should not be in the future")
}
