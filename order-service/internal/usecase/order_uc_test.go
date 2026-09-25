package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

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

func (m *mockOrderRepo) List(ctx context.Context, filter domain.OrderFilter, p pagination.Params) ([]domain.Order, int64, error) {
	var list []domain.Order
	for _, o := range m.orders {
		if filter.Status != "" && o.Status != filter.Status {
			continue
		}
		if filter.UserID != "" && o.UserID != filter.UserID {
			continue
		}
		if filter.SKU != "" {
			hasSKU := false
			for _, item := range o.Items {
				if item.SKU == filter.SKU {
					hasSKU = true
					break
				}
			}
			if !hasSKU {
				continue
			}
		}
		list = append(list, *o)
	}
	return list, int64(len(list)), nil
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

func (m *mockOrderRepo) AppendNote(ctx context.Context, orderID string, currentStatus domain.OrderStatus, note string) error {
	m.histories[orderID] = append(m.histories[orderID], domain.OrderStatusHistory{
		ID:         "note_id",
		OrderID:    orderID,
		FromStatus: currentStatus,
		ToStatus:   currentStatus,
		Note:       note,
		CreatedAt:  time.Now().UTC(),
	})
	return nil
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

type mockOutboxRepo struct{}

func (m *mockOutboxRepo) Create(ctx context.Context, event *domain.OutboxEvent) error {
	return nil
}

type mockUserClient struct{}

func (m *mockUserClient) GetProfile(ctx context.Context, userID string) (*client.UserProfileDTO, error) {
	return &client.UserProfileDTO{
		UserID:          userID,
		Email:           "example@gmail.com",
		Phone:           "099999999",
		FullName:        "Nguyen Van A",
		Status:          "verified",
		IsEmailVerified: true,
	}, nil
}

type mockCartClient struct{}

func (m *mockCartClient) GetCart(ctx context.Context, userID string, token string) (*client.CartDTO, error) {
	return &client.CartDTO{
		UserID: userID,
		Items: []client.CartItemDTO{
			{
				SKU:      "SKU-001",
				Quantity: 2,
			},
		},
	}, nil
}

func (m *mockCartClient) RemoveItems(ctx context.Context, userID string, skus []string, token string) error {
	return nil
}

func (m *mockCartClient) ClearCart(ctx context.Context, userID string, token string) error {
	return nil
}

type mockInventoryClient struct{}

func (m *mockInventoryClient) CheckStock(ctx context.Context, items []client.SKUQty) (map[string]int, error) {
	res := make(map[string]int, len(items))
	for _, item := range items {
		res[item.SKU] = 100
	}
	return res, nil
}

func (m *mockInventoryClient) ReserveStock(ctx context.Context, orderID string, items []client.SKUQty) error {
	return nil
}

func (m *mockInventoryClient) ConfirmReservation(ctx context.Context, orderID string) error {
	return nil
}

func (m *mockInventoryClient) ReleaseReservation(ctx context.Context, orderID string) error {
	return nil
}

func TestCheckoutAndOrderLifecycle(t *testing.T) {
	repo := newMockOrderRepo()
	l := logger.New("error")
	uc := usecase.NewOrderUC(repo, &mockOutboxRepo{}, &mockUserClient{}, &mockCartClient{}, &mockCatalogClient{}, &mockInventoryClient{}, &mockTransactor{}, l)

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

func TestOrderUC_ListOrders(t *testing.T) {
	repo := newMockOrderRepo()
	uc := usecase.NewOrderUC(repo, nil, nil, nil, nil, nil, nil, logger.New("error"))

	ctx := context.Background()
	now := time.Now().UTC()

	o1 := &domain.Order{
		ID:     "order_1",
		UserID: "user_1",
		Status: domain.OrderStatusConfirmed,
		Items: []domain.OrderItem{
			{ID: "item_1", SKU: "SKU-RED-M", Quantity: 1},
		},
		CreatedAt: now,
	}
	o2 := &domain.Order{
		ID:     "order_2",
		UserID: "user_2",
		Status: domain.OrderStatusDelivered,
		Items: []domain.OrderItem{
			{ID: "item_2", SKU: "SKU-BLUE-L", Quantity: 2},
		},
		CreatedAt: now,
	}
	_ = repo.Create(ctx, o1, nil)
	_ = repo.Create(ctx, o2, nil)

	// Filter by status
	res, err := uc.ListOrders(ctx, domain.OrderFilter{Status: domain.OrderStatusConfirmed}, pagination.Params{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Items) != 1 || res.Items[0].ID != "order_1" {
		t.Fatalf("expected 1 order with ID order_1, got %v", res.Items)
	}

	// Filter by user ID
	res, err = uc.ListOrders(ctx, domain.OrderFilter{UserID: "user_2"}, pagination.Params{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Items) != 1 || res.Items[0].ID != "order_2" {
		t.Fatalf("expected 1 order with ID order_2, got %v", res.Items)
	}

	// Filter by SKU
	res, err = uc.ListOrders(ctx, domain.OrderFilter{SKU: "SKU-RED-M"}, pagination.Params{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Items) != 1 || res.Items[0].ID != "order_1" {
		t.Fatalf("expected 1 order with ID order_1, got %v", res.Items)
	}
}
