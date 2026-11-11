package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/stockflow/backend/internal/infra/database"
	"github.com/ritchieridanko/stockflow/backend/internal/models"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
)

type StockMovementRepository interface {
	Create(ctx context.Context, data *models.CreateStockMovement) (sm *models.StockMovement, err *ce.Error)
	GetAllByProductID(ctx context.Context, productID int) (sms []models.StockMovement, err *ce.Error)
	HasAny(ctx context.Context, productID int) (exists bool, err *ce.Error)
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
			product_id, type, quantity, note
		)
		VALUES (
			$1, $2, $3, $4
		)
		RETURNING
			id, product_id, type, quantity,
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

func (r *stockMovementRepository) GetAllByProductID(ctx context.Context, productID int) ([]models.StockMovement, *ce.Error) {
	query := `
		SELECT
			id, product_id, type, quantity,
			note, created_at, updated_at
		FROM
			stock_movements
		WHERE
			product_id = $1
		ORDER BY
			created_at DESC,
			id ASC
	`

	rows, err := r.database.QueryAll(
		ctx, query,
		productID,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all stock movements by product id: %w", err),
		)
	}
	defer rows.Close()

	sms := make([]models.StockMovement, 0)
	for rows.Next() {
		var sm models.StockMovement
		err := rows.Scan(
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
				fmt.Errorf("failed to get all stock movements by product id: %w", err),
			)
		}

		sms = append(sms, sm)
	}
	if err := rows.Err(); err != nil {
		return nil, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to get all stock movements by product id: %w", err),
		)
	}

	return sms, nil
}

func (r *stockMovementRepository) HasAny(ctx context.Context, productID int) (bool, *ce.Error) {
	query := "SELECT 1 FROM stock_movements WHERE product_id = $1"
	if r.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var exists int
	err := r.database.QueryOne(
		ctx, query,
		productID,
	).Scan(
		&exists,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return false, nil
		}
		return false, ce.NewError(
			ce.CodeInternal,
			ce.MsgInternalServer,
			fmt.Errorf("failed to check if product has any stock movements: %w", err),
		)
	}

	return true, nil
}
