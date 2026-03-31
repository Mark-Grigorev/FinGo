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

// summary godoc
// @Summary Get dashboard summary
// @Description Get financial summary including total balance, income, expenses, and recent transactions
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Dashboard summary retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dashboard/summary [get]
func (h *dashboardHandler) summary(c *gin.Context) {
	userID := currentUserID(c)
	s, err := h.svc.Summary(c.Request.Context(), userID)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, s)
}

// report godoc
// @Summary Get financial report
// @Description Get detailed financial report for date range (income/expense breakdown by category)
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Param from query string false "Start date (YYYY-MM-DD)" example:"2024-01-01"
// @Param to query string false "End date (YYYY-MM-DD)" example:"2024-12-31"
// @Success 200 {object} map[string]interface{} "Financial report retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid date format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dashboard/report [get]
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
