package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"inventory-service/internal/domain"
	invpg "inventory-service/pkg/postgres"
)

// StockLevelRepo implements stock level persistence with optimistic locking.
type StockLevelRepo struct {
	db *sql.DB
}

// NewStockLevelRepo creates a new StockLevelRepo.
func NewStockLevelRepo(db *sql.DB) *StockLevelRepo {
	return &StockLevelRepo{db: db}
}

// FindBySKU returns all stock levels across warehouses for a given SKU.
func (r *StockLevelRepo) FindBySKU(ctx context.Context, sku string) ([]domain.StockLevel, error) {
	exec := invpg.GetExecutor(ctx, r.db)

	rows, err := exec.QueryContext(ctx, `
		SELECT id, sku, warehouse_id, on_hand, reserved,
		       reorder_threshold, reorder_quantity, version,
		       created_at, updated_at
		FROM stock_levels
		WHERE sku = $1
		ORDER BY warehouse_id
	`, sku)
	if err != nil {
		return nil, fmt.Errorf("StockLevelRepo.FindBySKU - query: %w", err)
	}
	defer rows.Close()

	var levels []domain.StockLevel
	for rows.Next() {
		var sl domain.StockLevel
		if err := rows.Scan(
			&sl.ID, &sl.SKU, &sl.WarehouseID,
			&sl.OnHand, &sl.Reserved,
			&sl.ReorderThreshold, &sl.ReorderQuantity,
			&sl.Version, &sl.CreatedAt, &sl.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("StockLevelRepo.FindBySKU - scan: %w", err)
		}
		levels = append(levels, sl)
	}

	return levels, rows.Err()
}

// FindBySKUAndWarehouse returns a single stock level for a given SKU + warehouse pair.
func (r *StockLevelRepo) FindBySKUAndWarehouse(ctx context.Context, sku, warehouseID string) (*domain.StockLevel, error) {
	exec := invpg.GetExecutor(ctx, r.db)

	var sl domain.StockLevel
	err := exec.QueryRowContext(ctx, `
		SELECT id, sku, warehouse_id, on_hand, reserved,
		       reorder_threshold, reorder_quantity, version,
		       created_at, updated_at
		FROM stock_levels
		WHERE sku = $1 AND warehouse_id = $2
	`, sku, warehouseID).Scan(
		&sl.ID, &sl.SKU, &sl.WarehouseID,
		&sl.OnHand, &sl.Reserved,
		&sl.ReorderThreshold, &sl.ReorderQuantity,
		&sl.Version, &sl.CreatedAt, &sl.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("StockLevelRepo.FindBySKUAndWarehouse: %w", err)
	}

	return &sl, nil
}

// Create inserts a new stock level.
func (r *StockLevelRepo) Create(ctx context.Context, sl *domain.StockLevel) error {
	exec := invpg.GetExecutor(ctx, r.db)

	err := exec.QueryRowContext(ctx, `
		INSERT INTO stock_levels (sku, warehouse_id, on_hand, reserved, reorder_threshold, reorder_quantity, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`,
		sl.SKU, sl.WarehouseID, sl.OnHand, sl.Reserved,
		sl.ReorderThreshold, sl.ReorderQuantity, sl.Version,
		sl.CreatedAt, sl.UpdatedAt,
	).Scan(&sl.ID)
	if err != nil {
		return fmt.Errorf("StockLevelRepo.Create: %w", err)
	}

	return nil
}

// FindFirstAvailable finds the first warehouse with enough available stock for the given SKU and quantity.
func (r *StockLevelRepo) FindFirstAvailable(ctx context.Context, sku string, qty int) (*domain.StockLevel, error) {
	exec := invpg.GetExecutor(ctx, r.db)

	var sl domain.StockLevel
	err := exec.QueryRowContext(ctx, `
		SELECT id, sku, warehouse_id, on_hand, reserved,
		       reorder_threshold, reorder_quantity, version,
		       created_at, updated_at
		FROM stock_levels
		WHERE sku = $1 AND (on_hand - reserved) >= $2
		ORDER BY (on_hand - reserved) DESC
		LIMIT 1
	`, sku, qty).Scan(
		&sl.ID, &sl.SKU, &sl.WarehouseID,
		&sl.OnHand, &sl.Reserved,
		&sl.ReorderThreshold, &sl.ReorderQuantity,
		&sl.Version, &sl.CreatedAt, &sl.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrInsufficientStock
		}
		return nil, fmt.Errorf("StockLevelRepo.FindFirstAvailable: %w", err)
	}

	return &sl, nil
}

// Update persists a stock level change with optimistic locking on version.
func (r *StockLevelRepo) Update(ctx context.Context, sl *domain.StockLevel) error {
	exec := invpg.GetExecutor(ctx, r.db)

	result, err := exec.ExecContext(ctx, `
		UPDATE stock_levels
		SET on_hand = $1, reserved = $2, reorder_threshold = $3,
		    reorder_quantity = $4, version = $5, updated_at = $6
		WHERE id = $7 AND version = $8
	`,
		sl.OnHand, sl.Reserved, sl.ReorderThreshold,
		sl.ReorderQuantity, sl.Version, sl.UpdatedAt,
		sl.ID, sl.Version-1, // previous version for optimistic lock
	)
	if err != nil {
		return fmt.Errorf("StockLevelRepo.Update - exec: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("StockLevelRepo.Update - RowsAffected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("%w: stock_level id=%s", domain.ErrConcurrentModification, sl.ID)
	}

	return nil
}

// BulkCheckAvailability returns the total available quantity (on_hand - reserved)
// for each requested SKU, summed across all warehouses.
func (r *StockLevelRepo) BulkCheckAvailability(ctx context.Context, skus []string) (map[string]int, error) {
	exec := invpg.GetExecutor(ctx, r.db)

	if len(skus) == 0 {
		return map[string]int{}, nil
	}

	// Build query with ANY($1)
	rows, err := exec.QueryContext(ctx, `
		SELECT sku, COALESCE(SUM(GREATEST(on_hand - reserved, 0)), 0) AS available
		FROM stock_levels
		WHERE sku = ANY($1)
		GROUP BY sku
	`, pqStringArray(skus))
	if err != nil {
		return nil, fmt.Errorf("StockLevelRepo.BulkCheckAvailability - query: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int, len(skus))
	for rows.Next() {
		var sku string
		var avail int
		if err := rows.Scan(&sku, &avail); err != nil {
			return nil, fmt.Errorf("StockLevelRepo.BulkCheckAvailability - scan: %w", err)
		}
		result[sku] = avail
	}

	// SKUs with no stock_level rows will have 0
	for _, sku := range skus {
		if _, ok := result[sku]; !ok {
			result[sku] = 0
		}
	}

	return result, rows.Err()
}

// pqStringArray converts a []string to a driver.Valuer for PostgreSQL's text[].
type pqStringArray []string

func (a pqStringArray) Value() (interface{}, error) {
	return "{" + joinStrings(a) + "}", nil
}

func joinStrings(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	result := `"` + arr[0] + `"`
	for _, s := range arr[1:] {
		result += `,"` + s + `"`
	}
	return result
}
