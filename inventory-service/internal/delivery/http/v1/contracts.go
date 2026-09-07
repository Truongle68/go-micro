package v1

import (
	"context"

	"inventory-service/internal/domain"
	"inventory-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

// SupplierUC defines the methods the delivery layer needs from the supplier usecase.
type SupplierUC interface {
	CreateSupplier(ctx context.Context, code, name, phone, email string, address domain.SupplierAddress) (*domain.Supplier, error)
	GetSupplier(ctx context.Context, id string) (*domain.Supplier, error)
	UpdateSupplier(ctx context.Context, id, name, email, phone string, address domain.SupplierAddress) (*domain.Supplier, error)
	DeactivateSupplier(ctx context.Context, id string) (*domain.Supplier, error)
	ReactivateSupplier(ctx context.Context, id string) (*domain.Supplier, error)
	ListSuppliers(ctx context.Context, activeOnly bool, page pagination.Params) ([]domain.Supplier, error)
}

var _ SupplierUC = (*usecase.SupplierUC)(nil)

// PurchaseOrderUC defines the methods the delivery layer needs from the purchase order usecase.
type PurchaseOrderUC interface {
	CreatePurchaseOrder(ctx context.Context, input usecase.CreatePurchaseOrderInput) (*domain.PurchaseOrder, error)
	GetPurchaseOrder(ctx context.Context, id string) (*domain.PurchaseOrder, error)
	ConfirmPurchaseOrder(ctx context.Context, id string) (*domain.PurchaseOrder, error)
	ReceiveLine(ctx context.Context, poID, sku string, qty int) (*domain.PurchaseOrder, error)
	CancelPurchaseOrder(ctx context.Context, id string) (*domain.PurchaseOrder, error)
	ListPurchaseOrders(ctx context.Context, filter domain.PurchaseOrderFilter, page pagination.Params) ([]domain.PurchaseOrder, error)
}

var _ PurchaseOrderUC = (*usecase.PurchaseOrderUC)(nil)

type WarehouseUC interface {
	List(ctx context.Context) ([]domain.Warehouse, error)
}

var _ WarehouseUC = (*usecase.WarehouseUC)(nil)

// StockUC defines the methods the HTTP delivery layer needs from the stock usecase.
type StockUC interface {
	CheckStock(ctx context.Context, skus []string) (map[string]int, error)
	GetSKUAvailability(ctx context.Context, sku string) (int, error)
	ListStockLevels(ctx context.Context, filter domain.StockLevelFilter, page pagination.Params) ([]usecase.DetailedStockLevel, int64, error)
	GetStockLevel(ctx context.Context, id string) (*usecase.DetailedStockLevel, error)
	AdjustStock(ctx context.Context, input usecase.AdjustStockInput) (*usecase.DetailedStockLevel, error)
	UpdateThresholds(ctx context.Context, id string, reorderThreshold, reorderQuantity int) (*usecase.DetailedStockLevel, error)
	GetStockSummary(ctx context.Context, warehouseID string) (*domain.StockSummary, error)
	ListMovements(ctx context.Context, filter domain.StockMovementFilter, page pagination.Params) ([]usecase.DetailedStockMovement, int64, error)
}

var _ StockUC = (*usecase.StockUC)(nil)
