package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Mark-Grigorev/FinGo/internal/service"
)

type currencyHandler struct {
	svc *service.CurrencyService
	log *slog.Logger
}

func (h *currencyHandler) listRates(c *gin.Context) {
	rates, err := h.svc.ListRates(c.Request.Context(), currentUserID(c))
	if err != nil { writeError(c, h.log, err); return }
	c.JSON(http.StatusOK, gin.H{"items": rates})
}

func (h *currencyHandler) upsertRate(c *gin.Context) {
	currency := c.Param("currency")
	var in struct {
		Rate float64 `json:"rate" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}
	r, err := h.svc.UpsertRate(c.Request.Context(), currentUserID(c), currency, in.Rate)
	if err != nil { writeError(c, h.log, err); return }
	c.JSON(http.StatusOK, r)
}

func (h *currencyHandler) deleteRate(c *gin.Context) {
	if err := h.svc.DeleteRate(c.Request.Context(), currentUserID(c), c.Param("currency")); err != nil {
		writeError(c, h.log, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *currencyHandler) getBaseCurrency(c *gin.Context) {
	base, err := h.svc.GetBaseCurrency(c.Request.Context(), currentUserID(c))
	if err != nil { writeError(c, h.log, err); return }
	c.JSON(http.StatusOK, gin.H{"base_currency": base})
}

func (h *currencyHandler) setBaseCurrency(c *gin.Context) {
	var in struct {
		Currency string `json:"currency" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}
	if err := h.svc.SetBaseCurrency(c.Request.Context(), currentUserID(c), in.Currency); err != nil {
		writeError(c, h.log, err)
		return
	}
	c.Status(http.StatusNoContent)
}
