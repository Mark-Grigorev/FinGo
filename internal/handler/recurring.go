package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/internal/service"
)

type recurringHandler struct {
	svc *service.RecurringService
	log *slog.Logger
}

type recurringInput struct {
	AccountID       int64   `json:"account_id"  binding:"required"`
	CategoryID      *int64  `json:"category_id"`
	Name            string  `json:"name"        binding:"required"`
	Amount          float64 `json:"amount"      binding:"required,gt=0"`
	Frequency       string  `json:"frequency"`
	NextPaymentDate string  `json:"next_payment_date"`
}

func (h *recurringHandler) list(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context(), currentUserID(c))
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list})
}

func (h *recurringHandler) create(c *gin.Context) {
	var in recurringInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}
	nextDate := time.Now().AddDate(0, 1, 0)
	if in.NextPaymentDate != "" {
		if parsed, err := time.Parse("2006-01-02", in.NextPaymentDate); err == nil {
			nextDate = parsed
		}
	}
	r := domain.RecurringPayment{
		AccountID:       in.AccountID,
		CategoryID:      in.CategoryID,
		Name:            in.Name,
		Amount:          in.Amount,
		Frequency:       in.Frequency,
		NextPaymentDate: nextDate,
	}
	out, err := h.svc.Create(c.Request.Context(), currentUserID(c), r)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *recurringHandler) update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный id"})
		return
	}
	var in recurringInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}
	nextDate := time.Now().AddDate(0, 1, 0)
	if in.NextPaymentDate != "" {
		if parsed, err := time.Parse("2006-01-02", in.NextPaymentDate); err == nil {
			nextDate = parsed
		}
	}
	r := domain.RecurringPayment{
		AccountID:       in.AccountID,
		CategoryID:      in.CategoryID,
		Name:            in.Name,
		Amount:          in.Amount,
		Frequency:       in.Frequency,
		NextPaymentDate: nextDate,
	}
	out, err := h.svc.Update(c.Request.Context(), id, currentUserID(c), r)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *recurringHandler) delete(c *gin.Context) {
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
