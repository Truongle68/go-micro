package usecase

import (
	"context"
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
	outboxRepo     OutboxRepository
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
	outboxRepo OutboxRepository,
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
		outboxRepo:     outboxRepo,
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

type ReceiptLine struct {
	SKU      string
	Quantity int
}

func (uc *PurchaseOrderUC) ReceiveGoods(ctx context.Context, poID string, lines []domain.ReceiveLine) (*domain.PurchaseOrder, error) {
	var po *domain.PurchaseOrder

	err := uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		po, err = uc.poRepo.FindByID(txCtx, poID)
		if err != nil {
			return err
		}

		if err := po.ReceiveLines(lines); err != nil {
			return err
		}

		if err := uc.poRepo.Update(txCtx, po); err != nil {
			return fmt.Errorf("updating purchase order: %w", err)
		}

		// Atomically record domain event in outbox within the same transaction
		if uc.outboxRepo != nil {
			poLinesBySKU := make(map[string]domain.PurchaseOrderLine, len(po.Lines))
			for _, pl := range po.Lines {
				poLinesBySKU[pl.SKU] = pl
			}

			eventLines := make([]domain.GoodsReceivedLine, len(lines))
			for i, l := range lines {
				pl := poLinesBySKU[l.SKU]
				eventLines[i] = domain.GoodsReceivedLine{
					SKU:              l.SKU,
					ProductName:      pl.ProductName,
					QuantityOrdered:  pl.QuantityOrdered,
					QuantityReceived: l.Quantity,
					UnitCost:         pl.UnitCost,
				}
			}

			receivedAt := time.Now().UTC()
			if po.ReceivedAt != nil {
				receivedAt = *po.ReceivedAt
			}

			goodsReceivedEvent := domain.GoodsReceivedEvent{
				EventType:       domain.EventGoodsReceived,
				PurchaseOrderID: po.ID,
				POCode:          po.Code,
				SupplierID:      po.SupplierID,
				SupplierCode:    po.SupplierCode,
				WarehouseID:     po.WarehouseID,
				WarehouseCode:   po.WarehouseCode,
				ReceivedAt:      receivedAt,
				Lines:           eventLines,
			}

			outboxEvt, err := domain.NewOutboxEvent(
				domain.AggregatePurchaseOrder,
				po.ID,
				domain.EventGoodsReceived,
				goodsReceivedEvent,
			)
			if err != nil {
				return fmt.Errorf("creating outbox event: %w", err)
			}

			if err := uc.outboxRepo.Create(txCtx, outboxEvt); err != nil {
				return fmt.Errorf("saving outbox event: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("PurchaseOrderUC.ReceiveGoods: %w", err)
	}
	return po, nil
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
