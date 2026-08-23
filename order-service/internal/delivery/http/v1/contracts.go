package v1

import (
	"context"
	"order-service/internal/domain"
	"order-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type OrderUC interface {
	Checkout(ctx context.Context, userID string, input usecase.CheckoutInput, token string) (*domain.Order, error)
	GetOrder(ctx context.Context, orderID string, userID string) (*domain.Order, error)
	ListOrdersByUser(ctx context.Context, userID string, page pagination.Params) (*pagination.Result[domain.Order], error)
	GetTrackingTimeline(ctx context.Context, orderID string, userID string) ([]domain.OrderStatusHistory, error)
	ShipOrder(ctx context.Context, orderID string, trackingCode string) (*domain.Order, error)
	DeliverOrder(ctx context.Context, orderID string) (*domain.Order, error)
	CancelOrder(ctx context.Context, orderID string, userID string, reason string) (*domain.Order, error)
}

var _ OrderUC = (*usecase.OrderUC)(nil)
