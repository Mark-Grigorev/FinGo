package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/repository"
)

func TestDashboardSummary_Success(t *testing.T) {
	env := newTestEnv()
	env.store.getDashboardSummaryFn = func(_ context.Context, _ int64) (*repository.DashboardSummary, error) {
		return &repository.DashboardSummary{Balance: 10000, Income: 5000}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestDashboardReport_DefaultDates(t *testing.T) {
	env := newTestEnv()
	env.store.getReportFn = func(_ context.Context, _ int64, _ time.Time, _ time.Time) (*repository.ReportResult, error) {
		return &repository.ReportResult{}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/report", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestDashboardReport_WithDates(t *testing.T) {
	env := newTestEnv()
	env.store.getReportFn = func(_ context.Context, _ int64, _ time.Time, _ time.Time) (*repository.ReportResult, error) {
		return &repository.ReportResult{}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/report?from=2024-01-01&to=2024-01-31", nil)
	req.Header.Set("Authorization", env.bearerToken(1))
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}
