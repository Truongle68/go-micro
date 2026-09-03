package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"inventory-service/internal/domain"
	invpg "inventory-service/pkg/postgres"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)

type SupplierRepo struct {
	db *sql.DB
}

func NewSupplierRepo(db *sql.DB) *SupplierRepo {
	return &SupplierRepo{db: db}
}

func (r *SupplierRepo) Create(ctx context.Context, s *domain.Supplier) error {
	executor := invpg.GetExecutor(ctx, r.db)

	if s.ID == "" {
		s.ID = uuid.New().String()
	}

	query := `
		INSERT INTO suppliers (
			id, version, code, name, email, phone,
			address_line1, address_ward, address_district, address_city,
			is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := executor.ExecContext(ctx, query,
		s.ID, s.Version, s.Code, s.Name, s.Email, s.Phone,
		s.Address.Line1, s.Address.Ward, s.Address.District, s.Address.City,
		s.IsActive, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrDuplicateSuppCode
		}
		if strings.Contains(err.Error(), "suppliers_code_key") || (strings.Contains(err.Error(), "23505") && strings.Contains(err.Error(), "code")) {
			return domain.ErrDuplicateSuppCode
		}
		return fmt.Errorf("SupplierRepo.Create: %w", err)
	}
	return nil
}

func (r *SupplierRepo) FindByID(ctx context.Context, id string) (*domain.Supplier, error) {
	executor := invpg.GetExecutor(ctx, r.db)

	query := `
		SELECT
			id, version, code, name, email, phone,
			address_line1, address_ward, address_district, address_city,
			is_active, created_at, updated_at
		FROM suppliers
		WHERE id = $1`

	var s domain.Supplier
	err := executor.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.Version, &s.Code, &s.Name, &s.Email, &s.Phone,
		&s.Address.Line1, &s.Address.Ward, &s.Address.District, &s.Address.City,
		&s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSuppNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("SupplierRepo.FindByID: %w", err)
	}
	return &s, nil
}

func (r *SupplierRepo) Update(ctx context.Context, s *domain.Supplier) error {
	executor := invpg.GetExecutor(ctx, r.db)

	query := `
		UPDATE suppliers
		SET version = $1, name = $2, email = $3, phone = $4,
			address_line1 = $5, address_ward = $6, address_district = $7, address_city = $8,
			is_active = $9, updated_at = $10
		WHERE id = $11 AND version = $12 - 1`

	result, err := executor.ExecContext(ctx, query,
		s.Version, s.Name, s.Email, s.Phone,
		s.Address.Line1, s.Address.Ward, s.Address.District, s.Address.City,
		s.IsActive, s.UpdatedAt,
		s.ID, s.Version,
	)
	if err != nil {
		return fmt.Errorf("SupplierRepo.Update: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("SupplierRepo.Update - RowsAffected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("SupplierRepo.Update: optimistic lock conflict or not found")
	}
	return nil
}

func (r *SupplierRepo) List(ctx context.Context, activeOnly bool, page pagination.Params) ([]domain.Supplier, error) {
	executor := invpg.GetExecutor(ctx, r.db)

	normParams := page.Normalize()

	baseWhere := ""
	args := make([]interface{}, 0, 4)
	argIdx := 1

	if activeOnly {
		baseWhere = fmt.Sprintf(" WHERE is_active = $%d", argIdx)
		args = append(args, true)
		argIdx++
	}

	query := fmt.Sprintf(`
		SELECT
			id, version, code, name, email, phone,
			address_line1, address_ward, address_district, address_city,
			is_active, created_at, updated_at
		FROM suppliers
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, baseWhere, argIdx, argIdx+1)

	args = append(args, normParams.Limit, normParams.Skip())

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("SupplierRepo.List: %w", err)
	}
	defer rows.Close()

	suppliers := make([]domain.Supplier, 0)
	for rows.Next() {
		var s domain.Supplier
		if err := rows.Scan(
			&s.ID, &s.Version, &s.Code, &s.Name, &s.Email, &s.Phone,
			&s.Address.Line1, &s.Address.Ward, &s.Address.District, &s.Address.City,
			&s.IsActive, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("SupplierRepo.List - scan: %w", err)
		}
		suppliers = append(suppliers, s)
	}
	return suppliers, nil
}
