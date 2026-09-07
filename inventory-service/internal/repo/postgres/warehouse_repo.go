package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"inventory-service/internal/domain"
	invpg "inventory-service/pkg/postgres"
)

// WarehouseRepo implements warehouse persistence in PostgreSQL.
type WarehouseRepo struct {
	db *sql.DB
}

// NewWarehouseRepo creates a new WarehouseRepo.
func NewWarehouseRepo(db *sql.DB) *WarehouseRepo {
	return &WarehouseRepo{db: db}
}

// FindByID retrieves a warehouse by its UUID.
func (r *WarehouseRepo) FindByID(ctx context.Context, id string) (*domain.Warehouse, error) {
	exec := invpg.GetExecutor(ctx, r.db)

	var w domain.Warehouse
	err := exec.QueryRowContext(ctx, `
		SELECT id, code, name, region, address_line1, address_ward, address_district, address_city, lat, lng, is_active, created_at, updated_at
		FROM warehouses
		WHERE id = $1
	`, id).Scan(
		&w.ID, &w.Code, &w.Name, &w.Region,
		&w.Address.Line1, &w.Address.Ward, &w.Address.District, &w.Address.City,
		&w.Address.Lat, &w.Address.Lng,
		&w.IsActive, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("WarehouseRepo.FindByID: %w", err)
	}

	return &w, nil
}

func (r *WarehouseRepo) Find(ctx context.Context) ([]domain.Warehouse, error) {
	exec := invpg.GetExecutor(ctx, r.db)

	rows, err := exec.QueryContext(ctx, `
		SELECT id, code, name, region, address_line1, address_ward, address_district, address_city, lat, lng, is_active, created_at, updated_at
		FROM warehouses
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("WarehouseRepo.Find - query: %w", err)
	}
	defer rows.Close()

	var warehouses []domain.Warehouse
	for rows.Next() {
		var w domain.Warehouse
		if err := rows.Scan(
			&w.ID, &w.Code, &w.Name, &w.Region,
			&w.Address.Line1, &w.Address.Ward, &w.Address.District, &w.Address.City,
			&w.Address.Lat, &w.Address.Lng,
			&w.IsActive, &w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("WarehouseRepo.Find - scan: %w", err)
		}
		warehouses = append(warehouses, w)
	}
	return warehouses, rows.Err()
}
