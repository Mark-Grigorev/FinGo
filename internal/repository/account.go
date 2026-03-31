package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func (s *Store) ListAccounts(ctx context.Context, userID int64) ([]domain.Account, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, name, type, currency, balance, created_at
		 FROM accounts WHERE user_id = $1 ORDER BY created_at`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]domain.Account, 0)
	for rows.Next() {
		var a domain.Account
		if err := rows.Scan(&a.ID, &a.UserID, &a.Name, &a.Type, &a.Currency, &a.Balance, &a.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (s *Store) GetAccount(ctx context.Context, id, userID int64) (*domain.Account, error) {
	a := &domain.Account{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, name, type, currency, balance, created_at
		 FROM accounts WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&a.ID, &a.UserID, &a.Name, &a.Type, &a.Currency, &a.Balance, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return a, err
}

func (s *Store) CreateAccount(ctx context.Context, a *domain.Account) (*domain.Account, error) {
	res := &domain.Account{}
	err := s.pool.QueryRow(ctx,
		`INSERT INTO accounts (user_id, name, type, currency, balance)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, name, type, currency, balance, created_at`,
		a.UserID, a.Name, a.Type, a.Currency, a.Balance,
	).Scan(&res.ID, &res.UserID, &res.Name, &res.Type, &res.Currency, &res.Balance, &res.CreatedAt)
	return res, err
}

func (s *Store) UpdateAccount(ctx context.Context, id, userID int64, name, typ, currency string) (*domain.Account, error) {
	a := &domain.Account{}
	err := s.pool.QueryRow(ctx,
		`UPDATE accounts SET name = $1, type = $2, currency = $3
		 WHERE id = $4 AND user_id = $5
		 RETURNING id, user_id, name, type, currency, balance, created_at`,
		name, typ, currency, id, userID,
	).Scan(&a.ID, &a.UserID, &a.Name, &a.Type, &a.Currency, &a.Balance, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return a, err
}

func (s *Store) DeleteAccount(ctx context.Context, id, userID int64) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM accounts WHERE id = $1 AND user_id = $2`, id, userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
