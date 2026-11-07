package repositories

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/stockflow/backend/internal/infra/database"
	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
)

type StockMovementRepository interface {
	Create(ctx context.Context, data *models.CreateStockMovement) (sm *models.StockMovement, err *ce.Error)
}

type stockMovementRepository struct {
	database *database.Database
}

func NewStockMovementRepository(db *database.Database) StockMovementRepository {
	return &stockMovementRepository{database: db}
}

func (r *stockMovementRepository) Create(ctx context.Context, data *models.CreateStockMovement) (*models.StockMovement, *ce.Error) {
	query := `
		INSERT INTO stock_movements (
			product_id, type, qty, note
		)
		VALUES (
			$1, $2, $3, $4
		)
		RETURNING
			id, product_id, type, qty,
			note, created_at, updated_at
	`

	var sm models.StockMovement
	err := r.database.QueryOne(
		ctx, query,
		data.ProductID,
		data.UpdateType,
		data.Quantity,
		data.Note,
	).Scan(
		&sm.ID,
		&sm.ProductID,
		&sm.UpdateType,
		&sm.Quantity,
		&sm.Note,
		&sm.CreatedAt,
		&sm.UpdatedAt,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create stock movement: %w", err),
		)
	}

	return &sm, nil
}
