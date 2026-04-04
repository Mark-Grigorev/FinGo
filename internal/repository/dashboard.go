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

	// Base currency
	base, _ := s.GetBaseCurrency(ctx, userID)
	if base == "" { base = "RUB" }

	// Exchange rates map
	rates, _ := s.GetRatesMap(ctx, userID)

	// Sum balances converting to base currency
	rows, err := s.pool.Query(ctx,
		`SELECT currency, COALESCE(SUM(balance), 0) FROM accounts WHERE user_id = $1 GROUP BY currency`,
		userID)
	if err != nil { return nil, err }
	defer rows.Close()
	for rows.Next() {
		var curr string; var bal float64
		if err := rows.Scan(&curr, &bal); err != nil { return nil, err }
		summary.Balance += convertAmount(bal, curr, base, rates)
	}

	// This month income/expenses (keep as-is for now, converted below)
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	incomeRows, err := s.pool.Query(ctx, `
		SELECT COALESCE(SUM(t.amount), 0), a.currency
		FROM transactions t
		JOIN accounts a ON a.id = t.account_id
		WHERE t.user_id = $1 AND t.type = 'income' AND t.date >= $2
		GROUP BY a.currency`,
		userID, firstOfMonth)
	if err != nil { return nil, err }
	defer incomeRows.Close()
	for incomeRows.Next() {
		var amt float64; var curr string
		if err := incomeRows.Scan(&amt, &curr); err != nil { return nil, err }
		summary.Income += convertAmount(amt, curr, base, rates)
	}

	expRows, err := s.pool.Query(ctx, `
		SELECT COALESCE(SUM(t.amount), 0), a.currency
		FROM transactions t
		JOIN accounts a ON a.id = t.account_id
		WHERE t.user_id = $1 AND t.type = 'expense' AND t.date >= $2
		GROUP BY a.currency`,
		userID, firstOfMonth)
	if err != nil { return nil, err }
	defer expRows.Close()
	for expRows.Next() {
		var amt float64; var curr string
		if err := expRows.Scan(&amt, &curr); err != nil { return nil, err }
		summary.Expenses += convertAmount(amt, curr, base, rates)
	}

	net := summary.Income - summary.Expenses
	if net > 0 { summary.Savings = net }
	if summary.Income > 0 { summary.SavingsPct = (summary.Savings / summary.Income) * 100 }
	return summary, nil
}

// convertAmount converts amount from srcCurrency to dstCurrency using rates map.
// If src == dst or no rate available, returns amount as-is.
func convertAmount(amount float64, src, dst string, rates map[string]float64) float64 {
	if src == dst || amount == 0 { return amount }
	rate, ok := rates[src]
	if !ok { return amount } // no rate — return as-is
	return amount * rate
}
