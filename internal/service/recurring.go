package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/internal/repository"
)

type RecurringService struct {
	store repository.Storer
	log   *slog.Logger
}

func NewRecurring(store repository.Storer, log *slog.Logger) *RecurringService {
	return &RecurringService{store: store, log: log}
}

var validFrequencies = map[string]bool{"monthly": true, "weekly": true, "yearly": true}

func (s *RecurringService) List(ctx context.Context, userID int64) ([]domain.RecurringPayment, error) {
	return s.store.ListRecurring(ctx, userID)
}

func (s *RecurringService) Create(ctx context.Context, userID int64, r domain.RecurringPayment) (*domain.RecurringPayment, error) {
	if strings.TrimSpace(r.Name) == "" || r.Amount <= 0 || r.AccountID == 0 {
		return nil, domain.ErrInvalidInput
	}
	if !validFrequencies[r.Frequency] {
		r.Frequency = "monthly"
	}
	r.UserID = userID
	return s.store.CreateRecurring(ctx, &r)
}

func (s *RecurringService) Update(ctx context.Context, id, userID int64, r domain.RecurringPayment) (*domain.RecurringPayment, error) {
	if strings.TrimSpace(r.Name) == "" || r.Amount <= 0 || r.AccountID == 0 {
		return nil, domain.ErrInvalidInput
	}
	if !validFrequencies[r.Frequency] {
		r.Frequency = "monthly"
	}
	return s.store.UpdateRecurring(ctx, id, userID, r.Name, r.Amount, r.Frequency, r.NextPaymentDate, r.AccountID, r.CategoryID)
}

func (s *RecurringService) Delete(ctx context.Context, id, userID int64) error {
	return s.store.DeleteRecurring(ctx, id, userID)
}
