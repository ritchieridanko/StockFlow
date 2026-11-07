package usecases

import (
	"context"
	"strings"

	"github.com/ritchieridanko/stockflow/backend/internal/constants"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/database"
	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/repositories"
	"github.com/ritchieridanko/stockflow/backend/internal/utils"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/validator"
	"go.opentelemetry.io/otel/trace"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *models.CreateProductReq) (p *models.Product, err *ce.Error)
	GetAllProductStockMovements(ctx context.Context, productID int) (sms []models.StockMovement, err *ce.Error)
	UpdateProduct(ctx context.Context, req *models.UpdateProductReq) (p *models.Product, err *ce.Error)
	UpdateStock(ctx context.Context, req *models.UpdateStockReq) (p *models.Product, err *ce.Error)
}

type productUsecase struct {
	pr         repositories.ProductRepository
	smr        repositories.StockMovementRepository
	transactor *database.Transactor
	validator  *validator.Validator
	tracer     trace.Tracer
}

func NewProductUsecase(
	pr repositories.ProductRepository,
	smr repositories.StockMovementRepository,
	tx *database.Transactor,
	v *validator.Validator,
	t trace.Tracer,
) ProductUsecase {
	return &productUsecase{
		pr:         pr,
		smr:        smr,
		transactor: tx,
		validator:  v,
		tracer:     t,
	}
}

func (u *productUsecase) CreateProduct(ctx context.Context, req *models.CreateProductReq) (*models.Product, *ce.Error) {
	ctx, span := u.tracer.Start(ctx, "product.create")
	defer span.End()

	// Data Normalization
	sku := strings.TrimSpace(req.SKU)
	name := strings.TrimSpace(req.Name)

	// Data Validation
	if ok, why := u.validator.SKU(sku); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}
	if ok, why := u.validator.Name(name); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}
	if ok, why := u.validator.Price(req.Price); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}

	// Product Creation
	return u.pr.Create(
		ctx,
		&models.CreateProduct{
			SKU:   sku,
			Name:  name,
			Price: req.Price,
		},
	)
}

func (u *productUsecase) GetAllProductStockMovements(ctx context.Context, productID int) ([]models.StockMovement, *ce.Error) {
	return u.smr.GetAllByProductID(ctx, productID)
}

func (u *productUsecase) UpdateProduct(ctx context.Context, req *models.UpdateProductReq) (*models.Product, *ce.Error) {
	ctx, span := u.tracer.Start(ctx, "product.update")
	defer span.End()

	// Data Normalization
	sku := utils.TrimSpacePtr(req.SKU)
	name := utils.TrimSpacePtr(req.Name)

	// Data Validation
	if sku != nil {
		if ok, why := u.validator.SKU(*sku); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
		}
	}
	if name != nil {
		if ok, why := u.validator.Name(*name); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
		}
	}
	if req.Price != nil {
		if ok, why := u.validator.Price(*req.Price); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
		}
	}

	// Product Update
	return u.pr.Update(
		ctx,
		req.ProductID,
		&models.UpdateProduct{
			SKU:   sku,
			Name:  name,
			Price: req.Price,
		},
	)
}

func (u *productUsecase) UpdateStock(ctx context.Context, req *models.UpdateStockReq) (*models.Product, *ce.Error) {
	ctx, span := u.tracer.Start(ctx, "product.update_stock")
	defer span.End()

	// Data Normalization
	updateType := strings.ToLower(strings.TrimSpace(req.UpdateType))
	note := utils.TrimSpacePtr(req.Note)

	// Data Validation
	if ok, why := u.validator.StockUpdateType(updateType); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}
	if ok, why := u.validator.Quantity(req.Quantity); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}

	var p *models.Product
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		var stockDiff int
		if updateType == constants.StockUpdateTypeOut {
			// Product Stock Fetching
			stock, err := u.pr.GetStock(ctx, req.ProductID)
			if err != nil {
				return err
			}
			if stock < req.Quantity {
				return ce.NewError(ce.CodeInsufficientStock, ce.MsgInsufficientStock, nil)
			}

			stockDiff = -req.Quantity
		}
		if updateType == constants.StockUpdateTypeIn {
			stockDiff = req.Quantity
		}

		// Stock Movement Creation
		_, err := u.smr.Create(
			ctx,
			&models.CreateStockMovement{
				ProductID:  req.ProductID,
				UpdateType: updateType,
				Quantity:   req.Quantity,
				Note:       note,
			},
		)
		if err != nil {
			return err
		}

		// Product Stock Update
		p, err = u.pr.UpdateStock(ctx, req.ProductID, stockDiff)
		return err
	})

	return p, err
}
