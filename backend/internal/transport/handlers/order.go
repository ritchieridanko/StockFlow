package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/transport/dtos"
	"github.com/ritchieridanko/stockflow/backend/internal/usecases"
	"github.com/ritchieridanko/stockflow/backend/internal/utils"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
)

type OrderHandler struct {
	ou usecases.OrderUsecase
}

func NewOrderHandler(ou usecases.OrderUsecase) *OrderHandler {
	return &OrderHandler{ou: ou}
}

func (h *OrderHandler) Create(ctx *gin.Context) {
	var payload dtos.CreateOrderRequest
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

	items := make([]models.OrderItemReq, 0, len(payload.Items))
	for _, item := range payload.Items {
		items = append(items, *h.toOrderItemReq(&item))
	}

	o, err := h.ou.CreateOrder(
		ctx.Request.Context(),
		&models.CreateOrderReq{
			Customer: payload.Customer,
			Items:    items,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetResponse(
		ctx,
		http.StatusCreated,
		"Order created successfully",
		dtos.CreateOrderResponse{
			Order: h.toOrder(o),
		},
	)
}

func (h *OrderHandler) toOrder(o *models.Order) *dtos.Order {
	if o == nil {
		return nil
	}

	items := make([]dtos.OrderItem, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, *h.toOrderItem(&item))
	}

	return &dtos.Order{
		ID: o.ID,

		Customer: o.Customer,
		Status:   o.Status,
		Total:    o.Total,

		Items: items,

		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
		CancelledAt: o.CancelledAt,
	}
}

func (h *OrderHandler) toOrderItem(oi *models.OrderItem) *dtos.OrderItem {
	if oi == nil {
		return nil
	}
	return &dtos.OrderItem{
		ID:      oi.ID,
		OrderID: oi.OrderID,

		SKU:       oi.SKU,
		Quantity:  oi.Quantity,
		UnitPrice: oi.UnitPrice,

		CreatedAt: oi.CreatedAt,
		UpdatedAt: oi.UpdatedAt,
	}
}

func (h *OrderHandler) toOrderItemReq(req *dtos.OrderItemRequest) *models.OrderItemReq {
	if req == nil {
		return nil
	}
	return &models.OrderItemReq{
		SKU:      req.SKU,
		Quantity: req.Quantity,
	}
}
