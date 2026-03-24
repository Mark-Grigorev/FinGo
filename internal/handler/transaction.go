package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Mark-Grigorev/FinGo/internal/repository"
	"github.com/Mark-Grigorev/FinGo/internal/service"
)

type transactionHandler struct {
	svc *service.TransactionService
	log *slog.Logger
}

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

func (h *transactionHandler) create(c *gin.Context) {
	userID := currentUserID(c)
	var in service.CreateTransactionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}
	tx, err := h.svc.Create(c.Request.Context(), userID, in)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusCreated, tx)
}
