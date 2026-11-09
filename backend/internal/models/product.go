package models

import "time"

type (
	Product struct {
		ID int

		SKU      string
		Name     string
		Price    int
		Stock    int
		IsActive bool

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	ProductListItem struct {
		ID int

		SKU      string
		Name     string
		Price    int
		Stock    int
		IsActive bool

		LastMovementAt *string

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	CreateProduct struct {
		SKU   string
		Name  string
		Price int
	}

	GetAllProducts struct {
		// Queries
		Search *string

		// Pagination
		Pagination
	}

	UpdateProduct struct {
		SKU   *string
		Name  *string
		Price *int
	}
)
