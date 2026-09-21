package usecase_test

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"inventory-service/internal/domain"
	"inventory-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type mockPORepo struct {
	po *domain.PurchaseOrder
}

func (m *mockPORepo) Create(ctx context.Context, po *domain.PurchaseOrder) error { return nil }
func (m *mockPORepo) FindByID(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	return m.po, nil
}
func (m *mockPORepo) Update(ctx context.Context, po *domain.PurchaseOrder) error {
	m.po = po
	return nil
}
func (m *mockPORepo) List(ctx context.Context, filter domain.PurchaseOrderFilter, page pagination.Params) ([]domain.PurchaseOrder, int64, error) {
	return nil, 0, nil
}

type mockStockLevelRepo struct {
	level *domain.StockLevel
}

func (m *mockStockLevelRepo) Create(ctx context.Context, sl *domain.StockLevel) error { return nil }
func (m *mockStockLevelRepo) FindByID(ctx context.Context, id string) (*domain.StockLevel, error) {
	return nil, nil
}
func (m *mockStockLevelRepo) FindBySKU(ctx context.Context, sku string) ([]domain.StockLevel, error) {
	return nil, nil
}
func (m *mockStockLevelRepo) FindBySKUAndWarehouse(ctx context.Context, sku, warehouseID string) (*domain.StockLevel, error) {
	return m.level, nil
}
func (m *mockStockLevelRepo) FindFirstAvailable(ctx context.Context, sku string, qty int) (*domain.StockLevel, error) {
	return nil, nil
}
func (m *mockStockLevelRepo) List(ctx context.Context, filter domain.StockLevelFilter, page pagination.Params) ([]domain.StockLevel, int64, error) {
	return nil, 0, nil
}
func (m *mockStockLevelRepo) GetSummary(ctx context.Context, warehouseID string) (*domain.StockSummary, error) {
	return nil, nil
}
func (m *mockStockLevelRepo) Update(ctx context.Context, sl *domain.StockLevel) error {
	m.level = sl
	return nil
}
func (m *mockStockLevelRepo) BulkCheckAvailability(ctx context.Context, skus []string) (map[string]int, error) {
	return nil, nil
}

type mockMovementRepo struct {
	movements []*domain.StockMovement
}

func (m *mockMovementRepo) Create(ctx context.Context, sm *domain.StockMovement) error {
	m.movements = append(m.movements, sm)
	return nil
}
func (m *mockMovementRepo) List(ctx context.Context, filter domain.StockMovementFilter, page pagination.Params) ([]domain.StockMovement, int64, error) {
	return nil, 0, nil
}

type mockOutboxRepo struct {
	mu     sync.Mutex
	events []*domain.OutboxEvent
}

func (m *mockOutboxRepo) Create(ctx context.Context, event *domain.OutboxEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
	return nil
}

type mockTransactor struct{}

func (m *mockTransactor) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestReceiveGoods_TransactionalOutbox(t *testing.T) {
	po := &domain.PurchaseOrder{
		ID:            "po-uuid-1",
		Code:          "PO-2026-001",
		SupplierID:    "supp-1",
		SupplierCode:  "SUPP-01",
		WarehouseID:   "wh-1",
		WarehouseCode: "WH-NORTH",
		Status:        domain.POStatusOrdered,
		Lines: []domain.PurchaseOrderLine{
			{
				SKU:              "SKU-IPHONE-15",
				ProductName:      "iPhone 15",
				QuantityOrdered:  10,
				QuantityReceived: 0,
				UnitCost:         20000000,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	stockLevel := &domain.StockLevel{
		ID:          "sl-1",
		SKU:         "SKU-IPHONE-15",
		WarehouseID: "wh-1",
		OnHand:      5,
		Reserved:    0,
		Version:     1,
	}

	poRepo := &mockPORepo{po: po}
	stockLevelRepo := &mockStockLevelRepo{level: stockLevel}
	movRepo := &mockMovementRepo{}
	outboxRepo := &mockOutboxRepo{}
	transactor := &mockTransactor{}
	l := logger.New("error")

	uc := usecase.NewPurchaseOrderUC(
		poRepo,
		nil,
		nil,
		stockLevelRepo,
		movRepo,
		nil,
		transactor,
		outboxRepo,
		l,
	)

	// Receive 6 units
	receiveLines := []domain.ReceiveLine{
		{
			SKU:      "SKU-IPHONE-15",
			Quantity: 6,
		},
	}

	updatedPO, err := uc.ReceiveGoods(context.Background(), po.ID, receiveLines)
	if err != nil {
		t.Fatalf("ReceiveGoods failed: %v", err)
	}

	// Verify PO updated
	if updatedPO.Status != domain.POStatusPartiallyReceived {
		t.Errorf("expected status %s, got %s", domain.POStatusPartiallyReceived, updatedPO.Status)
	}
	if updatedPO.Lines[0].QuantityReceived != 6 {
		t.Errorf("expected QuantityReceived 6, got %d", updatedPO.Lines[0].QuantityReceived)
	}

	// Verify outbox event written
	if len(outboxRepo.events) != 1 {
		t.Fatalf("expected 1 outbox event, got %d", len(outboxRepo.events))
	}

	outboxEvt := outboxRepo.events[0]
	if outboxEvt.AggregateType != domain.AggregatePurchaseOrder {
		t.Errorf("expected aggregate type %s, got %s", domain.AggregatePurchaseOrder, outboxEvt.AggregateType)
	}
	if outboxEvt.AggregateID != po.ID {
		t.Errorf("expected aggregate ID %s, got %s", po.ID, outboxEvt.AggregateID)
	}
	if outboxEvt.EventType != domain.EventGoodsReceived {
		t.Errorf("expected event type %s, got %s", domain.EventGoodsReceived, outboxEvt.EventType)
	}

	var payload domain.GoodsReceivedEvent
	if err := json.Unmarshal(outboxEvt.Payload, &payload); err != nil {
		t.Fatalf("failed to unmarshal outbox payload: %v", err)
	}

	// Process event via StockUC as event consumer would
	stockUC := usecase.NewStockUC(stockLevelRepo, nil, nil, movRepo, nil, transactor, nil, l)
	if err := stockUC.ApplyGoodsReceived(context.Background(), payload); err != nil {
		t.Fatalf("ApplyGoodsReceived failed: %v", err)
	}

	// Verify stock level adjusted
	if stockLevel.OnHand != 11 {
		t.Errorf("expected OnHand 11, got %d", stockLevel.OnHand)
	}

	// Verify movement recorded
	if len(movRepo.movements) != 1 {
		t.Fatalf("expected 1 stock movement, got %d", len(movRepo.movements))
	}
	if movRepo.movements[0].Quantity != 6 {
		t.Errorf("expected movement quantity 6, got %d", movRepo.movements[0].Quantity)
	}
}
