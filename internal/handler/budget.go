package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Mark-Grigorev/FinGo/internal/service"
)

type budgetHandler struct {
	svc *service.BudgetService
	log *slog.Logger
}

// list godoc
// @Summary List budgets
// @Description Get all budgets for a specific month
// @Tags budgets
// @Produce json
// @Security BearerAuth
// @Param month query string false "Month (YYYY-MM)" example:"2024-01"
// @Success 200 {object} map[string]interface{} "Budgets retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /budgets [get]
func (h *budgetHandler) list(c *gin.Context) {
	userID := currentUserID(c)
	monthStr := c.Query("month")
	month := time.Now()
	if monthStr != "" {
		if parsed, err := time.Parse("2006-01", monthStr); err == nil {
			month = parsed
		}
	}
	list, err := h.svc.List(c.Request.Context(), userID, month)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list})
}

// create godoc
// @Summary Create budget
// @Description Create a new budget limit for a category
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body interface{} true "Budget data"
// @Success 201 {object} map[string]interface{} "Budget created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied (category belongs to another user)"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Failure 409 {object} map[string]interface{} "Budget already exists for this category and month"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /budgets [post]
func (h *budgetHandler) create(c *gin.Context) {
	userID := currentUserID(c)
	var in struct {
		CategoryID int64   `json:"category_id" binding:"required"`
		Month      string  `json:"month"`
		Limit      float64 `json:"limit" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}
	month := time.Now()
	if in.Month != "" {
		if parsed, err := time.Parse("2006-01", in.Month); err == nil {
			month = parsed
		}
	}
	b, err := h.svc.Create(c.Request.Context(), userID, in.CategoryID, month, in.Limit)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusCreated, b)
}

// update godoc
// @Summary Update budget
// @Description Update budget limit by ID
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Budget ID"
// @Param request body interface{} true "Budget update data"
// @Success 200 {object} map[string]interface{} "Budget updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied (budget belongs to another user)"
// @Failure 404 {object} map[string]interface{} "Budget not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /budgets/{id} [put]
func (h *budgetHandler) update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный id"})
		return
	}
	var in struct {
		Limit float64 `json:"limit" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}
	b, err := h.svc.Update(c.Request.Context(), id, currentUserID(c), in.Limit)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, b)
}

// delete godoc
// @Summary Delete budget
// @Description Delete budget by ID
// @Tags budgets
// @Produce json
// @Security BearerAuth
// @Param id path int true "Budget ID"
// @Success 204 "Budget deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied (budget belongs to another user)"
// @Failure 404 {object} map[string]interface{} "Budget not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /budgets/{id} [delete]
func (h *budgetHandler) delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id, currentUserID(c)); err != nil {
		writeError(c, h.log, err)
		return
	}
	c.Status(http.StatusNoContent)
}
