package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"order-service/internal/domain"
	orderpg "order-service/pkg/postgres"

	"github.com/google/uuid"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (r *OrderRepo) Create(ctx context.Context, order *domain.Order, history *domain.OrderStatusHistory) error {
	executor := orderpg.GetExecutor(ctx, r.db)

	queryOrder := `
		INSERT INTO orders (
			id, user_id, status, subtotal, shipping_fee, total,
			payment_method, payment_id, tracking_code,
			full_name, phone, street, ward, district, city, country,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16,
			$17, $18
		)`

	_, err := executor.ExecContext(ctx, queryOrder,
		order.ID, order.UserID, string(order.Status), order.Subtotal, order.ShippingFee, order.Total,
		order.PaymentMethod, order.PaymentID, order.TrackingCode,
		order.ShippingAddress.FullName, order.ShippingAddress.Phone, order.ShippingAddress.Street,
		order.ShippingAddress.Ward, order.ShippingAddress.District, order.ShippingAddress.City, order.ShippingAddress.Country,
		order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("OrderRepo.Create - insert order: %w", err)
	}

	queryItem := `
		INSERT INTO order_items (
			id, order_id, product_id, variant_id, sku, product_name, image, variant_attrs, unit_price, quantity, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	for _, item := range order.Items {
		variantJSON, err := json.Marshal(item.VariantAttrs)
		if err != nil {
			return fmt.Errorf("OrderRepo.Create - marshal variant attrs: %w", err)
		}

		_, err = executor.ExecContext(ctx, queryItem,
			item.ID, item.OrderID, item.ProductID, item.VariantID, item.SKU, item.ProductName, item.Image, variantJSON, item.UnitPrice, item.Quantity, order.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("OrderRepo.Create - insert order item %s: %w", item.SKU, err)
		}
	}

	if history != nil {
		queryHist := `
			INSERT INTO order_status_history (
				id, order_id, from_status, to_status, note, created_at
			) VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = executor.ExecContext(ctx, queryHist,
			history.ID, history.OrderID, string(history.FromStatus), string(history.ToStatus), history.Note, history.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("OrderRepo.Create - insert status history: %w", err)
		}
	}

	return nil
}

func (r *OrderRepo) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	executor := orderpg.GetExecutor(ctx, r.db)

	queryOrder := `
		SELECT
			id, user_id, status, subtotal, shipping_fee, total,
			payment_method, payment_id, tracking_code,
			full_name, phone, street, ward, district, city, country,
			created_at, updated_at
		FROM orders
		WHERE id = $1`

	var o domain.Order
	var statusStr string
	err := executor.QueryRowContext(ctx, queryOrder, id).Scan(
		&o.ID, &o.UserID, &statusStr, &o.Subtotal, &o.ShippingFee, &o.Total,
		&o.PaymentMethod, &o.PaymentID, &o.TrackingCode,
		&o.ShippingAddress.FullName, &o.ShippingAddress.Phone, &o.ShippingAddress.Street,
		&o.ShippingAddress.Ward, &o.ShippingAddress.District, &o.ShippingAddress.City, &o.ShippingAddress.Country,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("OrderRepo.FindByID - query order: %w", err)
	}
	o.Status = domain.OrderStatus(statusStr)

	items, err := r.findItemsByOrderID(ctx, executor, o.ID)
	if err != nil {
		return nil, err
	}
	o.Items = items

	return &o, nil
}

func (r *OrderRepo) FindByUserID(ctx context.Context, userID string, limit int64, offset int64) ([]domain.Order, int64, error) {
	executor := orderpg.GetExecutor(ctx, r.db)

	var total int64
	countQuery := `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	if err := executor.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("OrderRepo.FindByUserID - count: %w", err)
	}
	if total == 0 {
		return []domain.Order{}, 0, nil
	}

	queryOrders := `
		SELECT
			id, user_id, status, subtotal, shipping_fee, total,
			payment_method, payment_id, tracking_code,
			full_name, phone, street, ward, district, city, country,
			created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := executor.QueryContext(ctx, queryOrders, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("OrderRepo.FindByUserID - query orders: %w", err)
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)
	for rows.Next() {
		var o domain.Order
		var statusStr string
		if err := rows.Scan(
			&o.ID, &o.UserID, &statusStr, &o.Subtotal, &o.ShippingFee, &o.Total,
			&o.PaymentMethod, &o.PaymentID, &o.TrackingCode,
			&o.ShippingAddress.FullName, &o.ShippingAddress.Phone, &o.ShippingAddress.Street,
			&o.ShippingAddress.Ward, &o.ShippingAddress.District, &o.ShippingAddress.City, &o.ShippingAddress.Country,
			&o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("OrderRepo.FindByUserID - scan: %w", err)
		}
		o.Status = domain.OrderStatus(statusStr)

		items, err := r.findItemsByOrderID(ctx, executor, o.ID)
		if err != nil {
			return nil, 0, err
		}
		o.Items = items

		orders = append(orders, o)
	}

	return orders, total, nil
}

