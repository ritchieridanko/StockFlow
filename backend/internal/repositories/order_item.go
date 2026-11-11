package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/ritchieridanko/stockflow/backend/internal/infra/database"
	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
)

type OrderItemRepository interface {
	CreateMany(ctx context.Context, data *models.CreateOrderItems) (ois []models.OrderItem, err *ce.Error)
	GetAll(ctx context.Context, orderID int) (ois []models.OrderItem, err *ce.Error)
}

type orderItemRepository struct {
	database *database.Database
}

func NewOrderItemRepository(db *database.Database) OrderItemRepository {
	return &orderItemRepository{database: db}
}

func (r *orderItemRepository) CreateMany(ctx context.Context, data *models.CreateOrderItems) ([]models.OrderItem, *ce.Error) {
	totalItems := len(data.Items)
	if totalItems == 0 {
		return nil, ce.NewError(ce.CodeInvalidPayload, ce.MsgEmptyOrderItems, nil)
	}

	valClauses := []string{}
	args := []any{}
	numCols := 4

	for i, item := range data.Items {
		base := i * numCols

		valClauses = append(valClauses,
			fmt.Sprintf(
				"($%d, $%d, $%d, $%d)",
				base+1, base+2, base+3, base+4,
			),
		)
		args = append(args,
			data.OrderID,
			item.SKU,
			item.Quantity,
			item.UnitPrice,
		)
	}

	query := fmt.Sprintf(
		`
			INSERT INTO sales_order_items (
				order_id, sku, quantity, unit_price
			)
			VALUES
				%s
			RETURNING
				id, order_id, sku, quantity,
				unit_price, created_at, updated_at
		`,
		strings.Join(valClauses, ", "),
	)

	rows, err := r.database.QueryAll(
		ctx, query,
		args...,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create order items: %w", err),
		)
	}
	defer rows.Close()

	ois := make([]models.OrderItem, 0, totalItems)
	for rows.Next() {
		var oi models.OrderItem
		err := rows.Scan(
			&oi.ID,
			&oi.OrderID,
			&oi.SKU,
			&oi.Quantity,
			&oi.UnitPrice,
			&oi.CreatedAt,
			&oi.UpdatedAt,
		)
		if err != nil {
			return nil, ce.NewError(
				ce.CodeInternal,
				ce.MsgInternalServer,
				fmt.Errorf("failed to create order items: %w", err),
			)
		}

		ois = append(ois, oi)
	}
	if err := rows.Err(); err != nil {
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create order items: %w", err),
		)
	}

	return ois, nil
}

func (r *orderItemRepository) GetAll(ctx context.Context, orderID int) ([]models.OrderItem, *ce.Error) {
	query := `
		SELECT
			id, order_id, sku, quantity,
			unit_price, created_at, updated_at
		FROM
			sales_order_items
		WHERE
			order_id = $1
	`

	rows, err := r.database.QueryAll(
		ctx, query,
		orderID,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all order items: %w", err),
		)
	}
	defer rows.Close()

	ois := make([]models.OrderItem, 0)
	for rows.Next() {
		var oi models.OrderItem
		err := rows.Scan(
			&oi.ID,
			&oi.OrderID,
			&oi.SKU,
			&oi.Quantity,
			&oi.UnitPrice,
			&oi.CreatedAt,
			&oi.UpdatedAt,
		)
		if err != nil {
			return nil, ce.NewError(
				ce.CodeInternal,
				ce.MsgInternalServer,
				fmt.Errorf("failed to get all order items: %w", err),
			)
		}

		ois = append(ois, oi)
	}
	if err := rows.Err(); err != nil {
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all order items: %w", err),
		)
	}

	return ois, nil
}
