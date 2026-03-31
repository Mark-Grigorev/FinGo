package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func (s *Store) ListRecurring(ctx context.Context, userID int64) ([]domain.RecurringPayment, error) {
	rows, err := s.pool.Query(ctx, `
        SELECT r.id, r.user_id, r.account_id, r.category_id,
               COALESCE(c.name, '') AS category_name,
               r.name, r.amount, r.frequency, r.next_payment_date, r.is_active, r.created_at
        FROM recurring_payments r
        LEFT JOIN categories c ON c.id = r.category_id
        WHERE r.user_id = $1 AND r.is_active = TRUE
        ORDER BY r.next_payment_date ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.RecurringPayment
	for rows.Next() {
		var r domain.RecurringPayment
		if err := rows.Scan(&r.ID, &r.UserID, &r.AccountID, &r.CategoryID,
			&r.CategoryName, &r.Name, &r.Amount, &r.Frequency,
			&r.NextPaymentDate, &r.IsActive, &r.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	if list == nil {
		list = []domain.RecurringPayment{}
	}
	return list, rows.Err()
}

func (s *Store) CreateRecurring(ctx context.Context, r *domain.RecurringPayment) (*domain.RecurringPayment, error) {
	out := &domain.RecurringPayment{}
	err := s.pool.QueryRow(ctx, `
        INSERT INTO recurring_payments (user_id, account_id, category_id, name, amount, frequency, next_payment_date)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, user_id, account_id, category_id, name, amount, frequency, next_payment_date, is_active, created_at`,
		r.UserID, r.AccountID, r.CategoryID, r.Name, r.Amount, r.Frequency, r.NextPaymentDate,
	).Scan(&out.ID, &out.UserID, &out.AccountID, &out.CategoryID,
		&out.Name, &out.Amount, &out.Frequency, &out.NextPaymentDate, &out.IsActive, &out.CreatedAt)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) UpdateRecurring(ctx context.Context, id, userID int64, name string, amount float64, frequency string, nextDate time.Time, accountID int64, categoryID *int64) (*domain.RecurringPayment, error) {
	out := &domain.RecurringPayment{}
	err := s.pool.QueryRow(ctx, `
        UPDATE recurring_payments
        SET name=$3, amount=$4, frequency=$5, next_payment_date=$6, account_id=$7, category_id=$8
        WHERE id=$1 AND user_id=$2
        RETURNING id, user_id, account_id, category_id, name, amount, frequency, next_payment_date, is_active, created_at`,
		id, userID, name, amount, frequency, nextDate, accountID, categoryID,
	).Scan(&out.ID, &out.UserID, &out.AccountID, &out.CategoryID,
		&out.Name, &out.Amount, &out.Frequency, &out.NextPaymentDate, &out.IsActive, &out.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return out, err
}

func (s *Store) DeleteRecurring(ctx context.Context, id, userID int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM recurring_payments WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
