package service

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/repository"
)

func TestDashboardSummary_Success(t *testing.T) {
	want := &repository.DashboardSummary{Balance: 10000, Income: 5000, Expenses: 2000}
	store := &mockStore{
		getDashboardSummaryFn: func(_ context.Context, _ int64) (*repository.DashboardSummary, error) {
			return want, nil
		},
	}
	svc := NewDashboard(store, slog.Default())
	got, err := svc.Summary(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Balance != 10000 {
		t.Errorf("Balance = %v, want 10000", got.Balance)
	}
	if got.Income != 5000 {
		t.Errorf("Income = %v, want 5000", got.Income)
	}
}

func TestDashboardReport_Success(t *testing.T) {
	want := &repository.ReportResult{Comparison: repository.ComparisonData{Income: 8000}}
	store := &mockStore{
		getReportFn: func(_ context.Context, _ int64, from, to time.Time) (*repository.ReportResult, error) {
			return want, nil
		},
	}
	svc := NewDashboard(store, slog.Default())
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	got, err := svc.Report(context.Background(), 1, from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Comparison.Income != 8000 {
		t.Errorf("Comparison.Income = %v, want 8000", got.Comparison.Income)
	}
}
