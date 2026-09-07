package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"inventory-service/internal/domain"
	invpg "inventory-service/pkg/postgres"

	"github.com/TruongLe68/go-micro/pkg/pagination"
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

// List returns paginated stock movements matching the given filter.
func (r *StockMovementRepo) List(ctx context.Context, filter domain.StockMovementFilter, p pagination.Params) ([]domain.StockMovement, int64, error) {
	exec := invpg.GetExecutor(ctx, r.db)
	normParams := p.Normalize()

	var conditions []string
	args := make([]interface{}, 0, 4)
	argIdx := 1

	if filter.WarehouseID != "" {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", argIdx))
		args = append(args, filter.WarehouseID)
		argIdx++
	}
	if filter.SKU != "" {
		conditions = append(conditions, fmt.Sprintf("sku = $%d", argIdx))
		args = append(args, filter.SKU)
		argIdx++
	}
	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, string(filter.Type))
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM stock_movements" + whereClause
	var total int64
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("StockMovementRepo.List - count: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, sku, warehouse_id, type, quantity, reference_type, reference_id, note, created_by, created_at
		FROM stock_movements
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	queryArgs := append(args, normParams.Limit, normParams.Skip())

	rows, err := exec.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("StockMovementRepo.List - query: %w", err)
	}
	defer rows.Close()

	movements := make([]domain.StockMovement, 0)
	for rows.Next() {
		var m domain.StockMovement
		var typeStr string
		if err := rows.Scan(
			&m.ID, &m.SKU, &m.WarehouseID, &typeStr, &m.Quantity,
			&m.ReferenceType, &m.ReferenceID, &m.Note, &m.CreatedBy, &m.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("StockMovementRepo.List - scan: %w", err)
		}
		m.Type = domain.MovementType(typeStr)
		movements = append(movements, m)
	}

	return movements, total, rows.Err()
}
