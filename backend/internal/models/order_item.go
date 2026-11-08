package models

import "time"

type (
	OrderItem struct {
		ID      int
		OrderID int

		SKU       string
		Quantity  int
		UnitPrice int

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	CreateOrderItem struct {
		SKU       string
		Quantity  int
		UnitPrice int
	}

	CreateOrderItems struct {
		OrderID int
		Items   []CreateOrderItem
	}
)
