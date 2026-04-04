package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, email, name, password_hash, base_currency, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.BaseCurrency, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (s *Store) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	u := &domain.User{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, email, name, password_hash, base_currency, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.BaseCurrency, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (s *Store) UpdateUser(ctx context.Context, id int64, name, email string) (*domain.User, error) {
	u := &domain.User{}
	err := s.pool.QueryRow(ctx,
		`UPDATE users SET name=$2, email=$3 WHERE id=$1 RETURNING id, email, name, base_currency, created_at`,
		id, name, email,
	).Scan(&u.ID, &u.Email, &u.Name, &u.BaseCurrency, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrAlreadyExists
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (s *Store) UpdatePassword(ctx context.Context, id int64, hash string) error {
	tag, err := s.pool.Exec(ctx,
		`UPDATE users SET password_hash=$2 WHERE id=$1`,
		id, hash,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) CreateUser(ctx context.Context, email, name, hash string) (*domain.User, error) {
	u := &domain.User{}
	err := s.pool.QueryRow(ctx,
		`INSERT INTO users (email, name, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id, email, name, base_currency, created_at`,
		email, name, hash,
	).Scan(&u.ID, &u.Email, &u.Name, &u.BaseCurrency, &u.CreatedAt)
	return u, err
}
