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

// list godoc
// @Summary List recurring payments
// @Description Get all recurring payments for authenticated user
// @Tags recurring
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Recurring payments retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /recurring [get]
func (h *recurringHandler) list(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context(), currentUserID(c))
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list})
}

// create godoc
// @Summary Create recurring payment
// @Description Create a new recurring payment (subscription, regular expense, etc.)
// @Tags recurring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body interface{} true "Recurring payment data"
// @Success 201 {object} map[string]interface{} "Recurring payment created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied (account/category belongs to another user)"
// @Failure 404 {object} map[string]interface{} "Account not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /recurring [post]
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

// update godoc
// @Summary Update recurring payment
// @Description Update existing recurring payment by ID
// @Tags recurring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Recurring payment ID"
// @Param request body interface{} true "Recurring payment update data"
// @Success 200 {object} map[string]interface{} "Recurring payment updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied (payment belongs to another user)"
// @Failure 404 {object} map[string]interface{} "Recurring payment not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /recurring/{id} [put]
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

// delete godoc
// @Summary Delete recurring payment
// @Description Delete recurring payment by ID
// @Tags recurring
// @Produce json
// @Security BearerAuth
// @Param id path int true "Recurring payment ID"
// @Success 204 "Recurring payment deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied (payment belongs to another user)"
// @Failure 404 {object} map[string]interface{} "Recurring payment not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /recurring/{id} [delete]
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
