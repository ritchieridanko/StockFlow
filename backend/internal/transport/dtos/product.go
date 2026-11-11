package dtos

import "time"

type Product struct {
	ID int `json:"id"`

	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Stock    int    `json:"stock"`
	IsActive bool   `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductListItem struct {
	ID int `json:"id"`

	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Stock    int    `json:"stock"`
	IsActive bool   `json:"is_active"`

	LastMovementAt *string `json:"last_movement_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StockMovement struct {
	ID        int `json:"id"`
	ProductID int `json:"product_id"`

	UpdateType string  `json:"type"`
	Quantity   int     `json:"qty"`
	Note       *string `json:"note"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type (
	CreateProductRequest struct {
		SKU   string `json:"sku" binding:"required"`
		Name  string `json:"name" binding:"required"`
		Price int    `json:"price" binding:"required"`
	}

	GetAllProductsRequest struct {
		// Queries
		Search *string `form:"search"`

		// Pagination
		PaginationParams
	}

	UpdateProductRequest struct {
		SKU   *string `json:"sku"`
		Name  *string `json:"name"`
		Price *int    `json:"price"`
	}

	UpdateStockRequest struct {
		UpdateType string  `json:"type" binding:"required"`
		Quantity   int     `json:"qty" binding:"required"`
		Note       *string `json:"note"`
	}
)
