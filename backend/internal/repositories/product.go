package repositories

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/database"
	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
)

type ProductRepository interface {
	Create(ctx context.Context, data *models.CreateProduct) (p *models.Product, err *ce.Error)
	GetStock(ctx context.Context, productID int) (stock int, err *ce.Error)
	Update(ctx context.Context, productID int, data *models.UpdateProduct) (p *models.Product, err *ce.Error)
	UpdateStock(ctx context.Context, productID, stockDiff int) (p *models.Product, err *ce.Error)
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

func (r *productRepository) GetStock(ctx context.Context, productID int) (int, *ce.Error) {
	query := "SELECT stock FROM products WHERE id = $1"
	if r.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var stock int
	err := r.database.QueryOne(
		ctx, query,
		productID,
	).Scan(
		&stock,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get product stock: %w", err)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, ce.NewError(
				ce.CodeProductNotFound,
				ce.MsgProductNotFound,
				wrappedErr,
			)
		}
		return 0, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return stock, nil
}

func (r *productRepository) Update(ctx context.Context, productID int, data *models.UpdateProduct) (*models.Product, *ce.Error) {
	setClauses := []string{}
	args := []any{}
	argPos := 1

	if data.SKU != nil {
		setClauses = append(setClauses, "sku = $"+strconv.Itoa(argPos))
		args = append(args, *data.SKU)
		argPos++
	}
	if data.Name != nil {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(argPos))
		args = append(args, *data.Name)
		argPos++
	}
	if data.Price != nil {
		setClauses = append(setClauses, "price = $"+strconv.Itoa(argPos))
		args = append(args, *data.Price)
		argPos++
	}
	if len(setClauses) == 0 {
		return nil, ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, nil)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, productID)
	query := fmt.Sprintf(
		`
			UPDATE
				products
			SET
				%s
			WHERE
				id = $%d
			RETURNING
				id, sku, name, price, stock,
				is_active, created_at, updated_at
		`,
		strings.Join(setClauses, ", "),
		argPos,
	)

	var p models.Product
	err := r.database.QueryOne(
		ctx, query,
		args...,
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
		wrappedErr := fmt.Errorf("failed to update product: %w", err)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeProductNotFound,
				ce.MsgProductNotFound,
				wrappedErr,
			)
		}
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

func (r *productRepository) UpdateStock(ctx context.Context, productID, stockDiff int) (*models.Product, *ce.Error) {
	query := `
		UPDATE
			products
		SET
			stock = stock + $2,
			updated_at = NOW()
		WHERE
			id = $1
		RETURNING
			id, sku, name, price, stock,
			is_active, created_at, updated_at
	`

	var p models.Product
	err := r.database.QueryOne(
		ctx, query,
		productID,
		stockDiff,
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
		wrappedErr := fmt.Errorf("failed to update product stock: %w", err)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeProductNotFound,
				ce.MsgProductNotFound,
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
