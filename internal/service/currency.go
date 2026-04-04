package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/internal/repository"
)

var supportedCurrencies = map[string]bool{
	"RUB": true, "USD": true, "EUR": true, "CNY": true,
	"GBP": true, "JPY": true, "CHF": true, "TRY": true,
}

type CurrencyService struct {
	store repository.Storer
	log   *slog.Logger
}

func NewCurrency(store repository.Storer, log *slog.Logger) *CurrencyService {
	return &CurrencyService{store: store, log: log}
}

func (s *CurrencyService) ListRates(ctx context.Context, userID int64) ([]domain.ExchangeRate, error) {
	return s.store.ListExchangeRates(ctx, userID)
}

func (s *CurrencyService) UpsertRate(ctx context.Context, userID int64, currency string, rate float64) (*domain.ExchangeRate, error) {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if !supportedCurrencies[currency] { return nil, domain.ErrInvalidInput }
	if rate <= 0 { return nil, domain.ErrInvalidInput }
	return s.store.UpsertExchangeRate(ctx, userID, currency, rate)
}

func (s *CurrencyService) DeleteRate(ctx context.Context, userID int64, currency string) error {
	return s.store.DeleteExchangeRate(ctx, userID, strings.ToUpper(currency))
}

func (s *CurrencyService) GetBaseCurrency(ctx context.Context, userID int64) (string, error) {
	return s.store.GetBaseCurrency(ctx, userID)
}

func (s *CurrencyService) SetBaseCurrency(ctx context.Context, userID int64, currency string) error {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if !supportedCurrencies[currency] { return domain.ErrInvalidInput }
	return s.store.SetBaseCurrency(ctx, userID, currency)
}
