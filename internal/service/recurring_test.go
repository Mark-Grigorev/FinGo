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

func TestRecurringList_Success(t *testing.T) {
	want := []domain.RecurringPayment{{ID: 1, Name: "Netflix", Amount: 500}}
	store := &mockStore{
		listRecurringFn: func(_ context.Context, _ int64) ([]domain.RecurringPayment, error) {
			return want, nil
		},
	}
	svc := NewRecurring(store, slog.Default())
	got, err := svc.List(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Netflix", got[0].Name)
}

func TestRecurringCreate_InvalidInput(t *testing.T) {
	svc := NewRecurring(&mockStore{}, slog.Default())
	cases := []domain.RecurringPayment{
		{Name: "", Amount: 100, AccountID: 1},    // пустое имя
		{Name: "  ", Amount: 100, AccountID: 1},  // пробельное имя
		{Name: "Sub", Amount: 0, AccountID: 1},   // нулевая сумма
		{Name: "Sub", Amount: -10, AccountID: 1}, // отрицательная сумма
		{Name: "Sub", Amount: 100, AccountID: 0}, // нулевой аккаунт
	}
	for _, r := range cases {
		_, err := svc.Create(context.Background(), 1, r)
		assert.ErrorIsf(t, err, domain.ErrInvalidInput, "input: %+v", r)
	}
}

func TestRecurringCreate_DefaultFrequency(t *testing.T) {
	var gotFreq string
	store := &mockStore{
		createRecurringFn: func(_ context.Context, r *domain.RecurringPayment) (*domain.RecurringPayment, error) {
			gotFreq = r.Frequency
			return r, nil
		},
	}
	svc := NewRecurring(store, slog.Default())
	_, err := svc.Create(context.Background(), 1, domain.RecurringPayment{
		Name: "Netflix", Amount: 500, AccountID: 1, Frequency: "biweekly",
	})
	require.NoError(t, err)
	assert.Equal(t, "monthly", gotFreq)
}

func TestRecurringCreate_Success(t *testing.T) {
	want := &domain.RecurringPayment{ID: 7, Name: "Rent", Amount: 30000, Frequency: "monthly"}
	store := &mockStore{
		createRecurringFn: func(_ context.Context, _ *domain.RecurringPayment) (*domain.RecurringPayment, error) {
			return want, nil
		},
	}
	svc := NewRecurring(store, slog.Default())
	got, err := svc.Create(context.Background(), 1, domain.RecurringPayment{
		Name: "Rent", Amount: 30000, AccountID: 1, Frequency: "monthly",
	})
	require.NoError(t, err)
	assert.Equal(t, int64(7), got.ID)
}

func TestRecurringUpdate_InvalidInput(t *testing.T) {
	svc := NewRecurring(&mockStore{}, slog.Default())
	_, err := svc.Update(context.Background(), 1, 1, domain.RecurringPayment{
		Name: "", Amount: 100, AccountID: 1,
	})
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestRecurringUpdate_Success(t *testing.T) {
	want := &domain.RecurringPayment{ID: 3, Name: "Gym", Amount: 1500}
	store := &mockStore{
		updateRecurringFn: func(_ context.Context, id, _ int64, _ string, _ float64, _ string, _ time.Time, _ int64, _ *int64) (*domain.RecurringPayment, error) {
			return want, nil
		},
	}
	svc := NewRecurring(store, slog.Default())
	got, err := svc.Update(context.Background(), 3, 1, domain.RecurringPayment{
		Name: "Gym", Amount: 1500, AccountID: 2, Frequency: "monthly",
	})
	require.NoError(t, err)
	assert.Equal(t, int64(3), got.ID)
}

func TestRecurringDelete_Success(t *testing.T) {
	called := false
	store := &mockStore{
		deleteRecurringFn: func(_ context.Context, _, _ int64) error {
			called = true
			return nil
		},
	}
	svc := NewRecurring(store, slog.Default())
	require.NoError(t, svc.Delete(context.Background(), 1, 1))
	assert.True(t, called, "DeleteRecurring was not called")
}
