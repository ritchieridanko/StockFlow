package usecases

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/ritchieridanko/stockflow/backend/internal/constants"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/database"
	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/repositories"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/validator"
	"go.opentelemetry.io/otel/trace"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, req *models.CreateOrderReq) (o *models.Order, err *ce.Error)
	GetOrderByID(ctx context.Context, orderID int) (o *models.Order, err *ce.Error)
}

type orderUsecase struct {
	or         repositories.OrderRepository
	oir        repositories.OrderItemRepository
	pr         repositories.ProductRepository
	smr        repositories.StockMovementRepository
	transactor *database.Transactor
	validator  *validator.Validator
	tracer     trace.Tracer
}

func NewOrderUsecase(
	or repositories.OrderRepository,
	oir repositories.OrderItemRepository,
	pr repositories.ProductRepository,
	smr repositories.StockMovementRepository,
	tx *database.Transactor,
	v *validator.Validator,
	t trace.Tracer,
) OrderUsecase {
	return &orderUsecase{
		or:         or,
		oir:        oir,
		pr:         pr,
		smr:        smr,
		transactor: tx,
		validator:  v,
		tracer:     t,
	}
}

func (u *orderUsecase) CreateOrder(ctx context.Context, req *models.CreateOrderReq) (*models.Order, *ce.Error) {
	ctx, span := u.tracer.Start(ctx, "order.create")
	defer span.End()

	// Data Normalization And Validation
	customer := strings.TrimSpace(req.Customer)
	if ok, why := u.validator.Customer(customer); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}
	if len(req.Items) == 0 {
		return nil, ce.NewError(ce.CodeInvalidPayload, ce.MsgEmptyOrderItems, nil)
	}

	skuSet := make(map[string]bool)
	skus := make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		sku := strings.TrimSpace(item.SKU)

		if ok, why := u.validator.SKU(sku); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
		}
		if ok, why := u.validator.Quantity(item.Quantity); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
		}
		if skuSet[sku] {
			return nil, ce.NewError(ce.CodeInvalidPayload, "Duplicate SKU in order: "+sku, nil)
		}

		skuSet[sku] = true
		skus = append(skus, sku)
	}
	sort.Strings(skus)

	var o *models.Order
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// All Products Fetching By SKUs
		ps, err := u.pr.GetAllBySKUsForUpdate(ctx, skus)
		if err != nil {
			return err
		}

		productMap := make(map[string]models.Product, len(ps))
		for _, p := range ps {
			productMap[p.SKU] = p
		}

		// Stock Availability Validation
		total := 0
		orderItems := make([]models.CreateOrderItem, 0, len(req.Items))

		for _, item := range req.Items {
			p, exists := productMap[item.SKU]
			if !exists {
				return ce.NewError(ce.CodeInvalidPayload, "Product with SKU not found: "+item.SKU, nil)
			}
			if p.Stock < item.Quantity {
				return ce.NewError(ce.CodeInsufficientStock, "Insufficient stock for SKU: "+item.SKU, nil)
			}

			total += p.Price * item.Quantity
			orderItems = append(orderItems,
				models.CreateOrderItem{
					SKU:       item.SKU,
					Quantity:  item.Quantity,
					UnitPrice: p.Price,
				},
			)
		}

		// Order Creation
		o, err = u.or.Create(
			ctx,
			&models.CreateOrder{
				Customer: customer,
				Total:    total,
			},
		)
		if err != nil {
			return err
		}

		// Order Items Creation
		ois, err := u.oir.CreateMany(
			ctx,
			&models.CreateOrderItems{
				OrderID: o.ID,
				Items:   orderItems,
			},
		)
		if err != nil {
			return err
		}
		o.Items = ois

		for _, item := range req.Items {
			p := productMap[item.SKU]
			note := "Sales Order #" + strconv.Itoa(o.ID)

			// Stock Movement Creation
			_, err := u.smr.Create(
				ctx,
				&models.CreateStockMovement{
					ProductID:  p.ID,
					UpdateType: constants.StockUpdateTypeOut,
					Quantity:   item.Quantity,
					Note:       &note,
				},
			)
			if err != nil {
				return err
			}

			// Product Stock Update
			_, err = u.pr.UpdateStock(ctx, p.ID, -item.Quantity)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (u *orderUsecase) GetOrderByID(ctx context.Context, orderID int) (*models.Order, *ce.Error) {
	ctx, span := u.tracer.Start(ctx, "order.get_by_id")
	defer span.End()

	o, err := u.or.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	ois, err := u.oir.GetAll(ctx, orderID)
	if err != nil {
		return nil, err
	}

	o.Items = ois
	return o, nil
}
