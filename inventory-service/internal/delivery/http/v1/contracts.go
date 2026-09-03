package v1

import (
	"context"

	"inventory-service/internal/domain"
	"inventory-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

// SupplierUC defines the methods the delivery layer needs from the supplier usecase.
type SupplierUC interface {
	CreateSupplier(ctx context.Context, code, name, phone string, address domain.SupplierAddress) (*domain.Supplier, error)
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
