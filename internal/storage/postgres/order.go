package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
)

func (s *Storage) CreateOrder(ctx context.Context, order models.Order) (int, error) {
	const fn = "storage.postgres.order.CreateOrder"

	var id int
	err := s.db.QueryRow(ctx, `INSERT INTO orders (user_id, total_price) VALUES ($1, $2) RETURNING id`, order.CustomerID, order.TotalPrice).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}
	return id, nil
}

func (s *Storage) AddOrderItem(ctx context.Context, orderItem models.OrderItem) error {
	const fn = "storage.postgres.order.AddOrderItem"

	_, err := s.db.Exec(ctx, `INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity)

	if err != nil {
		return fmt.Errorf("%s:%w", fn, err)
	}

	return nil
}

func (s *Storage) GetOrderByID(ctx context.Context, id int) (models.Order, error) {
	const fn = "storage.order.GetOrderByID"

	row := s.db.QueryRow(ctx, `SELECT id, user_id, total_price, created_at FROM orders WHERE id = $1`, id)

	var o models.Order

	if err := row.Scan(&o.ID, &o.CustomerID, &o.TotalPrice, &o.CreatedAt); err != nil {
		return models.Order{}, fmt.Errorf("%s:%w", fn, err)
	}

	return o, nil
}

func (s *Storage) GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error) {
	const fn = "storage.order.GetOrdersByUserEmail"

	rows, err := s.db.Query(ctx, `SELECT o.id, user_id, total_price, created_at FROM orders o 
	JOIN users u ON o.user_id = u.id WHERE u.email = $1`, email)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", fn, err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.TotalPrice, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s:%w", fn, err)
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s:%w", fn, err)
	}

	return orders, nil
}

func (s *Storage) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	const fn = "storage.postgres.order.GetOrderItemsByOrderID"

	rows, err := s.db.Query(ctx, `SELECT id, order_id, product_id, quantity FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", fn, err)
	}
	defer rows.Close()

	var order_items []models.OrderItem

	for rows.Next() {
		var o_i models.OrderItem
		if err := rows.Scan(&o_i.ID, &o_i.OrderID, &o_i.ProductID, &o_i.Quantity); err != nil {
			return nil, fmt.Errorf("%s:%w", fn, err)
		}

		order_items = append(order_items, o_i)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s:%w", fn, err)
	}

	return order_items, nil
}
