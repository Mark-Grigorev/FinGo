package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func TestCategoryList_Success(t *testing.T) {
	env := newTestEnv()
	env.store.listCategoriesFn = func(_ context.Context, _ int64) ([]domain.Category, error) {
		return []domain.Category{{ID: 1, Name: "Food", Type: "expense"}}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryCreate_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCategoryCreate_InvalidType(t *testing.T) {
	env := newTestEnv()
	body := strings.NewReader(`{"name":"Food","type":"transfer"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/categories", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCategoryCreate_Success(t *testing.T) {
	env := newTestEnv()
	env.store.createCategoryFn = func(_ context.Context, c *domain.Category) (*domain.Category, error) {
		c.ID = 5
		return c, nil
	}

	body := strings.NewReader(`{"name":"Food","type":"expense","icon":"🍔","color":"#FF0000"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/categories", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCategoryUpdate_BadID(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/categories/abc", strings.NewReader(`{"name":"Food"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCategoryUpdate_BadJSON(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodPut, "/api/categories/1", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCategoryUpdate_Success(t *testing.T) {
	env := newTestEnv()
	env.store.updateCategoryFn = func(_ context.Context, id, _ int64, _, _, _ string) (*domain.Category, error) {
		return &domain.Category{ID: id, Name: "Transport"}, nil
	}

	body := strings.NewReader(`{"name":"Transport","icon":"🚗","color":"#0000FF"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/categories/1", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryDelete_BadID(t *testing.T) {
	env := newTestEnv()
	req := httptest.NewRequest(http.MethodDelete, "/api/categories/abc", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCategoryDelete_NotFound(t *testing.T) {
	env := newTestEnv()
	env.store.deleteCategoryFn = func(_ context.Context, _, _ int64) error {
		return domain.ErrNotFound
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/categories/99", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCategoryDelete_Success(t *testing.T) {
	env := newTestEnv()
	env.store.deleteCategoryFn = func(_ context.Context, _, _ int64) error { return nil }

	req := httptest.NewRequest(http.MethodDelete, "/api/categories/1", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}
