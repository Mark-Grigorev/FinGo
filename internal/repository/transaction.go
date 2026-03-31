package repository

import (
	"context"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/jackc/pgx/v5"
)

type TransactionFilter struct {
	Page       int
	Limit      int
	CategoryID int64
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
	if f.CategoryID > 0 {
		if err := s.pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM transactions WHERE user_id = $1 AND category_id = $2`, userID, f.CategoryID,
		).Scan(&total); err != nil {
			return nil, 0, err
		}
	} else {
		if err := s.pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM transactions WHERE user_id = $1`, userID,
		).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	var rows pgx.Rows
	var err error
	if f.CategoryID > 0 {
		rows, err = s.pool.Query(ctx, `
			SELECT t.id, t.user_id, t.account_id,
			       COALESCE(a.name, ''),
			       t.category_id,
			       COALESCE(c.name, ''), COALESCE(c.color, ''), COALESCE(c.icon, '💳'),
			       t.type, t.amount, t.name, t.date, t.created_at
			FROM transactions t
			LEFT JOIN categories c ON c.id = t.category_id
			LEFT JOIN accounts a ON a.id = t.account_id
			WHERE t.user_id = $1 AND t.category_id = $2
			ORDER BY t.date DESC, t.created_at DESC
			LIMIT $3 OFFSET $4`,
			userID, f.CategoryID, f.Limit, offset,
		)
	} else {
		rows, err = s.pool.Query(ctx, `
			SELECT t.id, t.user_id, t.account_id,
			       COALESCE(a.name, ''),
			       t.category_id,
			       COALESCE(c.name, ''), COALESCE(c.color, ''), COALESCE(c.icon, '💳'),
			       t.type, t.amount, t.name, t.date, t.created_at
			FROM transactions t
			LEFT JOIN categories c ON c.id = t.category_id
			LEFT JOIN accounts a ON a.id = t.account_id
			WHERE t.user_id = $1
			ORDER BY t.date DESC, t.created_at DESC
			LIMIT $2 OFFSET $3`,
			userID, f.Limit, offset,
		)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := make([]domain.Transaction, 0)
	for rows.Next() {
		var tx domain.Transaction
		if err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.AccountID, &tx.AccountName,
			&tx.CategoryID, &tx.CategoryName, &tx.CategoryColor, &tx.Icon,
			&tx.Type, &tx.Amount, &tx.Name, &tx.Date, &tx.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, tx)
	}
	return list, total, rows.Err()
}

func (s *Store) DeleteTransaction(ctx context.Context, id, userID int64) error {
	dbTx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer dbTx.Rollback(ctx)

	// Read the transaction first to reverse the balance
	var txType string
	var amount float64
	var accountID int64
	err = dbTx.QueryRow(ctx,
		`DELETE FROM transactions WHERE id=$1 AND user_id=$2 RETURNING type, amount, account_id`,
		id, userID,
	).Scan(&txType, &amount, &accountID)
	if err != nil {
		return domain.ErrNotFound
	}

	// Reverse the balance change
	delta := -amount
	if txType == "expense" {
		delta = amount
	}
	if _, err = dbTx.Exec(ctx,
		`UPDATE accounts SET balance = balance + $1 WHERE id = $2 AND user_id = $3`,
		delta, accountID, userID,
	); err != nil {
		return err
	}

	return dbTx.Commit(ctx)
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
