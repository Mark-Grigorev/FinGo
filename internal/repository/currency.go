package repository

import (
	"context"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func (s *Store) GetBaseCurrency(ctx context.Context, userID int64) (string, error) {
	var c string
	err := s.pool.QueryRow(ctx, `SELECT base_currency FROM users WHERE id = $1`, userID).Scan(&c)
	if err != nil { return "RUB", err }
	return c, nil
}

func (s *Store) SetBaseCurrency(ctx context.Context, userID int64, currency string) error {
	tag, err := s.pool.Exec(ctx, `UPDATE users SET base_currency = $2 WHERE id = $1`, userID, currency)
	if err != nil { return err }
	if tag.RowsAffected() == 0 { return domain.ErrNotFound }
	return nil
}

func (s *Store) ListExchangeRates(ctx context.Context, userID int64) ([]domain.ExchangeRate, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, currency, rate, updated_at FROM exchange_rates WHERE user_id = $1 ORDER BY currency`,
		userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []domain.ExchangeRate
	for rows.Next() {
		var r domain.ExchangeRate
		if err := rows.Scan(&r.ID, &r.UserID, &r.Currency, &r.Rate, &r.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	if list == nil { list = []domain.ExchangeRate{} }
	return list, rows.Err()
}

func (s *Store) UpsertExchangeRate(ctx context.Context, userID int64, currency string, rate float64) (*domain.ExchangeRate, error) {
	r := &domain.ExchangeRate{}
	err := s.pool.QueryRow(ctx, `
		INSERT INTO exchange_rates (user_id, currency, rate, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, currency) DO UPDATE SET rate = EXCLUDED.rate, updated_at = NOW()
		RETURNING id, user_id, currency, rate, updated_at`,
		userID, currency, rate,
	).Scan(&r.ID, &r.UserID, &r.Currency, &r.Rate, &r.UpdatedAt)
	return r, err
}

func (s *Store) DeleteExchangeRate(ctx context.Context, userID int64, currency string) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM exchange_rates WHERE user_id = $1 AND currency = $2`, userID, currency)
	if err != nil { return err }
	if tag.RowsAffected() == 0 { return domain.ErrNotFound }
	return nil
}

func (s *Store) GetRatesMap(ctx context.Context, userID int64) (map[string]float64, error) {
	rates := make(map[string]float64)
	rows, err := s.pool.Query(ctx,
		`SELECT currency, rate FROM exchange_rates WHERE user_id = $1`, userID)
	if err != nil { return rates, err }
	defer rows.Close()
	for rows.Next() {
		var c string; var r float64
		if err := rows.Scan(&c, &r); err != nil { return rates, err }
		rates[c] = r
	}
	return rates, rows.Err()
}
