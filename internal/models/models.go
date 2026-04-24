package models

import "time"

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int // количество на складе
}

type Customer struct {
	ID    int
	Name  string
	Email string
}

type Order struct {
	ID         int
	CustomerID int
	TotalPrice int
	CreatedAt  time.Time
}

type OrderItem struct {
	ID        int
	OrderID   int
	ProductID int
	Quantity  int
}

type OrderItemsList struct {
	Email      string
	OrderItems []OrderItem
}

type OrderDetail struct {
	Order
	OrderItems        []OrderItem
	TransactionStatus string
}

type PopularProduct struct {
	Product
	Quantity int
}
