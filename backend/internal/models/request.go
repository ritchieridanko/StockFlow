package models

type (
	CreateProductReq struct {
		SKU   string
		Name  string
		Price int
	}

	UpdateProductReq struct {
		// Params
		ProductID int

		// Payload
		SKU   *string
		Name  *string
		Price *int
	}

	UpdateStockReq struct {
		// Params
		ProductID int

		// Payload
		UpdateType string
		Quantity   int
		Note       *string
	}
)
