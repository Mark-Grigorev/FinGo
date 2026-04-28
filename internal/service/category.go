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

type seedCategory struct {
	name, icon, color, typ string
}

var defaultCategories = []seedCategory{
	// Расходы
	{"Продукты", "🛒", "#4CAF50", "expense"},
	{"Транспорт", "🚗", "#2196F3", "expense"},
	{"Жильё", "🏠", "#9C27B0", "expense"},
	{"Рестораны", "🍽️", "#FF5722", "expense"},
	{"Здоровье", "💊", "#E91E63", "expense"},
	{"Одежда и обувь", "👗", "#FF9800", "expense"},
	{"Развлечения", "🎮", "#00BCD4", "expense"},
	{"Связь", "📱", "#607D8B", "expense"},
	{"Образование", "📚", "#795548", "expense"},
	{"Питомцы", "🐾", "#8BC34A", "expense"},
	// Доходы
	{"Зарплата", "💼", "#4ade80", "income"},
	{"Фриланс", "💻", "#60a5fa", "income"},
	{"Инвестиции", "📈", "#f59e0b", "income"},
	{"Подарки", "🎁", "#a78bfa", "income"},
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

func (s *CategoryService) SeedDefaultCategories(ctx context.Context, userID int64) error {
	for _, dc := range defaultCategories {
		c := &domain.Category{UserID: userID, Name: dc.name, Icon: dc.icon, Color: dc.color, Type: dc.typ}
		if _, err := s.store.CreateCategory(ctx, c); err != nil {
			s.log.Error("seed category failed", "name", dc.name, "err", err)
		}
	}
	return nil
}
