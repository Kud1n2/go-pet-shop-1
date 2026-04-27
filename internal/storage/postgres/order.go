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

func (s *Storage) PlaceOrder(ctx context.Context, userEmail string, items []models.OrderItem) (orderID int, err error) {
	const fn = "storage.postgres.order.PlaceOrder"
	//Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}
	//rollback
	defer tx.Rollback(ctx)

	var total_price int
	//Проверяем и уменьшаем stock товаров
	for _, orderItem := range items {
		var price int
		err = tx.QueryRow(ctx, `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1 RETURNING price`, orderItem.Quantity, orderItem.ProductID).Scan(&price)
		if err != nil {
			return 0, fmt.Errorf("%s:%w", fn, err)
		}
		total_price += orderItem.Quantity * price
	}
	//Получаем id пользователя
	var id int
	err = tx.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, userEmail).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}
	//Создаем order
	var order_id int
	err = tx.QueryRow(ctx, `INSERT INTO orders (user_id, total_price) VALUES ($1, $2) RETURNING id`, id, total_price).Scan(&order_id)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	//Добавляем orderItems
	for _, orderItem := range items {
		_, err := tx.Exec(ctx, `INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1,$2,$3)`, order_id, orderItem.ProductID, orderItem.Quantity)
		if err != nil {
			return 0, fmt.Errorf("%s:%w", fn, err)
		}
	}

	//Добавляем transactions
	_, err = tx.Exec(ctx, `INSERT INTO transactions(order_id, amount, status) VALUES ($1, $2, $3)`, order_id, total_price, "Completed")
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}
	return order_id, nil
}
