package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/internal/repository"
)

type BudgetService struct {
	store repository.Storer
	log   *slog.Logger
}

func NewBudget(store repository.Storer, log *slog.Logger) *BudgetService {
	return &BudgetService{store: store, log: log}
}

func (s *BudgetService) List(ctx context.Context, userID int64, month time.Time) ([]domain.Budget, error) {
	return s.store.ListBudgets(ctx, userID, month)
}

func (s *BudgetService) Create(ctx context.Context, userID int64, categoryID int64, month time.Time, limit float64) (*domain.Budget, error) {
	if limit <= 0 {
		return nil, domain.ErrInvalidInput
	}
	b := &domain.Budget{
		UserID:     userID,
		CategoryID: categoryID,
		Month:      month.Format("2006-01-02"),
		Limit:      limit,
	}
	return s.store.CreateBudget(ctx, b)
}

func (s *BudgetService) Update(ctx context.Context, id, userID int64, limit float64) (*domain.Budget, error) {
	if limit <= 0 {
		return nil, domain.ErrInvalidInput
	}
	return s.store.UpdateBudget(ctx, id, userID, limit)
}

func (s *BudgetService) Delete(ctx context.Context, id, userID int64) error {
	return s.store.DeleteBudget(ctx, id, userID)
}
