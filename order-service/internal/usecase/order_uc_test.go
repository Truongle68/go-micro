package usecase_test

import (
	"context"
	"errors"
	"testing"

	"order-service/internal/client"
	"order-service/internal/domain"
	"order-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type mockOrderRepo struct {
	orders    map[string]*domain.Order
	histories map[string][]domain.OrderStatusHistory
}

func newMockOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{
		orders:    make(map[string]*domain.Order),
		histories: make(map[string][]domain.OrderStatusHistory),
	}
}

func (m *mockOrderRepo) Create(ctx context.Context, order *domain.Order, history *domain.OrderStatusHistory) error {
	m.orders[order.ID] = order
	if history != nil {
		m.histories[order.ID] = append(m.histories[order.ID], *history)
	}
	return nil
}

func (m *mockOrderRepo) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	o, ok := m.orders[id]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}
	return o, nil
}

func (m *mockOrderRepo) FindByUserID(ctx context.Context, userID string, limit int64, offset int64) ([]domain.Order, int64, error) {
	var userOrders []domain.Order
	for _, o := range m.orders {
		if o.UserID == userID {
			userOrders = append(userOrders, *o)
		}
	}
	return userOrders, int64(len(userOrders)), nil
}

func (m *mockOrderRepo) UpdateStatus(ctx context.Context, order *domain.Order, history *domain.OrderStatusHistory) error {
	m.orders[order.ID] = order
	if history != nil {
		m.histories[order.ID] = append(m.histories[order.ID], *history)
	}
	return nil
}

func (m *mockOrderRepo) GetTrackingHistory(ctx context.Context, orderID string) ([]domain.OrderStatusHistory, error) {
	h, ok := m.histories[orderID]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}
	return h, nil
}

type mockCatalogClient struct{}

func (m *mockCatalogClient) GetVariantsBySKUs(ctx context.Context, skus []string) ([]client.VariantDTO, error) {
	variants := make([]client.VariantDTO, len(skus))
	for i, sku := range skus {
		variants[i] = client.VariantDTO{
			ID:          "var_1",
			ProductID:   "prod_1",
			ProductName: "Test Product",
			SKU:         sku,
			Price:       client.Price{Amount: 500000, Currency: "VND"},
			IsActive:    true,
		}
	}
	return variants, nil
}

type mockTransactor struct{}

func (m *mockTransactor) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestCheckoutAndOrderLifecycle(t *testing.T) {
	repo := newMockOrderRepo()
	l := logger.New("error")
	uc := usecase.NewOrderUC(repo, nil, &mockCatalogClient{}, &mockTransactor{}, l)

	ctx := context.Background()
	userID := "usr_1001"

	input := usecase.CheckoutInput{
		Items: []usecase.CheckoutItemInput{
			{
				SKU:      "SKU-001",
				Quantity: 1,
			},
		},
		ShippingAddress: domain.AddressSnapshot{
			FullName: "Alice",
			Phone:    "0912345678",
			Street:   "456 High St",
			City:     "Hanoi",
		},
		ShippingFee:   30000,
		PaymentMethod: "cod",
	}

	// 1. Checkout -> Creates order and transitions to confirmed
	order, err := uc.Checkout(ctx, userID, input, "")
	if err != nil {
		t.Fatalf("unexpected error during checkout: %v", err)
	}

	if order.Status != domain.OrderStatusConfirmed {
		t.Fatalf("expected order status confirmed, got %s", order.Status)
	}
	if order.Total != 530000 {
		t.Fatalf("expected total 530000, got %d", order.Total)
	}

	// 2. Get Order
	fetchedOrder, err := uc.GetOrder(ctx, order.ID, userID)
	if err != nil {
		t.Fatalf("unexpected error getting order: %v", err)
	}
	if fetchedOrder.ID != order.ID {
		t.Fatalf("expected order ID %s, got %s", order.ID, fetchedOrder.ID)
	}

	// 3. Get Tracking Timeline
	timeline, err := uc.GetTrackingTimeline(ctx, order.ID, userID)
	if err != nil {
		t.Fatalf("unexpected error getting tracking timeline: %v", err)
	}
	if len(timeline) != 2 {
		t.Fatalf("expected 2 status entries (created & confirmed), got %d", len(timeline))
	}

	// 4. Ship Order
	shippedOrder, err := uc.ShipOrder(ctx, order.ID, "EXPRESS-12345")
	if err != nil {
		t.Fatalf("unexpected error shipping order: %v", err)
	}
	if shippedOrder.Status != domain.OrderStatusShipped || shippedOrder.TrackingCode != "EXPRESS-12345" {
		t.Fatalf("expected shipped status with tracking code, got %s, %s", shippedOrder.Status, shippedOrder.TrackingCode)
	}

	// 5. Deliver Order
	deliveredOrder, err := uc.DeliverOrder(ctx, order.ID)
	if err != nil {
		t.Fatalf("unexpected error delivering order: %v", err)
	}
	if deliveredOrder.Status != domain.OrderStatusDelivered {
		t.Fatalf("expected delivered status, got %s", deliveredOrder.Status)
	}

	// 6. List Orders
	result, err := uc.ListOrdersByUser(ctx, userID, pagination.Params{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error listing orders: %v", err)
	}
	if result.Meta.TotalCount != 1 {
		t.Fatalf("expected 1 total order, got %d", result.Meta.TotalCount)
	}

	// 7. Cancel delivered order should fail
	_, err = uc.CancelOrder(ctx, order.ID, userID, "Changed mind")
	if !errors.Is(err, domain.ErrCannotCancelDeliveriedOrder) {
		t.Fatalf("expected ErrCannotCancelDeliveriedOrder, got %v", err)
	}
}
