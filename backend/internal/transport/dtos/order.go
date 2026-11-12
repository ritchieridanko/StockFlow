package dtos

import "time"

type Order struct {
	ID int `json:"id"`

	Customer string `json:"customer"`
	Status   string `json:"status"`
	Total    int    `json:"total"`

	Items []OrderItem `json:"items"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CancelledAt *time.Time `json:"cancelled_at"`
}

type OrderItem struct {
	ID      int `json:"id"`
	OrderID int `json:"order_id"`

	SKU       string `json:"sku"`
	Quantity  int    `json:"quantity"`
	UnitPrice int    `json:"unit_price"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type (
	OrderItemRequest struct {
		SKU      string `json:"sku" binding:"required"`
		Quantity int    `json:"quantity" binding:"required"`
	}

	CreateOrderRequest struct {
		Customer string             `json:"customer" binding:"required"`
		Items    []OrderItemRequest `json:"items"`
	}

	GetAllOrdersRequest struct {
		// Queries
		Search *string `form:"search"`

		// Pagination
		PaginationParams
	}
)
