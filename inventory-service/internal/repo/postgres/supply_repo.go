package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"inventory-service/internal/domain"
	"strings"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type SupplyRepo struct {
	db *sql.DB
}

func NewSupplyRepo(db *sql.DB) *SupplyRepo {
	return &SupplyRepo{
		db: db,
	}
}

func (r *SupplyRepo) CreateBatch(ctx context.Context, supplies []*domain.Supply) error {
	if len(supplies) == 0 {
		return nil
	}

	query := `
	INSERT INTO supplies
		(supplier_id, sku, price, lead_time, min_order_qty, is_preferred)
	VALUES` + buildPlaceholder(len(supplies), 6) + `
	ON CONFLICT (supplier_id, sku) DO UPDATE SET
		price = EXCLUDED.price,
		lead_time = EXCLUDED.lead_time,
		min_order_qty = EXCLUDED.min_order_qty,
		is_preferred = EXCLUDED.is_preferred,
		updated_at = NOW()`

	args := make([]any, 0, len(supplies)*6)
	for _, s := range supplies {
		rawLeadTime, err := json.Marshal(s.LeadTime)
		if err != nil {
			return fmt.Errorf("SupplyRepo.CreateBatch - marshal lead time: %w", err)
		}
		args = append(args, s.SupplierID, s.SKU, s.Price, rawLeadTime, s.MinOrderQty, s.IsPreferred)
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("SupplyRepo.CreateBatch - error: %w", err)
	}
	return nil
}

func (r *SupplyRepo) FindByID(ctx context.Context, id string) (*domain.Supply, error) {
	query := `
		SELECT id, supplier_id, sku, price::BIGINT AS price, lead_time, min_order_qty, is_preferred, created_at, updated_at
		FROM supplies
		WHERE id=$1
	`

	var s domain.Supply
	var rawLeadTime []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.SupplierID, &s.SKU, &s.Price, &rawLeadTime, &s.MinOrderQty, &s.IsPreferred, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSupplyNotFound
		}
		return nil, fmt.Errorf("SupplyRepo.FindByID - error: %w", err)
	}

	if err := json.Unmarshal(rawLeadTime, &s.LeadTime); err != nil {
		return nil, fmt.Errorf("SupplyRepo.FindByID - unmarshal lead time: %w", err)
	}

	return &s, nil
}

func (r *SupplyRepo) FindBySupplierID(ctx context.Context, supplierID string, p pagination.Params) ([]domain.Supply, int64, error) {
	normPage := p.Normalize()
	whereClause := " WHERE supplier_id = $1 "
	countQuery := "SELECT COUNT(*) FROM supplies" + whereClause
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, supplierID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("SupplyRepo.FindBySupplierID - count: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, supplier_id, sku, price::BIGINT AS price, lead_time, min_order_qty, is_preferred, created_at, updated_at
		FROM supplies
		%s
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, supplierID, normPage.Limit, normPage.Skip())
	if err != nil {
		return nil, 0, fmt.Errorf("SupplyRepo.FindBySupplierID - error: %w", err)
	}
	defer rows.Close()

	supplies := make([]domain.Supply, 0)
	for rows.Next() {
		var rawLeadTime []byte
		var s domain.Supply

		if err := rows.Scan(
			&s.ID,
			&s.SupplierID,
			&s.SKU,
			&s.Price,
			&rawLeadTime,
			&s.MinOrderQty,
			&s.IsPreferred,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("SupplyRepo.FindBySupplierID - scan:%w", err)
		}

		_ = json.Unmarshal(rawLeadTime, &s.LeadTime)
		supplies = append(supplies, s)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("SupplyRepo.FindBySupplierID - rows: %w", err)
	}

	return supplies, total, nil
}

// ($1,$2,$3),($4,$5,$6),..
func buildPlaceholder(rows, cols int) string {
	var sb strings.Builder
	arg := 1
	for i := 0; i < rows; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString("(")
		for j := 0; j < cols; j++ {
			if j > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(fmt.Sprintf("$%d", arg))
			arg++
		}
		sb.WriteString(")")
	}
	return sb.String()
}
