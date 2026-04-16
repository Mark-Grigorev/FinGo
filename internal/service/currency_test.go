package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestCurrencyListRates_Success(t *testing.T) {
	want := []domain.ExchangeRate{{ID: 1, Currency: "USD", Rate: 90.5}}
	store := &mockStore{
		listExchangeRatesFn: func(_ context.Context, _ int64) ([]domain.ExchangeRate, error) {
			return want, nil
		},
	}
	svc := NewCurrency(store, slog.Default())
	got, err := svc.ListRates(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "USD", got[0].Currency)
}

func TestCurrencyUpsertRate_UnsupportedCurrency(t *testing.T) {
	svc := NewCurrency(&mockStore{}, slog.Default())
	_, err := svc.UpsertRate(context.Background(), 1, "XYZ", 1.5)
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestCurrencyUpsertRate_InvalidRate(t *testing.T) {
	svc := NewCurrency(&mockStore{}, slog.Default())
	for _, rate := range []float64{0, -1} {
		_, err := svc.UpsertRate(context.Background(), 1, "USD", rate)
		assert.ErrorIsf(t, err, domain.ErrInvalidInput, "rate=%v", rate)
	}
}

func TestCurrencyUpsertRate_Success(t *testing.T) {
	want := &domain.ExchangeRate{Currency: "USD", Rate: 90.5}
	store := &mockStore{
		upsertExchangeRateFn: func(_ context.Context, _ int64, _ string, _ float64) (*domain.ExchangeRate, error) {
			return want, nil
		},
	}
	svc := NewCurrency(store, slog.Default())
	got, err := svc.UpsertRate(context.Background(), 1, "usd", 90.5)
	require.NoError(t, err)
	assert.Equal(t, "USD", got.Currency)
}

func TestCurrencyDeleteRate_Success(t *testing.T) {
	var gotCurrency string
	store := &mockStore{
		deleteExchangeRateFn: func(_ context.Context, _ int64, currency string) error {
			gotCurrency = currency
			return nil
		},
	}
	svc := NewCurrency(store, slog.Default())
	require.NoError(t, svc.DeleteRate(context.Background(), 1, "eur"))
	assert.Equal(t, "EUR", gotCurrency)
}

func TestCurrencyGetBaseCurrency_Success(t *testing.T) {
	store := &mockStore{
		getBaseCurrencyFn: func(_ context.Context, _ int64) (string, error) {
			return "USD", nil
		},
	}
	svc := NewCurrency(store, slog.Default())
	got, err := svc.GetBaseCurrency(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "USD", got)
}

func TestCurrencySetBaseCurrency_Unsupported(t *testing.T) {
	svc := NewCurrency(&mockStore{}, slog.Default())
	err := svc.SetBaseCurrency(context.Background(), 1, "XYZ")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestCurrencySetBaseCurrency_Success(t *testing.T) {
	var gotCurrency string
	store := &mockStore{
		setBaseCurrencyFn: func(_ context.Context, _ int64, currency string) error {
			gotCurrency = currency
			return nil
		},
	}
	svc := NewCurrency(store, slog.Default())
	require.NoError(t, svc.SetBaseCurrency(context.Background(), 1, "eur"))
	assert.Equal(t, "EUR", gotCurrency)
}
