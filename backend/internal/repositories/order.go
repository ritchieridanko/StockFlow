package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/stockflow/backend/internal/infra/database"
	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
)

type OrderRepository interface {
	Create(ctx context.Context, data *models.CreateOrder) (o *models.Order, err *ce.Error)
	GetByID(ctx context.Context, orderID int) (o *models.Order, err *ce.Error)
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
