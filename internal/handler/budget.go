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
