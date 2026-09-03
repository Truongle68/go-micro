package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"inventory-service/internal/domain"
	invpg "inventory-service/pkg/postgres"
)

// StockMovementRepo implements stock movement audit trail persistence.
type StockMovementRepo struct {
	db *sql.DB
}

// NewStockMovementRepo creates a new StockMovementRepo.
func NewStockMovementRepo(db *sql.DB) *StockMovementRepo {
	return &StockMovementRepo{db: db}
}

// Create inserts a new stock movement record.
func (r *StockMovementRepo) Create(ctx context.Context, m *domain.StockMovement) error {
	exec := invpg.GetExecutor(ctx, r.db)

	err := exec.QueryRowContext(ctx, `
		INSERT INTO stock_movements
			(sku, warehouse_id, type, quantity, reference_type, reference_id, note, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`,
		m.SKU, m.WarehouseID, string(m.Type), m.Quantity,
		m.ReferenceType, m.ReferenceID, m.Note, m.CreatedBy, m.CreatedAt,
	).Scan(&m.ID)
	if err != nil {
		return fmt.Errorf("StockMovementRepo.Create: %w", err)
	}

	return nil
}
