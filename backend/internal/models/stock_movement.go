package models

import "time"

type (
	StockMovement struct {
		ID        int
		ProductID int

		UpdateType string
		Quantity   int
		Note       *string

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	CreateStockMovement struct {
		ProductID  int
		UpdateType string
		Quantity   int
		Note       *string
	}
)
