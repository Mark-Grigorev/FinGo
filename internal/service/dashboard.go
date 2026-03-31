package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/repository"
)

type DashboardService struct {
	store repository.Storer
	log   *slog.Logger
}

func NewDashboard(store repository.Storer, log *slog.Logger) *DashboardService {
	return &DashboardService{store: store, log: log}
}

func (s *DashboardService) Summary(ctx context.Context, userID int64) (*repository.DashboardSummary, error) {
	return s.store.GetDashboardSummary(ctx, userID)
}

func (s *DashboardService) Report(ctx context.Context, userID int64, from, to time.Time) (*repository.ReportResult, error) {
	return s.store.GetReport(ctx, userID, from, to)
}
