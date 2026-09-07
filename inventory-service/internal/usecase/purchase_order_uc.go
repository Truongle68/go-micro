package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"inventory-service/internal/client"
	"inventory-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type CreatePurchaseOrderInput struct {
	Code          string
	SupplierID    string
	WarehouseID   string
	CreatedBy     string
	CreatedByName string
	Lines         []domain.NewPurchaseOrderLineInput
}

type PurchaseOrderUC struct {
	poRepo         PurchaseOrderRepository
	supplierRepo   SupplierRepository
	warehouseRepo  WarehouseRepository
	stockLevelRepo StockLevelRepository
	movRepo        StockMovementRepository
	catalogClient  CatalogClient
	transactor     Transactor
	logger         logger.Interface
}

func NewPurchaseOrderUC(
	poRepo PurchaseOrderRepository,
	supplierRepo SupplierRepository,
	warehouseRepo WarehouseRepository,
	stockLevelRepo StockLevelRepository,
	movRepo StockMovementRepository,
	catalogClient CatalogClient,
	transactor Transactor,
	l logger.Interface,
) *PurchaseOrderUC {
	return &PurchaseOrderUC{
		poRepo:         poRepo,
		supplierRepo:   supplierRepo,
		warehouseRepo:  warehouseRepo,
		stockLevelRepo: stockLevelRepo,
		movRepo:        movRepo,
		catalogClient:  catalogClient,
		transactor:     transactor,
		logger:         l,
	}
}

func (uc *PurchaseOrderUC) CreatePurchaseOrder(ctx context.Context, input CreatePurchaseOrderInput) (*domain.PurchaseOrder, error) {
	supplier, err := uc.supplierRepo.FindByID(ctx, input.SupplierID)
	if err != nil {
		return nil, fmt.Errorf("PurchaseOrderUC.CreatePurchaseOrder - find supplier: %w", err)
	}
	if supplier == nil {
		return nil, domain.ErrSuppNotFound
	}
	if !supplier.IsActive {
		return nil, fmt.Errorf("PurchaseOrderUC.CreatePurchaseOrder: %w", domain.ErrSuppAlreadyInactive)
	}

	var whCode, whName string
	if uc.warehouseRepo != nil {
		wh, err := uc.warehouseRepo.FindByID(ctx, input.WarehouseID)
		if err != nil {
			return nil, fmt.Errorf("PurchaseOrderUC.CreatePurchaseOrder - find warehouse: %w", err)
		}
		if wh == nil {
			return nil, domain.ErrWhNotFound
		}
		if !wh.IsActive {
			return nil, fmt.Errorf("PurchaseOrderUC.CreatePurchaseOrder: %w", domain.ErrWhAlreadyInactive)
		}
		whCode = wh.Code
		whName = wh.Name
	}

	// Validate SKUs against catalog-service and auto-populate ProductName
	if uc.catalogClient != nil && len(input.Lines) > 0 {
		skus := make([]string, len(input.Lines))
		for i, l := range input.Lines {
			skus[i] = l.SKU
		}

		variants, err := uc.catalogClient.GetVariantsBySKUs(ctx, skus)
		if err != nil {
			return nil, fmt.Errorf("PurchaseOrderUC.CreatePurchaseOrder - catalogClient.GetVariantsBySKUs: %w", err)
		}

		bySKU := make(map[string]client.VariantDTO, len(variants))
		for _, v := range variants {
			bySKU[v.SKU] = v
		}

		for i := range input.Lines {
			v, exists := bySKU[input.Lines[i].SKU]
			if !exists {
				return nil, fmt.Errorf("%w: %s", domain.ErrSKUNotFound, input.Lines[i].SKU)
			}
			if !v.IsActive {
				return nil, fmt.Errorf("%w: %s", domain.ErrInactiveVariant, input.Lines[i].SKU)
			}
			input.Lines[i].ProductName = v.ProductName
		}
	}

	po, err := domain.NewPurchaseOrder(domain.NewPurchaseOrderParams{
		Code:          input.Code,
		SupplierID:    supplier.ID,
		SupplierCode:  supplier.Code,
		SupplierName:  supplier.Name,
		WarehouseID:   input.WarehouseID,
		WarehouseCode: whCode,
		WarehouseName: whName,
		CreatedBy:     input.CreatedBy,
		CreatedByName: input.CreatedByName,
		Lines:         input.Lines,
	})
	if err != nil {
		return nil, err
	}

	if err := uc.poRepo.Create(ctx, po); err != nil {
		return nil, fmt.Errorf("PurchaseOrderUC.CreatePurchaseOrder - repo.Create: %w", err)
	}
	return po, nil
}

