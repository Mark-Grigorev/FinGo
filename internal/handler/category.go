package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Mark-Grigorev/FinGo/internal/service"
)

type categoryHandler struct {
	svc *service.CategoryService
	log *slog.Logger
}

// list godoc
// @Summary List categories
// @Description Get all categories for authenticated user (income and expense types)
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Categories retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories [get]
func (h *categoryHandler) list(c *gin.Context) {
	userID := currentUserID(c)
	categories, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": categories})
}

// create godoc
// @Summary Create category
// @Description Create a new transaction category (income or expense)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body interface{} true "Category data"
// @Success 201 {object} map[string]interface{} "Category created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 409 {object} map[string]interface{} "Category with this name already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories [post]
func (h *categoryHandler) create(c *gin.Context) {
	var in struct {
		Name  string `json:"name"  binding:"required"`
		Icon  string `json:"icon"`
		Color string `json:"color"`
		Type  string `json:"type"  binding:"required,oneof=income expense"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}

	userID := currentUserID(c)
	cat, err := h.svc.Create(c.Request.Context(), userID, in.Name, in.Icon, in.Color, in.Type)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusCreated, cat)
}

// update godoc
// @Summary Update category
// @Description Update existing category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body interface{} true "Category update data"
// @Success 200 {object} map[string]interface{} "Category updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied (category belongs to another user)"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Failure 409 {object} map[string]interface{} "Category with this name already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories/{id} [put]
func (h *categoryHandler) update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный id"})
		return
	}

	var in struct {
		Name  string `json:"name"  binding:"required"`
		Icon  string `json:"icon"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}

	userID := currentUserID(c)
	cat, err := h.svc.Update(c.Request.Context(), id, userID, in.Name, in.Icon, in.Color)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, cat)
}

// delete godoc
// @Summary Delete category
// @Description Delete category by ID (only if no transactions linked)
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 204 "Category deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied (category belongs to another user)"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Failure 409 {object} map[string]interface{} "Cannot delete category with existing transactions"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories/{id} [delete]
func (h *categoryHandler) delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный id"})
		return
	}

	userID := currentUserID(c)
	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		writeError(c, h.log, err)
		return
	}
	c.Status(http.StatusNoContent)
}
