package usecase

import (
	"context"

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
}

var _ OrderRepository = (*postgres.OrderRepo)(nil)

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

var _ Transactor = (*pgtransactor.PostgresTransactor)(nil)
