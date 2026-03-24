package service

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/internal/repository"
)

type TransactionService struct {
	store *repository.Store
	log   *slog.Logger
}

func NewTransaction(store *repository.Store, log *slog.Logger) *TransactionService {
	return &TransactionService{store: store, log: log}
}

type CreateTransactionInput struct {
	AccountID  int64   `json:"account_id"`
	CategoryID *int64  `json:"category_id,omitempty"`
	Type       string  `json:"type"` // income | expense
	Amount     float64 `json:"amount"`
	Name       string  `json:"name"`
	Date       string  `json:"date"` // YYYY-MM-DD, optional
}

func (s *TransactionService) List(ctx context.Context, userID int64, f repository.TransactionFilter) ([]domain.Transaction, int, error) {
	return s.store.ListTransactions(ctx, userID, f)
}

func (s *TransactionService) Create(ctx context.Context, userID int64, in CreateTransactionInput) (*domain.Transaction, error) {
	if in.Type != "income" && in.Type != "expense" {
		return nil, domain.ErrInvalidInput
	}
	if in.Amount <= 0 {
		return nil, domain.ErrInvalidInput
	}
	if in.AccountID == 0 {
		return nil, domain.ErrInvalidInput
	}

	date := time.Now()
	if in.Date != "" {
		if parsed, err := time.Parse("2006-01-02", in.Date); err == nil {
			date = parsed
		}
	}

	return s.store.CreateTransaction(ctx, &domain.Transaction{
		UserID:     userID,
		AccountID:  in.AccountID,
		CategoryID: in.CategoryID,
		Type:       in.Type,
		Amount:     in.Amount,
		Name:       strings.TrimSpace(in.Name),
		Date:       date,
	})
}
