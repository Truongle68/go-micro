package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"inventory-service/internal/domain"
	invpg "inventory-service/pkg/postgres"
)

// StockReservationRepo implements stock reservation persistence.
type StockReservationRepo struct {
	db *sql.DB
}

// NewStockReservationRepo creates a new StockReservationRepo.
func NewStockReservationRepo(db *sql.DB) *StockReservationRepo {
	return &StockReservationRepo{db: db}
}

// Create inserts a new stock reservation.
func (r *StockReservationRepo) Create(ctx context.Context, res *domain.StockReservation) error {
	exec := invpg.GetExecutor(ctx, r.db)

	err := exec.QueryRowContext(ctx, `
		INSERT INTO stock_reservations (sku, warehouse_id, quantity, order_id, status, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`,
		res.SKU, res.WarehouseID, res.Quantity,
		res.OrderID, string(res.Status), res.ExpiresAt, res.CreatedAt,
	).Scan(&res.ID)
	if err != nil {
		return fmt.Errorf("StockReservationRepo.Create: %w", err)
	}

	return nil
}

// FindByOrderID returns all reservations for a given order ID.
func (r *StockReservationRepo) FindByOrderID(ctx context.Context, orderID string) ([]domain.StockReservation, error) {
	exec := invpg.GetExecutor(ctx, r.db)

	rows, err := exec.QueryContext(ctx, `
		SELECT id, sku, warehouse_id, quantity, order_id, status,
		       expires_at, created_at, confirmed_at, released_at
		FROM stock_reservations
		WHERE order_id = $1
		ORDER BY created_at
	`, orderID)
	if err != nil {
		return nil, fmt.Errorf("StockReservationRepo.FindByOrderID - query: %w", err)
	}
	defer rows.Close()

	var reservations []domain.StockReservation
	for rows.Next() {
		var res domain.StockReservation
		var status string
		if err := rows.Scan(
			&res.ID, &res.SKU, &res.WarehouseID, &res.Quantity,
			&res.OrderID, &status, &res.ExpiresAt, &res.CreatedAt,
			&res.ConfirmedAt, &res.ReleasedAt,
		); err != nil {
			return nil, fmt.Errorf("StockReservationRepo.FindByOrderID - scan: %w", err)
		}
		res.Status = domain.ReservationStatus(status)
		reservations = append(reservations, res)
	}

	return reservations, rows.Err()
}

// FindPendingByOrderID returns only pending reservations for a given order ID.
func (r *StockReservationRepo) FindPendingByOrderID(ctx context.Context, orderID string) ([]domain.StockReservation, error) {
	exec := invpg.GetExecutor(ctx, r.db)

	rows, err := exec.QueryContext(ctx, `
		SELECT id, sku, warehouse_id, quantity, order_id, status,
		       expires_at, created_at, confirmed_at, released_at
		FROM stock_reservations
		WHERE order_id = $1 AND status = 'pending'
		ORDER BY created_at
	`, orderID)
	if err != nil {
		return nil, fmt.Errorf("StockReservationRepo.FindPendingByOrderID - query: %w", err)
	}
	defer rows.Close()

	var reservations []domain.StockReservation
	for rows.Next() {
		var res domain.StockReservation
		var status string
		if err := rows.Scan(
			&res.ID, &res.SKU, &res.WarehouseID, &res.Quantity,
			&res.OrderID, &status, &res.ExpiresAt, &res.CreatedAt,
			&res.ConfirmedAt, &res.ReleasedAt,
		); err != nil {
			return nil, fmt.Errorf("StockReservationRepo.FindPendingByOrderID - scan: %w", err)
		}
		res.Status = domain.ReservationStatus(status)
		reservations = append(reservations, res)
	}

	return reservations, rows.Err()
}

// Update persists changes to a stock reservation (status, confirmed_at, released_at).
func (r *StockReservationRepo) Update(ctx context.Context, res *domain.StockReservation) error {
	exec := invpg.GetExecutor(ctx, r.db)

	_, err := exec.ExecContext(ctx, `
		UPDATE stock_reservations
		SET status = $1, confirmed_at = $2, released_at = $3
		WHERE id = $4
	`,
		string(res.Status), res.ConfirmedAt, res.ReleasedAt, res.ID,
	)
	if err != nil {
		return fmt.Errorf("StockReservationRepo.Update: %w", err)
	}

	return nil
}

// FindExpiredPending returns all pending reservations that have passed their expiry time.
func (r *StockReservationRepo) FindExpiredPending(ctx context.Context) ([]domain.StockReservation, error) {
	exec := invpg.GetExecutor(ctx, r.db)

	rows, err := exec.QueryContext(ctx, `
		SELECT id, sku, warehouse_id, quantity, order_id, status,
		       expires_at, created_at, confirmed_at, released_at
		FROM stock_reservations
		WHERE status = 'pending' AND expires_at < NOW()
		ORDER BY expires_at
	`)
	if err != nil {
		return nil, fmt.Errorf("StockReservationRepo.FindExpiredPending - query: %w", err)
	}
	defer rows.Close()

	var reservations []domain.StockReservation
	for rows.Next() {
		var res domain.StockReservation
		var status string
		if err := rows.Scan(
			&res.ID, &res.SKU, &res.WarehouseID, &res.Quantity,
			&res.OrderID, &status, &res.ExpiresAt, &res.CreatedAt,
			&res.ConfirmedAt, &res.ReleasedAt,
		); err != nil {
			return nil, fmt.Errorf("StockReservationRepo.FindExpiredPending - scan: %w", err)
		}
		res.Status = domain.ReservationStatus(status)
		reservations = append(reservations, res)
	}

	return reservations, rows.Err()
}
