package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func (s *Store) ListBudgets(ctx context.Context, userID int64, month time.Time) ([]domain.Budget, error) {
	rows, err := s.pool.Query(ctx, `
        SELECT b.id, b.user_id, b.category_id,
               COALESCE(c.name, '') AS name,
               COALESCE(c.color, '#888888') AS color,
               TO_CHAR(b.month, 'YYYY-MM-DD') AS month,
               b.limit_amount,
               COALESCE(SUM(t.amount), 0) AS spent
        FROM budgets b
        LEFT JOIN categories c ON c.id = b.category_id
        LEFT JOIN transactions t
               ON t.category_id = b.category_id
              AND t.user_id = b.user_id
              AND t.type = 'expense'
              AND DATE_TRUNC('month', t.date) = DATE_TRUNC('month', b.month)
        WHERE b.user_id = $1
          AND DATE_TRUNC('month', b.month) = DATE_TRUNC('month', $2::date)
        GROUP BY b.id, b.user_id, b.category_id, c.name, c.color, b.month, b.limit_amount
        ORDER BY c.name`,
		userID, month,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Budget
	for rows.Next() {
		var b domain.Budget
		if err := rows.Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Name, &b.Color, &b.Month, &b.Limit, &b.Spent); err != nil {
			return nil, err
		}
		if b.Limit > 0 {
			b.Pct = (b.Spent / b.Limit) * 100
		}
		list = append(list, b)
	}
	if list == nil {
		list = []domain.Budget{}
	}
	return list, rows.Err()
}

func (s *Store) CreateBudget(ctx context.Context, b *domain.Budget) (*domain.Budget, error) {
	out := &domain.Budget{}
	err := s.pool.QueryRow(ctx, `
        INSERT INTO budgets (user_id, category_id, month, limit_amount)
        VALUES ($1, $2, DATE_TRUNC('month', $3::date), $4)
        ON CONFLICT (user_id, category_id, month) DO UPDATE SET limit_amount = EXCLUDED.limit_amount
        RETURNING id, user_id, category_id, TO_CHAR(month, 'YYYY-MM-DD'), limit_amount`,
		b.UserID, b.CategoryID, b.Month, b.Limit,
	).Scan(&out.ID, &out.UserID, &out.CategoryID, &out.Month, &out.Limit)
	if err != nil {
		return nil, err
	}
	out.Name = b.Name
	return out, nil
}

func (s *Store) UpdateBudget(ctx context.Context, id, userID int64, limit float64) (*domain.Budget, error) {
	out := &domain.Budget{}
	err := s.pool.QueryRow(ctx, `
        UPDATE budgets SET limit_amount = $3
        WHERE id = $1 AND user_id = $2
        RETURNING id, user_id, category_id, TO_CHAR(month, 'YYYY-MM-DD'), limit_amount`,
		id, userID, limit,
	).Scan(&out.ID, &out.UserID, &out.CategoryID, &out.Month, &out.Limit)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return out, err
}

func (s *Store) DeleteBudget(ctx context.Context, id, userID int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM budgets WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
