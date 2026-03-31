package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/internal/repository"
	"github.com/Mark-Grigorev/FinGo/internal/service"
)

type transactionHandler struct {
	svc *service.TransactionService
	log *slog.Logger
}

// list godoc
// @Summary List transactions
// @Description Get paginated list of user transactions with optional filters
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(20) minimum(1) maximum(100)
// @Param account_id query int false "Filter by account ID"
// @Param category_id query int false "Filter by category ID"
// @Param type query string false "Transaction type" Enums(income, expense)
// @Param start_date query string false "Start date (YYYY-MM-DD)" example:"2024-01-01"
// @Param end_date query string false "End date (YYYY-MM-DD)" example:"2024-12-31"
// @Success 200 {object} map[string]interface{} "Transactions retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions [get]
func (h *transactionHandler) list(c *gin.Context) {
	userID := currentUserID(c)

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	f := repository.TransactionFilter{Page: page, Limit: limit}

	list, total, err := h.svc.List(c.Request.Context(), userID, f)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": list,
		"total": total,
		"page":  f.Page,
		"limit": f.Limit,
	})
}

// delete godoc
// @Summary Delete transaction
// @Description Delete transaction by ID (can only delete own transactions)
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 204 "Transaction deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid transaction ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied (transaction belongs to another user)"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions/{id} [delete]
func (h *transactionHandler) delete(c *gin.Context) {
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

// create godoc
// @Summary Create transaction
// @Description Create a new income or expense transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body interface{} true "Transaction data"
// @Success 201 {object} map[string]interface{} "Transaction created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied (account/category belongs to another user)"
// @Failure 404 {object} map[string]interface{} "Account or category not found"
// @Failure 422 {object} map[string]interface{} "Insufficient funds for expense"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions [post]
func (h *transactionHandler) create(c *gin.Context) {
	userID := currentUserID(c)
	var in service.CreateTransactionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}
	tx, err := h.svc.Create(c.Request.Context(), userID, in)
	if err != nil {
		var insuf *domain.InsufficientFundsError
		if errors.As(err, &insuf) {
			alts := insuf.Alternatives
			if alts == nil {
				alts = []domain.Account{}
			}
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":         "insufficient_funds",
				"message":      "недостаточно средств на счёте «" + insuf.AccountName + "»",
				"balance":      insuf.Balance,
				"amount":       insuf.Amount,
				"alternatives": alts,
			})
			return
		}
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusCreated, tx)
}
