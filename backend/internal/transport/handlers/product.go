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
	)
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
