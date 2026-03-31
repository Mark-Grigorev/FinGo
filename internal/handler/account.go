package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Mark-Grigorev/FinGo/internal/service"
)

type accountHandler struct {
	svc *service.AccountService
	log *slog.Logger
}

// list godoc
// @Summary List all accounts
// @Description Get list of all accounts for authenticated user
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Accounts retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /accounts [get]
func (h *accountHandler) list(c *gin.Context) {
	userID := currentUserID(c)
	accounts, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": accounts})
}

// create godoc
// @Summary Create new account
// @Description Create a new financial account (debit, credit, cash)
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateAccountInput true "Account creation data"
// @Success 201 {object} interface{} "Account created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 409 {object} map[string]interface{} "Account with this name already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /accounts [post]
func (h *accountHandler) create(c *gin.Context) {
	userID := currentUserID(c)
	var in service.CreateAccountInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}
	acc, err := h.svc.Create(c.Request.Context(), userID, in)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusCreated, acc)
}

// update godoc
// @Summary Update account
// @Description Update existing account by ID
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Account ID"
// @Param request body service.UpdateAccountInput true "Account update data"
// @Success 201 {object} interface{} "Account created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied (account belongs to another user)"
// @Failure 404 {object} map[string]interface{} "Account not found"
// @Failure 409 {object} map[string]interface{} "Account with this name already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /accounts/{id} [put]
func (h *accountHandler) update(c *gin.Context) {
	userID := currentUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный id"})
		return
	}
	var in service.UpdateAccountInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}
	acc, err := h.svc.Update(c.Request.Context(), id, userID, in)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, acc)
}

// delete godoc
// @Summary Delete account
// @Description Delete account by ID (only if balance is zero)
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Param id path int true "Account ID"
// @Success 204 "Account deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Access denied (account belongs to another user)"
// @Failure 404 {object} map[string]interface{} "Account not found"
// @Failure 409 {object} map[string]interface{} "Cannot delete account with non-zero balance"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /accounts/{id} [delete]
func (h *accountHandler) delete(c *gin.Context) {
	userID := currentUserID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		writeError(c, h.log, err)
		return
	}
	c.Status(http.StatusNoContent)
}
