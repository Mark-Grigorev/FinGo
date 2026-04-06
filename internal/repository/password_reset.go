package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func (s *Store) CreatePasswordReset(ctx context.Context, token string, userID int64, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO password_resets (token, user_id, expires_at) VALUES ($1, $2, $3)`,
		token, userID, expiresAt,
	)
	return err
}

func (s *Store) GetPasswordReset(ctx context.Context, token string) (*domain.PasswordReset, error) {
	r := &domain.PasswordReset{}
	err := s.pool.QueryRow(ctx,
		`SELECT token, user_id, expires_at, used_at FROM password_resets WHERE token = $1`,
		token,
	).Scan(&r.Token, &r.UserID, &r.ExpiresAt, &r.UsedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return r, err
}

func (s *Store) MarkPasswordResetUsed(ctx context.Context, token string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE password_resets SET used_at = NOW() WHERE token = $1`,
		token,
	)
	return err
}
