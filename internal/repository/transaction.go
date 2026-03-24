package repository

import (
	"context"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

type TransactionFilter struct {
	Page  int
	Limit int
}

func (s *Store) ListTransactions(ctx context.Context, userID int64, f TransactionFilter) ([]domain.Transaction, int, error) {
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Page < 1 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	var total int
	if err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM transactions WHERE user_id = $1`, userID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx, `
		SELECT t.id, t.user_id, t.account_id, t.category_id,
		       COALESCE(c.name, ''), COALESCE(c.icon, '💳'),
		       t.type, t.amount, t.name, t.date, t.created_at
		FROM transactions t
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE t.user_id = $1
		ORDER BY t.date DESC, t.created_at DESC
		LIMIT $2 OFFSET $3`,
		userID, f.Limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := make([]domain.Transaction, 0)
	for rows.Next() {
		var tx domain.Transaction
		if err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.AccountID, &tx.CategoryID,
			&tx.CategoryName, &tx.Icon,
			&tx.Type, &tx.Amount, &tx.Name, &tx.Date, &tx.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, tx)
	}
	return list, total, rows.Err()
}

func (s *Store) CreateTransaction(ctx context.Context, t *domain.Transaction) (*domain.Transaction, error) {
	dbTx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer dbTx.Rollback(ctx)

	date := t.Date
	if date.IsZero() {
		date = time.Now()
	}

	result := &domain.Transaction{}
	err = dbTx.QueryRow(ctx, `
		INSERT INTO transactions (user_id, account_id, category_id, type, amount, name, date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, account_id, category_id, type, amount, name, date, created_at`,
		t.UserID, t.AccountID, t.CategoryID, t.Type, t.Amount, t.Name, date,
	).Scan(
		&result.ID, &result.UserID, &result.AccountID, &result.CategoryID,
		&result.Type, &result.Amount, &result.Name, &result.Date, &result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Update account balance atomically
	delta := t.Amount
	if t.Type == "expense" {
		delta = -t.Amount
	}
	if _, err = dbTx.Exec(ctx,
		`UPDATE accounts SET balance = balance + $1 WHERE id = $2 AND user_id = $3`,
		delta, t.AccountID, t.UserID,
	); err != nil {
		return nil, err
	}

	return result, dbTx.Commit(ctx)
}
