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
)
