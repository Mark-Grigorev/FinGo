package service

import (
	"context"
	"log/slog"
	"testing"

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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Currency != "USD" {
		t.Errorf("unexpected rates: %v", got)
	}
}

func TestCurrencyUpsertRate_UnsupportedCurrency(t *testing.T) {
	svc := NewCurrency(&mockStore{}, slog.Default())
	_, err := svc.UpsertRate(context.Background(), 1, "XYZ", 1.5)
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput for unsupported currency, got %v", err)
	}
}

func TestCurrencyUpsertRate_InvalidRate(t *testing.T) {
	svc := NewCurrency(&mockStore{}, slog.Default())
	for _, rate := range []float64{0, -1} {
		_, err := svc.UpsertRate(context.Background(), 1, "USD", rate)
		if err != domain.ErrInvalidInput {
			t.Errorf("rate=%v: expected ErrInvalidInput, got %v", rate, err)
		}
	}
}

func TestCurrencyUpsertRate_Success(t *testing.T) {
	want := &domain.ExchangeRate{Currency: "USD", Rate: 90.5}
	store := &mockStore{
		upsertExchangeRateFn: func(_ context.Context, _ int64, currency string, rate float64) (*domain.ExchangeRate, error) {
			return want, nil
		},
	}
	svc := NewCurrency(store, slog.Default())
	got, err := svc.UpsertRate(context.Background(), 1, "usd", 90.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Currency != "USD" {
		t.Errorf("currency = %q, want USD", got.Currency)
	}
}

func TestCurrencyDeleteRate_Success(t *testing.T) {
	called := false
	store := &mockStore{
		deleteExchangeRateFn: func(_ context.Context, _ int64, currency string) error {
			called = true
			if currency != "EUR" {
				t.Errorf("currency = %q, want EUR", currency)
			}
			return nil
		},
	}
	svc := NewCurrency(store, slog.Default())
	if err := svc.DeleteRate(context.Background(), 1, "eur"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("DeleteExchangeRate was not called")
	}
}

func TestCurrencyGetBaseCurrency_Success(t *testing.T) {
	store := &mockStore{
		getBaseCurrencyFn: func(_ context.Context, _ int64) (string, error) {
			return "USD", nil
		},
	}
	svc := NewCurrency(store, slog.Default())
	got, err := svc.GetBaseCurrency(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "USD" {
		t.Errorf("base currency = %q, want USD", got)
	}
}

func TestCurrencySetBaseCurrency_Unsupported(t *testing.T) {
	svc := NewCurrency(&mockStore{}, slog.Default())
	if err := svc.SetBaseCurrency(context.Background(), 1, "XYZ"); err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCurrencySetBaseCurrency_Success(t *testing.T) {
	called := false
	store := &mockStore{
		setBaseCurrencyFn: func(_ context.Context, _ int64, currency string) error {
			called = true
			if currency != "EUR" {
				t.Errorf("currency = %q, want EUR", currency)
			}
			return nil
		},
	}
	svc := NewCurrency(store, slog.Default())
	if err := svc.SetBaseCurrency(context.Background(), 1, "eur"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("SetBaseCurrency was not called")
	}
}
