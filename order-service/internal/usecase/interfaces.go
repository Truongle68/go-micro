package usecase

import (
	"context"

	"order-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type CheckoutItemInput struct {
	SKU          string            `json:"sku"`
	ProductName  string            `json:"product_name"`
	VariantAttrs map[string]string `json:"variant_attrs,omitempty"`
	UnitPrice    int64             `json:"unit_price"`
	Quantity     int               `json:"quantity"`
}

type CheckoutInput struct {
	Items           []CheckoutItemInput    `json:"items,omitempty"` // if empty, fetches cart from cart-service
	ShippingAddress domain.AddressSnapshot `json:"shipping_address"`
	ShippingFee     int64                  `json:"shipping_fee"`
	PaymentMethod   string                 `json:"payment_method"`
}

type Order interface {
	Checkout(ctx context.Context, userID string, input CheckoutInput, token string) (*domain.Order, error)
	GetOrder(ctx context.Context, orderID string, userID string) (*domain.Order, error)
	ListOrdersByUser(ctx context.Context, userID string, page pagination.Params) (*pagination.Result[domain.Order], error)
	GetTrackingTimeline(ctx context.Context, orderID string, userID string) ([]domain.OrderStatusHistory, error)
	ShipOrder(ctx context.Context, orderID string, trackingCode string) (*domain.Order, error)
	DeliverOrder(ctx context.Context, orderID string) (*domain.Order, error)
	CancelOrder(ctx context.Context, orderID string, userID string, reason string) (*domain.Order, error)
}
