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
	store repository.Storer
	log   *slog.Logger
}

func NewTransaction(store repository.Storer, log *slog.Logger) *TransactionService {
	return &TransactionService{store: store, log: log}
}

type CreateTransactionInput struct {
	AccountID  int64   `json:"account_id"`
	Amount     float64 `json:"amount"`
	CategoryID *int64  `json:"category_id,omitempty"`
	Date       string  `json:"date"` // YYYY-MM-DD, optional
	Type       string  `json:"type"` // income | expense
	Name       string  `json:"name"`
}

func (s *TransactionService) List(ctx context.Context, userID int64, f repository.TransactionFilter) ([]domain.Transaction, int, error) {
	return s.store.ListTransactions(ctx, userID, f)
}

func (s *TransactionService) Export(ctx context.Context, userID int64, from, to time.Time) ([]domain.Transaction, error) {
	return s.store.ExportTransactions(ctx, userID, from, to)
}

func (s *TransactionService) Delete(ctx context.Context, id, userID int64) error {
	return s.store.DeleteTransaction(ctx, id, userID)
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
	if in.CategoryID == nil {
		return nil, domain.ErrInvalidInput
	}

	if in.Type == "expense" {
		acc, err := s.store.GetAccount(ctx, in.AccountID, userID)
		if err != nil {
			return nil, err
		}
		if acc.Type != "credit" && acc.Balance < in.Amount {
			all, _ := s.store.ListAccounts(ctx, userID)
			var alts []domain.Account
			for _, a := range all {
				if a.ID != in.AccountID && (a.Type == "credit" || a.Balance >= in.Amount) {
					alts = append(alts, a)
				}
			}
			return nil, &domain.InsufficientFundsError{
				AccountName:  acc.Name,
				Balance:      acc.Balance,
				Amount:       in.Amount,
				Alternatives: alts,
			}
		}
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
