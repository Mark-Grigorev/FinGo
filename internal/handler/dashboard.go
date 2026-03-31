package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Mark-Grigorev/FinGo/internal/service"
)

type dashboardHandler struct {
	svc *service.DashboardService
	log *slog.Logger
}

func (h *dashboardHandler) summary(c *gin.Context) {
	userID := currentUserID(c)
	s, err := h.svc.Summary(c.Request.Context(), userID)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *dashboardHandler) report(c *gin.Context) {
	userID := currentUserID(c)

	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		now := time.Now()
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		to = time.Now()
	}
	// включаем конец дня для to
	to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 0, time.UTC)

	result, err := h.svc.Report(c.Request.Context(), userID, from, to)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