func (r *OrderRepo) findItemsByOrderID(ctx context.Context, executor interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}, orderID string) ([]domain.OrderItem, error) {
	queryItems := `
		SELECT id, order_id, product_id, variant_id, sku, product_name, image, variant_attrs, unit_price, quantity
		FROM order_items
		WHERE order_id=$1
		ORDER BY created_at ASC`

	rows, err := executor.QueryContext(ctx, queryItems, orderID)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo.findItemsByOrderID - query items: %w", err)
	}
	defer rows.Close()

	items := make([]domain.OrderItem, 0)
	for rows.Next() {
		var item domain.OrderItem
		var variantJSON []byte
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.VariantID, &item.SKU, &item.ProductName, &item.Image, &variantJSON, &item.UnitPrice, &item.Quantity); err != nil {
			return nil, fmt.Errorf("OrderRepo.findItemsByOrderID - scan item: %w", err)
		}
		if len(variantJSON) > 0 {
			_ = json.Unmarshal(variantJSON, &item.VariantAttrs)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, order *domain.Order, history *domain.OrderStatusHistory) error {
	executor := orderpg.GetExecutor(ctx, r.db)

	queryUpdate := `
		UPDATE orders
		SET status = $1, payment_id = $2, tracking_code = $3, updated_at = $4
		WHERE id = $5`

	_, err := executor.ExecContext(ctx, queryUpdate,
		string(order.Status), order.PaymentID, order.TrackingCode, order.UpdatedAt, order.ID,
	)
	if err != nil {
		return fmt.Errorf("OrderRepo.UpdateStatus - update order: %w", err)
	}

	if history != nil {
		queryHist := `
			INSERT INTO order_status_history (
				id, order_id, from_status, to_status, note, created_at
			) VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = executor.ExecContext(ctx, queryHist,
			history.ID, history.OrderID, string(history.FromStatus), string(history.ToStatus), history.Note, history.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("OrderRepo.UpdateStatus - insert status history: %w", err)
		}
	}

	return nil
}

func (r *OrderRepo) GetTrackingHistory(ctx context.Context, orderID string) ([]domain.OrderStatusHistory, error) {
	executor := orderpg.GetExecutor(ctx, r.db)

	query := `
		SELECT id, order_id, from_status, to_status, note, created_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY created_at ASC`

	rows, err := executor.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo.GetTrackingHistory - query history: %w", err)
	}
	defer rows.Close()

	histories := make([]domain.OrderStatusHistory, 0)
	for rows.Next() {
		var h domain.OrderStatusHistory
		var fromStr, toStr string
		if err := rows.Scan(&h.ID, &h.OrderID, &fromStr, &toStr, &h.Note, &h.CreatedAt); err != nil {
			return nil, fmt.Errorf("OrderRepo.GetTrackingHistory - scan: %w", err)
		}
		h.FromStatus = domain.OrderStatus(fromStr)
		h.ToStatus = domain.OrderStatus(toStr)
		histories = append(histories, h)
	}

	return histories, nil
}

func (r *OrderRepo) AppendNote(ctx context.Context, orderID string, currentStatus domain.OrderStatus, note string) error {
	executor := orderpg.GetExecutor(ctx, r.db)
	query := `
		INSERT INTO order_status_history (
			id, order_id, from_status, to_status, note, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := executor.ExecContext(ctx, query,
		uuid.New().String(), orderID, currentStatus, currentStatus, note, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("OrderRepo.AppendNote - insert: %w", err)
	}
	return nil
}
