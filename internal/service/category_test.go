package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestCategoryList_Success(t *testing.T) {
	want := []domain.Category{{ID: 1, Name: "Food", Type: "expense"}}
	store := &mockStore{
		listCategoriesFn: func(_ context.Context, _ int64) ([]domain.Category, error) {
			return want, nil
		},
	}
	svc := NewCategory(store, slog.Default())
	got, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Food" {
		t.Errorf("unexpected list: %v", got)
	}
}

func TestCategoryCreate_InvalidInput(t *testing.T) {
	svc := NewCategory(&mockStore{}, slog.Default())
	cases := []struct{ name, typ string }{
		{"", "expense"},
		{"Food", ""},
		{"Food", "transfer"},
		{"Food", "INCOME"},
	}
	for _, c := range cases {
		_, err := svc.Create(context.Background(), 1, c.name, "", "", c.typ)
		if err != domain.ErrInvalidInput {
			t.Errorf("Create(%q,%q): expected ErrInvalidInput, got %v", c.name, c.typ, err)
		}
	}
}

func TestCategoryCreate_Defaults(t *testing.T) {
	var got *domain.Category
	store := &mockStore{
		createCategoryFn: func(_ context.Context, c *domain.Category) (*domain.Category, error) {
			got = c
			return c, nil
		},
	}
	svc := NewCategory(store, slog.Default())
	_, err := svc.Create(context.Background(), 1, "Food", "", "", "expense")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Icon == "" {
		t.Error("expected default icon")
	}
	if got.Color == "" {
		t.Error("expected default color")
	}
}

func TestCategoryCreate_Success(t *testing.T) {
	want := &domain.Category{ID: 3, Name: "Salary", Type: "income"}
	store := &mockStore{
		createCategoryFn: func(_ context.Context, c *domain.Category) (*domain.Category, error) {
			return want, nil
		},
	}
	svc := NewCategory(store, slog.Default())
	got, err := svc.Create(context.Background(), 1, "Salary", "💰", "#00FF00", "income")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 3 {
		t.Errorf("got.ID = %d, want 3", got.ID)
	}
}

func TestCategoryUpdate_EmptyName(t *testing.T) {
	svc := NewCategory(&mockStore{}, slog.Default())
	_, err := svc.Update(context.Background(), 1, 1, "", "", "")
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCategoryUpdate_Success(t *testing.T) {
	want := &domain.Category{ID: 2, Name: "Transport"}
	store := &mockStore{
		updateCategoryFn: func(_ context.Context, id, _ int64, _, _, _ string) (*domain.Category, error) {
			return want, nil
		},
	}
	svc := NewCategory(store, slog.Default())
	got, err := svc.Update(context.Background(), 2, 1, "Transport", "🚗", "#FF0000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 2 {
		t.Errorf("got.ID = %d, want 2", got.ID)
	}
}

func TestCategoryDelete_Success(t *testing.T) {
	called := false
	store := &mockStore{
		deleteCategoryFn: func(_ context.Context, id, userID int64) error {
			called = true
			return nil
		},
	}
	svc := NewCategory(store, slog.Default())
	if err := svc.Delete(context.Background(), 1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("DeleteCategory was not called")
	}
}