func (uc *PurchaseOrderUC) GetPurchaseOrder(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	po, err := uc.poRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (uc *PurchaseOrderUC) ConfirmPurchaseOrder(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	po, err := uc.poRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := po.Confirm(); err != nil {
		return nil, err
	}

	if err := uc.poRepo.Update(ctx, po); err != nil {
		return nil, fmt.Errorf("PurchaseOrderUC.ConfirmPurchaseOrder - repo.Update: %w", err)
	}
	return po, nil
}

func (uc *PurchaseOrderUC) ReceiveLine(ctx context.Context, poID, sku string, qty int) (*domain.PurchaseOrder, error) {
	var po *domain.PurchaseOrder

	err := uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		po, err = uc.poRepo.FindByID(txCtx, poID)
		if err != nil {
			return err
		}

		if err := po.ReceiveLine(sku, qty); err != nil {
			return err
		}

		if err := uc.poRepo.Update(txCtx, po); err != nil {
			return fmt.Errorf("updating purchase order: %w", err)
		}

		// Adjust stock atomically alongside the PO update
		if err := uc.applyStockReceipt(txCtx, po.WarehouseID, sku, qty, poID); err != nil {
			return fmt.Errorf("adjusting stock for %s: %w", sku, err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("PurchaseOrderUC.ReceiveLine: %w", err)
	}
	return po, nil
}

func (uc *PurchaseOrderUC) applyStockReceipt(ctx context.Context, warehouseID, sku string, qty int, poID string) error {
	const maxRetries = 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		level, err := uc.stockLevelRepo.FindBySKUAndWarehouse(ctx, sku, warehouseID)
		if err != nil {
			return err
		}
		if level == nil {
			// First-ever receipt for this SKU at this warehouse — create the row.
			newLevel, err := domain.NewStockLevel(domain.NewStockLevelParams{
				SKU:              sku,
				WarehouseID:      warehouseID,
				ReorderThreshold: 0,
				ReorderQuantity:  0,
			})
			if err != nil {
				return err
			}
			if err := newLevel.AdjustOnHand(qty); err != nil {
				return err
			}
			if err := uc.stockLevelRepo.Create(ctx, newLevel); err != nil {
				return err
			}
			if uc.movRepo != nil {
				_ = uc.movRepo.Create(ctx, &domain.StockMovement{
					SKU:           sku,
					WarehouseID:   warehouseID,
					Type:          domain.MovementInbound,
					Quantity:      qty,
					ReferenceType: "purchase_order",
					ReferenceID:   poID,
					Note:          fmt.Sprintf("Received %d units from PO %s", qty, poID),
					CreatedAt:     time.Now().UTC(),
				})
			}
			return nil
		}

		if err := level.AdjustOnHand(qty); err != nil {
			return err
		}
		err = uc.stockLevelRepo.Update(ctx, level)
		if err == nil {
			if uc.movRepo != nil {
				_ = uc.movRepo.Create(ctx, &domain.StockMovement{
					SKU:           sku,
					WarehouseID:   warehouseID,
					Type:          domain.MovementInbound,
					Quantity:      qty,
					ReferenceType: "purchase_order",
					ReferenceID:   poID,
					Note:          fmt.Sprintf("Received %d units from PO %s", qty, poID),
					CreatedAt:     time.Now().UTC(),
				})
			}
			return nil
		}
		if !errors.Is(err, domain.ErrConcurrentModification) {
			return err
		}
	}
	return domain.ErrConcurrentModification
}

func (uc *PurchaseOrderUC) CancelPurchaseOrder(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	po, err := uc.poRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := po.Cancel(); err != nil {
		return nil, err
	}

	if err := uc.poRepo.Update(ctx, po); err != nil {
		return nil, fmt.Errorf("PurchaseOrderUC.CancelPurchaseOrder - repo.Update: %w", err)
	}
	return po, nil
}

func (uc *PurchaseOrderUC) ListPurchaseOrders(ctx context.Context, filter domain.PurchaseOrderFilter, page pagination.Params) ([]domain.PurchaseOrder, error) {
	orders, err := uc.poRepo.List(ctx, filter, page)
	if err != nil {
		return nil, fmt.Errorf("PurchaseOrderUC.ListPurchaseOrders - repo.List: %w", err)
	}
	return orders, nil
}
