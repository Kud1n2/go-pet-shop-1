package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
	"time"
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

func (s *Storage) GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error) {
	const fn = "storage.postgres.users.GetUserOrderHistory"

	var orderDetails []models.OrderDetail

	rows, err := s.db.Query(ctx, `select o.id, o.total_price, o.created_at, p.name, oi.quantity, t.status  from order_items oi
								join orders o on o.id = oi.order_id
								join transactions t on o.id = t.order_id
								join users u ON o.user_id = u.id 
								join products p on oi.product_id = p.id
								WHERE u.email = $1`, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	OrderDetailsMap := make(map[int]*models.OrderDetail)

	for rows.Next() {
		var (
			OrderID           int
			TotalPrice        int
			CreatedAt         time.Time
			TransactionStatus string
			Items             models.OrderDetailItems
		)
		if err = rows.Scan(&OrderID, &TotalPrice, &CreatedAt, &Items.ProductName, &Items.Quantity, &TransactionStatus); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		order, ok := OrderDetailsMap[OrderID]
		if !ok {
			order = &models.OrderDetail{
				OrderID:           OrderID,
				TotalPrice:        TotalPrice,
				CreatedAt:         CreatedAt,
				TransactionStatus: TransactionStatus,
			}
			OrderDetailsMap[OrderID] = order
		}
		order.Items = append(order.Items, Items)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	for _, value := range OrderDetailsMap {
		orderDetails = append(orderDetails, *value)
	}
	return orderDetails, nil
}
