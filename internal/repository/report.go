package repository

import (
	"context"
	"time"
)

type BarDataPoint struct {
	Label    string  `json:"label"`
	Income   float64 `json:"income"`
	Expenses float64 `json:"expenses"`
}

type PieCategoryData struct {
	Name   string  `json:"name"`
	Color  string  `json:"color"`
	Amount float64 `json:"amount"`
}

type ComparisonData struct {
	Income      float64 `json:"income"`
	IncomePct   float64 `json:"income_pct"`
	Expenses    float64 `json:"expenses"`
	ExpensesPct float64 `json:"expenses_pct"`
	Savings     float64 `json:"savings"`
	SavingsPct  float64 `json:"savings_pct"`
}

type ReportResult struct {
	BarData       []BarDataPoint    `json:"bar_data"`
	PieData       []PieCategoryData `json:"pie_data"`
	Comparison    ComparisonData    `json:"comparison"`
	TopCategories []PieCategoryData `json:"top_categories"`
}

func (s *Store) GetReport(ctx context.Context, userID int64, from, to time.Time) (*ReportResult, error) {
	result := &ReportResult{
		BarData:       []BarDataPoint{},
		PieData:       []PieCategoryData{},
		TopCategories: []PieCategoryData{},
	}

	// Bar chart — monthly breakdown
	barRows, err := s.pool.Query(ctx, `
		SELECT
			TO_CHAR(date_trunc('month', date), 'Mon YYYY') AS label,
			date_trunc('month', date) AS month,
			COALESCE(SUM(CASE WHEN type='income'  THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type='expense' THEN amount ELSE 0 END), 0) AS expenses
		FROM transactions
		WHERE user_id = $1 AND date BETWEEN $2 AND $3
		GROUP BY month, label
		ORDER BY month`,
		userID, from, to,
	)
	if err != nil {
		return nil, err
	}
	defer barRows.Close()
	for barRows.Next() {
		var p BarDataPoint
		var month time.Time
		if err := barRows.Scan(&p.Label, &month, &p.Income, &p.Expenses); err != nil {
			return nil, err
		}
		result.BarData = append(result.BarData, p)
	}
	if err := barRows.Err(); err != nil {
		return nil, err
	}

	// Pie chart — expenses by category
	pieRows, err := s.pool.Query(ctx, `
		SELECT
			COALESCE(c.name, 'Без категории') AS name,
			COALESCE(c.color, '#888888')       AS color,
			SUM(t.amount)                       AS amount
		FROM transactions t
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE t.user_id = $1 AND t.type = 'expense' AND t.date BETWEEN $2 AND $3
		GROUP BY c.name, c.color
		ORDER BY amount DESC`,
		userID, from, to,
	)
	if err != nil {
		return nil, err
	}
	defer pieRows.Close()
	for pieRows.Next() {
		var p PieCategoryData
		if err := pieRows.Scan(&p.Name, &p.Color, &p.Amount); err != nil {
			return nil, err
		}
		result.PieData = append(result.PieData, p)
	}
	if err := pieRows.Err(); err != nil {
		return nil, err
	}

	// Top 5 categories
	if len(result.PieData) > 5 {
		result.TopCategories = result.PieData[:5]
	} else {
		result.TopCategories = result.PieData
	}

	// Comparison — current period vs previous period of same length
	duration := to.Sub(from)
	prevTo := from.Add(-24 * time.Hour)
	prevFrom := prevTo.Add(-duration)

	var curIncome, curExpenses, prevIncome, prevExpenses float64

	if err := s.pool.QueryRow(ctx,
		`SELECT
			COALESCE(SUM(CASE WHEN type='income'  THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type='expense' THEN amount ELSE 0 END), 0)
		 FROM transactions WHERE user_id=$1 AND date BETWEEN $2 AND $3`,
		userID, from, to,
	).Scan(&curIncome, &curExpenses); err != nil {
		return nil, err
	}

	if err := s.pool.QueryRow(ctx,
		`SELECT
			COALESCE(SUM(CASE WHEN type='income'  THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type='expense' THEN amount ELSE 0 END), 0)
		 FROM transactions WHERE user_id=$1 AND date BETWEEN $2 AND $3`,
		userID, prevFrom, prevTo,
	).Scan(&prevIncome, &prevExpenses); err != nil {
		return nil, err
	}

	result.Comparison = ComparisonData{
		Income:   curIncome,
		Expenses: curExpenses,
		Savings:  max(0, curIncome-curExpenses),
	}
	if prevIncome > 0 {
		result.Comparison.IncomePct = ((curIncome - prevIncome) / prevIncome) * 100
	}
	if prevExpenses > 0 {
		result.Comparison.ExpensesPct = ((curExpenses - prevExpenses) / prevExpenses) * 100
	}
	prevSavings := max(0, prevIncome-prevExpenses)
	if prevSavings > 0 {
		result.Comparison.SavingsPct = ((result.Comparison.Savings - prevSavings) / prevSavings) * 100
	}

	return result, nil
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
