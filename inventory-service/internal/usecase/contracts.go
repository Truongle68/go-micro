package usecase

import (
	"context"
	"inventory-service/internal/client"
	"inventory-service/internal/domain"
	"inventory-service/internal/repo/postgres"
	invpg "inventory-service/pkg/postgres"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/TruongLe68/go-micro/pkg/rabbitmq/publisher"
)

type CatalogClient interface {
	GetVariantsBySKUs(ctx context.Context, skus []string) ([]client.VariantDTO, error)
}

type SupplierRepository interface {
	Create(ctx context.Context, s *domain.Supplier) error
	FindByID(ctx context.Context, id string) (*domain.Supplier, error)
	Update(ctx context.Context, s *domain.Supplier) error
	List(ctx context.Context, activeOnly bool, page pagination.Params) ([]domain.Supplier, error)
}

var _ SupplierRepository = (*postgres.SupplierRepo)(nil)

type PurchaseOrderRepository interface {
	Create(ctx context.Context, po *domain.PurchaseOrder) error
	FindByID(ctx context.Context, id string) (*domain.PurchaseOrder, error)
	Update(ctx context.Context, po *domain.PurchaseOrder) error
	List(ctx context.Context, filter domain.PurchaseOrderFilter, page pagination.Params) ([]domain.PurchaseOrder, error)
}

var _ PurchaseOrderRepository = (*postgres.PurchaseOrderRepo)(nil)

// StockLevelRepository defines persistence operations for stock levels.
type StockLevelRepository interface {
	Create(ctx context.Context, sl *domain.StockLevel) error
	FindBySKU(ctx context.Context, sku string) ([]domain.StockLevel, error)
	FindBySKUAndWarehouse(ctx context.Context, sku, warehouseID string) (*domain.StockLevel, error)
	FindFirstAvailable(ctx context.Context, sku string, qty int) (*domain.StockLevel, error)
	Update(ctx context.Context, sl *domain.StockLevel) error
	BulkCheckAvailability(ctx context.Context, skus []string) (map[string]int, error)
}

var _ StockLevelRepository = (*postgres.StockLevelRepo)(nil)

// WarehouseRepository defines persistence operations for warehouses.
type WarehouseRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Warehouse, error)
}

var _ WarehouseRepository = (*postgres.WarehouseRepo)(nil)

// StockReservationRepository defines persistence operations for stock reservations.
type StockReservationRepository interface {
	Create(ctx context.Context, res *domain.StockReservation) error
	FindByOrderID(ctx context.Context, orderID string) ([]domain.StockReservation, error)
	FindPendingByOrderID(ctx context.Context, orderID string) ([]domain.StockReservation, error)
	Update(ctx context.Context, res *domain.StockReservation) error
	FindExpiredPending(ctx context.Context) ([]domain.StockReservation, error)
}

var _ StockReservationRepository = (*postgres.StockReservationRepo)(nil)

// StockMovementRepository defines persistence operations for stock movements.
type StockMovementRepository interface {
	Create(ctx context.Context, m *domain.StockMovement) error
}

var _ StockMovementRepository = (*postgres.StockMovementRepo)(nil)

// EventPublisher abstracts the event publishing mechanism.
type EventPublisher interface {
	Publish(ctx context.Context, event publisher.Event) error
}

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

var _ Transactor = (*invpg.PostgresTransactor)(nil)
