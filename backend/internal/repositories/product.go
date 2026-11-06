package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/database"
	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
)

type ProductRepository interface {
	Create(ctx context.Context, data *models.CreateProduct) (p *models.Product, err *ce.Error)
}

type productRepository struct {
	database *database.Database
}

func NewProductRepository(db *database.Database) ProductRepository {
	return &productRepository{database: db}
}

func (r *productRepository) Create(ctx context.Context, data *models.CreateProduct) (*models.Product, *ce.Error) {
	query := `
		INSERT INTO products (
			sku, name, price
		)
		VALUES (
			$1, $2, $3
		)
		RETURNING
			id, sku, name, price, stock,
			is_active, created_at, updated_at
	`

	var p models.Product
	err := r.database.QueryOne(
		ctx, query,
		data.SKU,
		data.Name,
		data.Price,
	).Scan(
		&p.ID,
		&p.SKU,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.IsActive,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create product: %w", err)

		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ce.NewError(
				ce.CodeSKUAlreadyExists,
				ce.MsgSKUAlreadyExists,
				wrappedErr,
			)
		}
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return &p, nil
}
