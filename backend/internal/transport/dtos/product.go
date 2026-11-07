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

type (
	CreateProductRequest struct {
		SKU   string `json:"sku" binding:"required"`
		Name  string `json:"name" binding:"required"`
		Price int    `json:"price" binding:"required"`
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

type (
	CreateProductResponse struct {
		Product *Product `json:"product,omitempty"`
	}

	UpdateProductResponse struct {
		Product *Product `json:"product,omitempty"`
	}

	UpdateStockResponse struct {
		Product *Product `json:"product,omitempty"`
	}
)
