package repositories

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/database"
	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
)

type ProductRepository interface {
	Create(ctx context.Context, data *models.CreateProduct) (p *models.Product, err *ce.Error)
	GetStock(ctx context.Context, productID int) (stock int, err *ce.Error)
	GetAll(ctx context.Context, params *models.GetAllProducts) (plis []models.ProductListItem, total int64, err *ce.Error)
	GetAllBySKUsForUpdate(ctx context.Context, skus []string) (ps []models.Product, err *ce.Error)
	Update(ctx context.Context, productID int, data *models.UpdateProduct) (p *models.Product, err *ce.Error)
	UpdateStock(ctx context.Context, productID, stockDiff int) (p *models.Product, err *ce.Error)
	Delete(ctx context.Context, productID int) (err *ce.Error)
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

func (r *productRepository) GetAll(ctx context.Context, params *models.GetAllProducts) ([]models.ProductListItem, int64, *ce.Error) {
	isTrigramEligible := false
	if params.Search != nil {
		isTrigramEligible = len(*params.Search) >= 3
	}

	query := ""
	args := []any{}
	argPos := 1

	query += `
		SELECT
			p.id, p.sku, p.name, p.price, p.stock,
			p.is_active, p.created_at, p.updated_at,
			MAX(sm.created_at) AS last_movement_at,

			COUNT(*) OVER() AS total
		FROM
			products p
		LEFT JOIN
			stock_movements sm
			ON
				p.id = sm.product_id
	`

	searchArgPos := argPos
	if params.Search != nil {
		if isTrigramEligible {
			query += " WHERE $" + strconv.Itoa(argPos) + " <% p.sku OR $" + strconv.Itoa(argPos) + " <% p.name"
		} else {
			query += " WHERE (p.sku ILIKE '%' || $" + strconv.Itoa(argPos) + " || '%' OR p.name ILIKE '%' || $" + strconv.Itoa(argPos) + " || '%')"
		}
		args = append(args, *params.Search)
		argPos++
	}

	query += " GROUP BY p.id"
	sortClauses := []string{"last_movement_at DESC NULLS LAST"}
	if params.Search != nil && isTrigramEligible {
		sortClauses = append(
			sortClauses,
			"GREATEST(WORD_SIMILARITY($"+strconv.Itoa(searchArgPos)+", p.sku), WORD_SIMILARITY($"+strconv.Itoa(searchArgPos)+", p.name)) DESC",
		)
	}

	query += " ORDER BY " + strings.Join(sortClauses, ", ")
	query += " LIMIT $" + strconv.Itoa(argPos) + " OFFSET $" + strconv.Itoa(argPos+1)
	args = append(args, params.Size, params.Offset())

	rows, err := r.database.QueryAll(
		ctx, query,
		args...,
	)
	if err != nil {
		return nil, 0, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all products: %w", err),
		)
	}
	defer rows.Close()

	var total int64
	plis := make([]models.ProductListItem, 0, params.Size)
	for rows.Next() {
		var pli models.ProductListItem
		var lastMovementAt *time.Time

		err := rows.Scan(
			&pli.ID,
			&pli.SKU,
			&pli.Name,
			&pli.Price,
			&pli.Stock,
			&pli.IsActive,
			&pli.CreatedAt,
			&pli.UpdatedAt,
			&lastMovementAt,
			&total,
		)
		if err != nil {
			return nil, 0, ce.NewError(
				ce.CodeInternal,
				ce.MsgInternalServer,
				fmt.Errorf("failed to get all products: %w", err),
			)
		}
		if lastMovementAt != nil {
			isoStr := lastMovementAt.UTC().Format(time.RFC3339)
			pli.LastMovementAt = &isoStr
		}

		plis = append(plis, pli)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all products: %w", err),
		)
	}
	if len(plis) == 0 {
		countQuery := "SELECT COUNT(*) FROM products"
		countArgs := []any{}
		countArgPos := 1
		if params.Search != nil {
			if isTrigramEligible {
				countQuery += " WHERE $" + strconv.Itoa(countArgPos) + " <% sku OR $" + strconv.Itoa(countArgPos) + " <% name"
			} else {
				countQuery += " WHERE (sku ILIKE '%' || $" + strconv.Itoa(countArgPos) + " || '%' OR name ILIKE '%' || $" + strconv.Itoa(countArgPos) + " || '%')"
			}
			countArgs = append(countArgs, *params.Search)
			countArgPos++
		}

		err := r.database.QueryOne(
			ctx, countQuery,
			countArgs...,
		).Scan(
			&total,
		)
		if err != nil {
			return nil, 0, ce.NewError(
				ce.CodeInternal,
				ce.MsgInternalServer,
				fmt.Errorf("failed to get total products count: %w", err),
			)
		}
	}

	return plis, total, nil
}

func (r *productRepository) GetAllBySKUsForUpdate(ctx context.Context, skus []string) ([]models.Product, *ce.Error) {
	query := `
		SELECT
			id, sku, name, price, stock,
			is_active, created_at, updated_at
		FROM
			products
		WHERE
			sku = ANY($1)
	`
	if r.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	rows, err := r.database.QueryAll(
		ctx, query,
		skus,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all products by skus for update: %w", err),
		)
	}
	defer rows.Close()

	ps := make([]models.Product, 0, len(skus))
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
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
			return nil, ce.NewError(
				ce.CodeInternal,
				ce.MsgInternalServer,
				fmt.Errorf("failed to get all products by skus for update: %w", err),
			)
		}

		ps = append(ps, p)
	}
	if err := rows.Err(); err != nil {
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all products by skus for update: %w", err),
		)
	}

	return ps, nil
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

func (r *productRepository) Delete(ctx context.Context, productID int) *ce.Error {
	query := "DELETE FROM products WHERE id = $1"

	err := r.database.Execute(
		ctx, query,
		productID,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete product: %w", err)

		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(
				ce.CodeProductNotFound,
				ce.MsgProductNotFound,
				wrappedErr,
			)
		}
		return ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return nil
}
