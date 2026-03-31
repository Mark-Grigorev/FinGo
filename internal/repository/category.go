package repository

import (
	"context"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func (s *Store) ListCategories(ctx context.Context, userID int64) ([]domain.Category, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, COALESCE(user_id, 0), name, icon, color, type, is_system
		 FROM categories
		 WHERE user_id = $1 OR user_id IS NULL
		 ORDER BY is_system DESC, name`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Icon, &c.Color, &c.Type, &c.IsSystem); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if categories == nil {
		categories = []domain.Category{}
	}
	return categories, rows.Err()
}

func (s *Store) CreateCategory(ctx context.Context, c *domain.Category) (*domain.Category, error) {
	err := s.pool.QueryRow(ctx,
		`INSERT INTO categories (user_id, name, icon, color, type)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, name, icon, color, type, is_system`,
		c.UserID, c.Name, c.Icon, c.Color, c.Type,
	).Scan(&c.ID, &c.UserID, &c.Name, &c.Icon, &c.Color, &c.Type, &c.IsSystem)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Store) UpdateCategory(ctx context.Context, id, userID int64, name, icon, color string) (*domain.Category, error) {
	c := &domain.Category{}
	err := s.pool.QueryRow(ctx,
		`UPDATE categories
		 SET name=$3, icon=$4, color=$5
		 WHERE id=$1 AND user_id=$2 AND is_system=FALSE
		 RETURNING id, user_id, name, icon, color, type, is_system`,
		id, userID, name, icon, color,
	).Scan(&c.ID, &c.UserID, &c.Name, &c.Icon, &c.Color, &c.Type, &c.IsSystem)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return c, nil
}

func (s *Store) DeleteCategory(ctx context.Context, id, userID int64) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM categories WHERE id = $1 AND user_id = $2 AND is_system = FALSE`,
		id, userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
