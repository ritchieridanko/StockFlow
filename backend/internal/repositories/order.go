package repositories

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ritchieridanko/stockflow/backend/internal/infra/database"
	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
)

type OrderRepository interface {
	Create(ctx context.Context, data *models.CreateOrder) (o *models.Order, err *ce.Error)
	GetByID(ctx context.Context, orderID int) (o *models.Order, err *ce.Error)
	GetAll(ctx context.Context, params *models.GetAllOrders) (os []models.Order, total int64, err *ce.Error)
	Cancel(ctx context.Context, orderID int) (o *models.Order, err *ce.Error)
}

type orderRepository struct {
	database *database.Database
}

func NewOrderRepository(db *database.Database) OrderRepository {
	return &orderRepository{database: db}
}

func (r *orderRepository) Create(ctx context.Context, data *models.CreateOrder) (*models.Order, *ce.Error) {
	query := `
		INSERT INTO sales_orders (
			customer, total
		)
		VALUES (
			$1, $2
		)
		RETURNING
			id, customer, status, total,
			created_at, updated_at, cancelled_at
	`

	var o models.Order
	err := r.database.QueryOne(
		ctx, query,
		data.Customer,
		data.Total,
	).Scan(
		&o.ID,
		&o.Customer,
		&o.Status,
		&o.Total,
		&o.CreatedAt,
		&o.UpdatedAt,
		&o.CancelledAt,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create order: %w", err),
		)
	}

	return &o, nil
}

func (r *orderRepository) GetByID(ctx context.Context, orderID int) (*models.Order, *ce.Error) {
	query := `
		SELECT
			id, customer, status, total,
			created_at, updated_at, cancelled_at
		FROM
			sales_orders
		WHERE
			id = $1
	`
	if r.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var o models.Order
	err := r.database.QueryOne(
		ctx, query,
		orderID,
	).Scan(
		&o.ID,
		&o.Customer,
		&o.Status,
		&o.Total,
		&o.CreatedAt,
		&o.UpdatedAt,
		&o.CancelledAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get order by id: %w", err)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeOrderNotFound,
				ce.MsgOrderNotFound,
				wrappedErr,
			)
		}
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return &o, nil
}

func (r *orderRepository) GetAll(ctx context.Context, params *models.GetAllOrders) ([]models.Order, int64, *ce.Error) {
	isTrigramEligible := false
	if params.Search != nil {
		isTrigramEligible = len(*params.Search) >= 3
	}

	query := ""
	args := []any{}
	argPos := 1

	query += `
		SELECT
			id, customer, status, total,
			created_at, updated_at, cancelled_at,

			COUNT(*) OVER() AS total
		FROM
			sales_orders
	`

	searchArgPos := argPos
	if params.Search != nil {
		if isTrigramEligible {
			query += " WHERE $" + strconv.Itoa(argPos) + " <% customer"
		} else {
			query += " WHERE customer ILIKE '%' || $" + strconv.Itoa(argPos) + " || '%'"
		}
		args = append(args, *params.Search)
		argPos++
	}

	sortClauses := []string{"created_at DESC"}
	if params.Search != nil && isTrigramEligible {
		sortClauses = append(
			sortClauses,
			"WORD_SIMILARITY($"+strconv.Itoa(searchArgPos)+", customer) DESC",
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
			fmt.Errorf("failed to get all orders: %w", err),
		)
	}
	defer rows.Close()

	var total int64
	os := make([]models.Order, 0, params.Size)
	for rows.Next() {
		var o models.Order
		err := rows.Scan(
			&o.ID,
			&o.Customer,
			&o.Status,
			&o.Total,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.CancelledAt,
			&total,
		)
		if err != nil {
			return nil, 0, ce.NewError(
				ce.CodeInternal,
				ce.MsgInternalServer,
				fmt.Errorf("failed to get all orders: %w", err),
			)
		}

		os = append(os, o)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all orders: %w", err),
		)
	}
	if len(os) == 0 {
		countQuery := "SELECT COUNT(*) FROM sales_orders"
		countArgs := []any{}
		countArgPos := 1
		if params.Search != nil {
			if isTrigramEligible {
				countQuery += " WHERE $" + strconv.Itoa(countArgPos) + " <% customer"
			} else {
				countQuery += " WHERE customer ILIKE '%' || $" + strconv.Itoa(countArgPos) + " || '%'"
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
				fmt.Errorf("failed to get total orders count: %w", err),
			)
		}
	}

	return os, total, nil
}

func (r *orderRepository) Cancel(ctx context.Context, orderID int) (*models.Order, *ce.Error) {
	query := `
		UPDATE
			sales_orders
		SET
			status = 'CANCELLED',
			cancelled_at = NOW(),
			updated_at = NOW()
		WHERE
			id = $1
		RETURNING
			id, customer, status, total,
			created_at, updated_at, cancelled_at
	`

	var o models.Order
	err := r.database.QueryOne(
		ctx, query,
		orderID,
	).Scan(
		&o.ID,
		&o.Customer,
		&o.Status,
		&o.Total,
		&o.CreatedAt,
		&o.UpdatedAt,
		&o.CancelledAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to cancel order: %w", err)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeOrderNotFound,
				ce.MsgOrderNotFound,
				wrappedErr,
			)
		}
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return &o, nil
}
