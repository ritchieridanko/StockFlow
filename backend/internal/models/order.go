package models

import "time"

type (
	Order struct {
		ID int

		Customer string
		Status   string
		Total    int

		Items []OrderItem

		CreatedAt   time.Time
		UpdatedAt   time.Time
		CancelledAt *time.Time
	}

	CreateOrder struct {
		Customer string
		Total    int
	}

	GetAllOrders struct {
		// Queries
		Search *string

		// Pagination
		Pagination
	}
)
