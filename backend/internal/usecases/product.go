package usecases

import (
	"context"
	"strings"

	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/repositories"
	"github.com/ritchieridanko/stockflow/backend/internal/utils"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/validator"
	"go.opentelemetry.io/otel/trace"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *models.CreateProductReq) (p *models.Product, err *ce.Error)
	UpdateProduct(ctx context.Context, req *models.UpdateProductReq) (p *models.Product, err *ce.Error)
}

type productUsecase struct {
	pr        repositories.ProductRepository
	validator *validator.Validator
	tracer    trace.Tracer
}

func NewProductUsecase(
	pr repositories.ProductRepository,
	v *validator.Validator,
	t trace.Tracer,
) ProductUsecase {
	return &productUsecase{
		pr:        pr,
		validator: v,
		tracer:    t,
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
