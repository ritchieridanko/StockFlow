package models

type (
	CreateProductReq struct {
		SKU   string
		Name  string
		Price int
	}

	GetAllProductsReq struct {
		// Queries
		Search *string

		// Pagination
		Pagination
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

type (
	OrderItemReq struct {
		SKU      string
		Quantity int
	}

	CreateOrderReq struct {
		Customer string
		Items    []OrderItemReq
	}
)
