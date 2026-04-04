package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

type TransactionFilter struct {
	Page       int
	Limit      int
	CategoryID int64
	From       time.Time
	To         time.Time
}

func (s *Store) ListTransactions(ctx context.Context, userID int64, f TransactionFilter) ([]domain.Transaction, int, error) {
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Page < 1 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	// Build dynamic WHERE clause
	args := []any{userID}
	conds := []string{"t.user_id = $1"}
	n := 2
	if f.CategoryID > 0 {
		conds = append(conds, fmt.Sprintf("t.category_id = $%d", n))
		args = append(args, f.CategoryID)
		n++
	}
	if !f.From.IsZero() {
		conds = append(conds, fmt.Sprintf("t.date >= $%d", n))
		args = append(args, f.From)
		n++
	}
	if !f.To.IsZero() {
		conds = append(conds, fmt.Sprintf("t.date <= $%d", n))
		args = append(args, f.To)
		n++
	}
	where := strings.Join(conds, " AND ")

	var total int
	if err := s.pool.QueryRow(ctx,
		fmt.Sprintf(`SELECT COUNT(*) FROM transactions t WHERE %s`, where),
		args...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Add pagination params
	limitArg := n
	offsetArg := n + 1
	args = append(args, f.Limit, offset)

	rows, err := s.pool.Query(ctx, fmt.Sprintf(`
		SELECT t.id, t.user_id, t.account_id,
		       COALESCE(a.name, ''),
		       t.category_id,
		       COALESCE(c.name, ''), COALESCE(c.color, ''), COALESCE(c.icon, '💳'),
		       t.type, t.amount, t.name, t.date, t.created_at
		FROM transactions t
		LEFT JOIN categories c ON c.id = t.category_id
		LEFT JOIN accounts a ON a.id = t.account_id
		WHERE %s
		ORDER BY t.date DESC, t.created_at DESC
		LIMIT $%d OFFSET $%d`, where, limitArg, offsetArg),
		args...,
	)
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

func (s *Store) ExportTransactions(ctx context.Context, userID int64, from, to time.Time) ([]domain.Transaction, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT t.id, t.user_id, t.account_id,
		       COALESCE(a.name, ''),
		       t.category_id,
		       COALESCE(c.name, ''), COALESCE(c.color, ''), COALESCE(c.icon, '💳'),
		       t.type, t.amount, t.name, t.date, t.created_at
		FROM transactions t
		LEFT JOIN categories c ON c.id = t.category_id
		LEFT JOIN accounts a ON a.id = t.account_id
		WHERE t.user_id = $1 AND t.date BETWEEN $2 AND $3
		ORDER BY t.date DESC, t.created_at DESC`,
		userID, from, to,
	)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		list = append(list, tx)
	}
	return list, rows.Err()
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
