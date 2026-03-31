package repository

import (
	"context"
	"time"
)

type DashboardSummary struct {
	Balance     float64 `json:"balance"`
	Income      float64 `json:"income"`
	Expenses    float64 `json:"expenses"`
	Savings     float64 `json:"savings"`
	BalancePct  float64 `json:"balance_pct"`
	IncomePct   float64 `json:"income_pct"`
	ExpensesPct float64 `json:"expenses_pct"`
	SavingsPct  float64 `json:"savings_pct"`
}

func (s *Store) GetDashboardSummary(ctx context.Context, userID int64) (*DashboardSummary, error) {
	summary := &DashboardSummary{}

	// Total balance from all accounts
	if err := s.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(balance), 0) FROM accounts WHERE user_id = $1`,
		userID,
	).Scan(&summary.Balance); err != nil {
		return nil, err
	}

	// This month bounds
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	// This month income
	if err := s.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE user_id = $1 AND type = 'income' AND date >= $2`,
		userID, firstOfMonth,
	).Scan(&summary.Income); err != nil {
		return nil, err
	}

	// This month expenses
	if err := s.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE user_id = $1 AND type = 'expense' AND date >= $2`,
		userID, firstOfMonth,
	).Scan(&summary.Expenses); err != nil {
		return nil, err
	}

	net := summary.Income - summary.Expenses
	if net > 0 {
		summary.Savings = net
	}
	if summary.Income > 0 {
		summary.SavingsPct = (summary.Savings / summary.Income) * 100
	}
	return summary, nil
}
