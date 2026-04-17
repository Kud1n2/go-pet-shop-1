package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
)

func (s *Storage) GetAllUsers(ctx context.Context) ([]models.Customer, error) {
	const fn = "storage.postgres.user.GetAllUsers"

	rows, err := s.db.Query(ctx, `SELECT id, name, email FROM users`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var users []models.Customer

	for rows.Next() {
		var u models.Customer
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, fmt.Errorf("%s:%w", fn, err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s:%w", fn, err)
	}

	return users, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.Customer, error) {
	const fn = "storage.postgres.user.GetUserByEmail"

	row := s.db.QueryRow(ctx, `SELECT id, name, email FROM users WHERE email = $1`, email)

	var user models.Customer
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return models.Customer{}, fmt.Errorf("%s: %w: Email=%s", fn, ErrNotFound, user.Email)
	}
	return user, nil
}

func (s *Storage) CreateUser(ctx context.Context, user models.Customer) error {
	const fn = "storage.postgres.users.CreateUser"

	_, err := s.db.Exec(ctx, `INSERT INTO users (name, email) VALUES ($1, $2)`, user.Name, user.Email)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}
