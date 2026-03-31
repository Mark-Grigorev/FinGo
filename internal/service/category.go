package service

import (
	"context"
	"log/slog"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/internal/repository"
)

type CategoryService struct {
	store repository.Storer
	log   *slog.Logger
}

func NewCategory(store repository.Storer, log *slog.Logger) *CategoryService {
	return &CategoryService{store: store, log: log}
}

func (s *CategoryService) List(ctx context.Context, userID int64) ([]domain.Category, error) {
	return s.store.ListCategories(ctx, userID)
}

func (s *CategoryService) Create(ctx context.Context, userID int64, name, icon, color, typ string) (*domain.Category, error) {
	if name == "" || (typ != "income" && typ != "expense") {
		return nil, domain.ErrInvalidInput
	}
	if icon == "" {
		icon = "💳"
	}
	if color == "" {
		color = "#888888"
	}
	c := &domain.Category{UserID: userID, Name: name, Icon: icon, Color: color, Type: typ}
	return s.store.CreateCategory(ctx, c)
}

func (s *CategoryService) Update(ctx context.Context, id, userID int64, name, icon, color string) (*domain.Category, error) {
	if name == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.store.UpdateCategory(ctx, id, userID, name, icon, color)
}

func (s *CategoryService) Delete(ctx context.Context, id, userID int64) error {
	return s.store.DeleteCategory(ctx, id, userID)
}
