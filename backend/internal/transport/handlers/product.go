package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/transport/dtos"
	"github.com/ritchieridanko/stockflow/backend/internal/usecases"
	"github.com/ritchieridanko/stockflow/backend/internal/utils"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
)

type ProductHandler struct {
	pu usecases.ProductUsecase
}

func NewProductHandler(pu usecases.ProductUsecase) *ProductHandler {
	return &ProductHandler{pu: pu}
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var payload dtos.CreateProductRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(
			ce.CodeInvalidPayload,
			ce.MsgInvalidPayload,
			fmt.Errorf("failed to bind json payload: %w", err),
		).Bind(
			ctx,
		)
		return
	}

	p, err := h.pu.CreateProduct(
		ctx.Request.Context(),
		&models.CreateProductReq{
			SKU:   payload.SKU,
			Name:  payload.Name,
			Price: payload.Price,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetResponse(
		ctx,
		http.StatusCreated,
		"Product created successfully",
		dtos.CreateProductResponse{
			Product: h.toProduct(p),
		},
		nil,
	)
}

func (h *ProductHandler) GetAllProducts(ctx *gin.Context) {
	var params dtos.GetAllProductsRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ce.NewError(
			ce.CodeInvalidParam,
			ce.MsgInvalidParam,
			fmt.Errorf("failed to bind query params: %w", err),
		).Bind(
			ctx,
		)
		return
	}

	plis, total, err := h.pu.GetAllProducts(
		ctx.Request.Context(),
		&models.GetAllProductsReq{
			Search: params.Search,

			Pagination: models.Pagination{
				Page: params.Page,
				Size: params.Limit,
			},
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	products := make([]dtos.ProductListItem, 0, params.Limit)
	for _, pli := range plis {
		products = append(products, *h.toProductListItem(&pli))
	}

	utils.SetResponse(
		ctx,
		http.StatusOK,
		"Products retrieved successfully",
		dtos.GetAllProductsResponse{
			Products: products,
		},
		&utils.Pagination{
			Page:  params.Page,
			Limit: params.Limit,
			Total: total,
		},
	)
}

func (h *ProductHandler) GetAllProductStockMovements(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ce.NewError(
			ce.CodeInvalidParam,
			ce.MsgInvalidParam,
			fmt.Errorf("failed to convert url param: %w", err),
		).Bind(
			ctx,
		)
		return
	}

	sms, getErr := h.pu.GetAllProductStockMovements(
		ctx.Request.Context(),
		productID,
	)
	if getErr != nil {
		getErr.Bind(ctx)
		return
	}

	movements := make([]dtos.StockMovement, 0, len(sms))
	for _, sm := range sms {
		movements = append(movements, *h.toStockMovement(&sm))
	}

	utils.SetResponse(
		ctx,
		http.StatusOK,
		"Stock movements retrieved successfully",
		dtos.GetAllProductStockMovementsResponse{
			StockMovements: movements,
		},
		nil,
	)
}

func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ce.NewError(
			ce.CodeInvalidParam,
			ce.MsgInvalidParam,
			fmt.Errorf("failed to convert url param: %w", err),
		).Bind(
			ctx,
		)
		return
	}

	var payload dtos.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(
			ce.CodeInvalidPayload,
			ce.MsgInvalidPayload,
			fmt.Errorf("failed to bind json payload: %w", err),
		).Bind(
			ctx,
		)
		return
	}

	p, updateErr := h.pu.UpdateProduct(
		ctx.Request.Context(),
		&models.UpdateProductReq{
			ProductID: productID,

			SKU:   payload.SKU,
			Name:  payload.Name,
			Price: payload.Price,
		},
	)
	if updateErr != nil {
		updateErr.Bind(ctx)
		return
	}

	utils.SetResponse(
		ctx,
		http.StatusOK,
		"Product updated successfully",
		dtos.UpdateProductResponse{
			Product: h.toProduct(p),
		},
		nil,
	)
}

func (h *ProductHandler) UpdateStock(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ce.NewError(
			ce.CodeInvalidParam,
			ce.MsgInvalidParam,
			fmt.Errorf("failed to convert url param: %w", err),
		).Bind(
			ctx,
		)
		return
	}

	var payload dtos.UpdateStockRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(
			ce.CodeInvalidPayload,
			ce.MsgInvalidPayload,
			fmt.Errorf("failed to bind json payload: %w", err),
		).Bind(
			ctx,
		)
		return
	}

	p, updateErr := h.pu.UpdateStock(
		ctx.Request.Context(),
		&models.UpdateStockReq{
			ProductID: productID,

			UpdateType: payload.UpdateType,
			Quantity:   payload.Quantity,
			Note:       payload.Note,
		},
	)
	if updateErr != nil {
		updateErr.Bind(ctx)
		return
	}

	utils.SetResponse(
		ctx,
		http.StatusOK,
		"Stock updated successfully",
		dtos.UpdateStockResponse{
			Product: h.toProduct(p),
		},
		nil,
	)
}

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ce.NewError(
			ce.CodeInvalidParam,
			ce.MsgInvalidParam,
			fmt.Errorf("failed to convert url param: %w", err),
		).Bind(
			ctx,
		)
		return
	}

	delErr := h.pu.DeleteProduct(
		ctx.Request.Context(),
		productID,
	)
	if delErr != nil {
		delErr.Bind(ctx)
		return
	}

	utils.SetResponse[any](ctx, http.StatusNoContent, "", nil, nil)
}

func (h *ProductHandler) toProduct(p *models.Product) *dtos.Product {
	if p == nil {
		return nil
	}
	return &dtos.Product{
		ID: p.ID,

		SKU:      p.SKU,
		Name:     p.Name,
		Price:    p.Price,
		Stock:    p.Stock,
		IsActive: p.IsActive,

		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (h *ProductHandler) toProductListItem(pli *models.ProductListItem) *dtos.ProductListItem {
	if pli == nil {
		return nil
	}
	return &dtos.ProductListItem{
		ID: pli.ID,

		SKU:      pli.SKU,
		Name:     pli.Name,
		Price:    pli.Price,
		Stock:    pli.Stock,
		IsActive: pli.IsActive,

		LastMovementAt: pli.LastMovementAt,

		CreatedAt: pli.CreatedAt,
		UpdatedAt: pli.UpdatedAt,
	}
}

func (h *ProductHandler) toStockMovement(sm *models.StockMovement) *dtos.StockMovement {
	if sm == nil {
		return nil
	}
	return &dtos.StockMovement{
		ID:        sm.ID,
		ProductID: sm.ProductID,

		UpdateType: sm.UpdateType,
		Quantity:   sm.Quantity,
		Note:       sm.Note,

		CreatedAt: sm.CreatedAt,
		UpdatedAt: sm.UpdatedAt,
	}
}
