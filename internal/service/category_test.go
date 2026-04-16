package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Food", got[0].Name)
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
		assert.ErrorIsf(t, err, domain.ErrInvalidInput,
			"Create(%q, %q)", c.name, c.typ)
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
	require.NoError(t, err)
	assert.NotEmpty(t, got.Icon)
	assert.NotEmpty(t, got.Color)
}

func TestCategoryCreate_Success(t *testing.T) {
	want := &domain.Category{ID: 3, Name: "Salary", Type: "income"}
	store := &mockStore{
		createCategoryFn: func(_ context.Context, _ *domain.Category) (*domain.Category, error) {
			return want, nil
		},
	}
	svc := NewCategory(store, slog.Default())
	got, err := svc.Create(context.Background(), 1, "Salary", "💰", "#00FF00", "income")
	require.NoError(t, err)
	assert.Equal(t, int64(3), got.ID)
}

func TestCategoryUpdate_EmptyName(t *testing.T) {
	svc := NewCategory(&mockStore{}, slog.Default())
	_, err := svc.Update(context.Background(), 1, 1, "", "", "")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
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
	require.NoError(t, err)
	assert.Equal(t, int64(2), got.ID)
}

func TestCategoryDelete_Success(t *testing.T) {
	called := false
	store := &mockStore{
		deleteCategoryFn: func(_ context.Context, _, _ int64) error {
			called = true
			return nil
		},
	}
	svc := NewCategory(store, slog.Default())
	require.NoError(t, svc.Delete(context.Background(), 1, 1))
	assert.True(t, called, "DeleteCategory was not called")
}
