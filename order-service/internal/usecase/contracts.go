package usecase

import (
	"context"

	"order-service/internal/client"
	"order-service/internal/domain"
	"order-service/internal/repo/postgres"
	pgtransactor "order-service/pkg/postgres"
)

type CheckoutItemInput struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type CheckoutInput struct {
	Items           []CheckoutItemInput    `json:"items,omitempty"` // if empty, fetches cart from cart-service
	ShippingAddress domain.AddressSnapshot `json:"shipping_address"`
	ShippingFee     int64                  `json:"shipping_fee"`
	PaymentMethod   string                 `json:"payment_method"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order, history *domain.OrderStatusHistory) error
	FindByID(ctx context.Context, id string) (*domain.Order, error)
	FindByUserID(ctx context.Context, userID string, limit int64, offset int64) ([]domain.Order, int64, error)
	UpdateStatus(ctx context.Context, order *domain.Order, history *domain.OrderStatusHistory) error
	GetTrackingHistory(ctx context.Context, orderID string) ([]domain.OrderStatusHistory, error)
	AppendNote(ctx context.Context, orderID string, currentStatus domain.OrderStatus, note string) error
}

var _ OrderRepository = (*postgres.OrderRepo)(nil)

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

var _ Transactor = (*pgtransactor.PostgresTransactor)(nil)

// InventoryClient defines the inventory service operations needed by the order usecase.
type InventoryClient interface {
	CheckStock(ctx context.Context, items []client.SKUQty) (map[string]int, error)
	ReserveStock(ctx context.Context, orderID string, items []client.SKUQty) error
	ConfirmReservation(ctx context.Context, orderID string) error
	ReleaseReservation(ctx context.Context, orderID string) error
}
