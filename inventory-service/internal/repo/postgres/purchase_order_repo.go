package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"inventory-service/internal/domain"
	invpg "inventory-service/pkg/postgres"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type PurchaseOrderRepo struct {
	db *sql.DB
}

func NewPurchaseOrderRepo(db *sql.DB) *PurchaseOrderRepo {
	return &PurchaseOrderRepo{db: db}
}

func (r *PurchaseOrderRepo) Create(ctx context.Context, po *domain.PurchaseOrder) error {
	executor := invpg.GetExecutor(ctx, r.db)

	if po.ID == "" {
		po.ID = uuid.New().String()
	}

	linesJSON, err := json.Marshal(po.Lines)
	if err != nil {
		return fmt.Errorf("PurchaseOrderRepo.Create - marshal lines: %w", err)
	}

	query := `
		INSERT INTO purchase_orders (
			id, version, code, supplier_id, supplier_name, warehouse_id,
			lines, status, created_by, expected_at, received_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = executor.ExecContext(ctx, query,
		po.ID, po.Version, po.Code, po.SupplierID, po.SupplierName, po.WarehouseID,
		linesJSON, string(po.Status), po.CreatedBy, po.ExpectedAt, po.ReceivedAt,
		po.CreatedAt, po.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrDuplicatePOCode
		}
		if strings.Contains(err.Error(), "purchase_orders_code_key") || (strings.Contains(err.Error(), "23505") && strings.Contains(err.Error(), "code")) {
			return domain.ErrDuplicatePOCode
		}
		return fmt.Errorf("PurchaseOrderRepo.Create: %w", err)
	}
	return nil
}

func (r *PurchaseOrderRepo) FindByID(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	executor := invpg.GetExecutor(ctx, r.db)

	query := `
		SELECT
			id, version, code, supplier_id, supplier_name, warehouse_id,
			lines, status, created_by, expected_at, received_at,
			created_at, updated_at
		FROM purchase_orders
		WHERE id = $1`

	var po domain.PurchaseOrder
	var linesJSON []byte
	var statusStr string

	err := executor.QueryRowContext(ctx, query, id).Scan(
		&po.ID, &po.Version, &po.Code, &po.SupplierID, &po.SupplierName, &po.WarehouseID,
		&linesJSON, &statusStr, &po.CreatedBy, &po.ExpectedAt, &po.ReceivedAt,
		&po.CreatedAt, &po.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPONotFound
	}
	if err != nil {
		return nil, fmt.Errorf("PurchaseOrderRepo.FindByID: %w", err)
	}

	po.Status = domain.PurchaseOrderStatus(statusStr)
	if len(linesJSON) > 0 {
		if err := json.Unmarshal(linesJSON, &po.Lines); err != nil {
			return nil, fmt.Errorf("PurchaseOrderRepo.FindByID - unmarshal lines: %w", err)
		}
	}
	return &po, nil
}

func (r *PurchaseOrderRepo) Update(ctx context.Context, po *domain.PurchaseOrder) error {
	executor := invpg.GetExecutor(ctx, r.db)

	linesJSON, err := json.Marshal(po.Lines)
	if err != nil {
		return fmt.Errorf("PurchaseOrderRepo.Update - marshal lines: %w", err)
	}

	query := `
		UPDATE purchase_orders
		SET version = $1, lines = $2, status = $3,
			expected_at = $4, received_at = $5, updated_at = $6
		WHERE id = $7 AND version = $8 - 1`

	result, err := executor.ExecContext(ctx, query,
		po.Version, linesJSON, string(po.Status),
		po.ExpectedAt, po.ReceivedAt, po.UpdatedAt,
		po.ID, po.Version,
	)
	if err != nil {
		return fmt.Errorf("PurchaseOrderRepo.Update: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PurchaseOrderRepo.Update - RowsAffected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("PurchaseOrderRepo.Update: optimistic lock conflict or not found")
	}
	return nil
}

func (r *PurchaseOrderRepo) List(ctx context.Context, filter domain.PurchaseOrderFilter, page pagination.Params) ([]domain.PurchaseOrder, error) {
	executor := invpg.GetExecutor(ctx, r.db)

	normParams := page.Normalize()

	var conditions []string
	args := make([]interface{}, 0, 5)
	argIdx := 1

	if filter.SupplierID != "" {
		conditions = append(conditions, fmt.Sprintf("supplier_id = $%d", argIdx))
		args = append(args, filter.SupplierID)
		argIdx++
	}
	if filter.WarehouseID != "" {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", argIdx))
		args = append(args, filter.WarehouseID)
		argIdx++
	}
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT
			id, version, code, supplier_id, supplier_name, warehouse_id,
			lines, status, created_by, expected_at, received_at,
			created_at, updated_at
		FROM purchase_orders
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, normParams.Limit, normParams.Skip())

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PurchaseOrderRepo.List: %w", err)
	}
	defer rows.Close()

	orders := make([]domain.PurchaseOrder, 0)
	for rows.Next() {
		var po domain.PurchaseOrder
		var linesJSON []byte
		var statusStr string

		if err := rows.Scan(
			&po.ID, &po.Version, &po.Code, &po.SupplierID, &po.SupplierName, &po.WarehouseID,
			&linesJSON, &statusStr, &po.CreatedBy, &po.ExpectedAt, &po.ReceivedAt,
			&po.CreatedAt, &po.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("PurchaseOrderRepo.List - scan: %w", err)
		}

		po.Status = domain.PurchaseOrderStatus(statusStr)
		if len(linesJSON) > 0 {
			_ = json.Unmarshal(linesJSON, &po.Lines)
		}
		orders = append(orders, po)
	}
	return orders, nil
}
