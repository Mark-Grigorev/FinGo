package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/internal/repository"
)

type AccountService struct {
	store repository.Storer
	log   *slog.Logger
}

func NewAccount(store repository.Storer, log *slog.Logger) *AccountService {
	return &AccountService{store: store, log: log}
}

type CreateAccountInput struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

type UpdateAccountInput struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Currency string `json:"currency"`
}

func (s *AccountService) List(ctx context.Context, userID int64) ([]domain.Account, error) {
	return s.store.ListAccounts(ctx, userID)
}

func (s *AccountService) Get(ctx context.Context, id, userID int64) (*domain.Account, error) {
	return s.store.GetAccount(ctx, id, userID)
}

func (s *AccountService) Create(ctx context.Context, userID int64, in CreateAccountInput) (*domain.Account, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return nil, domain.ErrInvalidInput
	}
	if in.Type == "" {
		in.Type = "card"
	}
	if in.Currency == "" {
		in.Currency = "RUB"
	}
	return s.store.CreateAccount(ctx, &domain.Account{
		UserID:   userID,
		Name:     in.Name,
		Type:     in.Type,
		Currency: in.Currency,
		Balance:  in.Balance,
	})
}

func (s *AccountService) Update(ctx context.Context, id, userID int64, in UpdateAccountInput) (*domain.Account, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return nil, domain.ErrInvalidInput
	}
	if in.Type == "" {
		in.Type = "card"
	}
	if in.Currency == "" {
		in.Currency = "RUB"
	}
	return s.store.UpdateAccount(ctx, id, userID, in.Name, in.Type, in.Currency)
}

func (s *AccountService) Delete(ctx context.Context, id, userID int64) error {
	return s.store.DeleteAccount(ctx, id, userID)
}
