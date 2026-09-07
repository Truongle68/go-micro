package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"inventory-service/internal/client"
	"inventory-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/TruongLe68/go-micro/pkg/rabbitmq/publisher"
)

const _defaultReservationTTL = 15 * time.Minute

type SKUQty struct {
	SKU      string
	Quantity int
}

type DetailedStockLevel struct {
	domain.StockLevel
	WarehouseCode string `json:"warehouse_code"`
	WarehouseName string `json:"warehouse_name"`
	ProductName   string `json:"product_name"`
	ProductImage  string `json:"product_image,omitempty"`
	Available     int    `json:"available"`
	LowStock      bool   `json:"low_stock"`
}

type DetailedStockMovement struct {
	domain.StockMovement
	WarehouseCode string `json:"warehouse_code"`
	WarehouseName string `json:"warehouse_name"`
	ProductName   string `json:"product_name"`
}

type AdjustStockInput struct {
	SKU           string
	WarehouseID   string
	QuantityDelta int
	Reason        string
	Note          string
	CreatedBy     string
}

// StockUC handles stock checking, reservation, confirmation, release, and management.
type StockUC struct {
	stockRepo     StockLevelRepository
	warehouseRepo WarehouseRepository
	resRepo       StockReservationRepository
	movRepo       StockMovementRepository
	catalogClient CatalogClient
	transactor    Transactor
	publisher     EventPublisher
	logger        logger.Interface
}

// NewStockUC creates a new StockUC.
func NewStockUC(
	stockRepo StockLevelRepository,
	warehouseRepo WarehouseRepository,
	resRepo StockReservationRepository,
	movRepo StockMovementRepository,
	catalogClient CatalogClient,
	transactor Transactor,
	pub EventPublisher,
	l logger.Interface,
) *StockUC {
	return &StockUC{
		stockRepo:     stockRepo,
		warehouseRepo: warehouseRepo,
		resRepo:       resRepo,
		movRepo:       movRepo,
		catalogClient: catalogClient,
		transactor:    transactor,
		publisher:     pub,
		logger:        l,
	}
}

// CheckStock returns the available quantity for each requested SKU,
// summed across all warehouses.
func (uc *StockUC) CheckStock(ctx context.Context, skus []string) (map[string]int, error) {
	if len(skus) == 0 {
		return map[string]int{}, nil
	}

	avail, err := uc.stockRepo.BulkCheckAvailability(ctx, skus)
	if err != nil {
		return nil, fmt.Errorf("StockUC.CheckStock: %w", err)
	}

	return avail, nil
}

// ReserveStock reserves inventory for an order. For each item, it finds the first
// warehouse with sufficient stock, then atomically reserves the quantity.
func (uc *StockUC) ReserveStock(ctx context.Context, orderID string, items []SKUQty) error {
	if len(items) == 0 {
		return nil
	}

	var eventItems []publisher.SKUQtyItem

	err := uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		eventItems = make([]publisher.SKUQtyItem, 0, len(items))

		for _, item := range items {
			// Find the first warehouse with enough available stock
			sl, err := uc.stockRepo.FindFirstAvailable(txCtx, item.SKU, item.Quantity)
			if err != nil {
				return fmt.Errorf("ReserveStock - FindFirstAvailable sku=%s: %w", item.SKU, err)
			}

			// Reserve on the domain entity
			if err := sl.Reserve(item.Quantity); err != nil {
				return fmt.Errorf("ReserveStock - Reserve sku=%s: %w", item.SKU, err)
			}

			// Persist the updated stock level (optimistic lock)
			if err := uc.stockRepo.Update(txCtx, sl); err != nil {
				return fmt.Errorf("ReserveStock - Update stock_level sku=%s: %w", item.SKU, err)
			}

			// Create a reservation record
			reservation, err := domain.NewStockReservation(domain.NewStockReservationParams{
				SKU:         item.SKU,
				WarehouseID: sl.WarehouseID,
				Quantity:    item.Quantity,
				OrderID:     orderID,
				TTL:         _defaultReservationTTL,
			})
			if err != nil {
				return fmt.Errorf("ReserveStock - NewStockReservation sku=%s: %w", item.SKU, err)
			}
			if err := uc.resRepo.Create(txCtx, reservation); err != nil {
				return fmt.Errorf("ReserveStock - Create reservation sku=%s: %w", item.SKU, err)
			}

			// Record movement
			movement := &domain.StockMovement{
				SKU:           item.SKU,
				WarehouseID:   sl.WarehouseID,
				Type:          domain.MovementReserved,
				Quantity:      item.Quantity,
				ReferenceType: "order",
				ReferenceID:   orderID,
				Note:          fmt.Sprintf("Reserved %d units for order %s", item.Quantity, orderID),
				CreatedAt:     time.Now().UTC(),
			}
			if err := uc.movRepo.Create(txCtx, movement); err != nil {
				return fmt.Errorf("ReserveStock - Create movement sku=%s: %w", item.SKU, err)
			}

			eventItems = append(eventItems, publisher.SKUQtyItem{
				SKU:         item.SKU,
				WarehouseID: sl.WarehouseID,
				Quantity:    item.Quantity,
			})
		}

		return nil
	})
	if err != nil {
		return err
	}

	uc.publishEvent(ctx, publisher.EventStockReserved, publisher.StockReservedPayload{
		OrderID: orderID,
		Items:   eventItems,
	})

	return nil
}

// ConfirmReservation confirms all pending reservations for an order,
// deducting on-hand stock and marking reservations as confirmed.
func (uc *StockUC) ConfirmReservation(ctx context.Context, orderID string) error {
	var eventItems []publisher.SKUQtyItem

	err := uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		reservations, err := uc.resRepo.FindPendingByOrderID(txCtx, orderID)
		if err != nil {
			return fmt.Errorf("ConfirmReservation - FindPendingByOrderID: %w", err)
		}

		eventItems = make([]publisher.SKUQtyItem, 0, len(reservations))

		for i := range reservations {
			res := &reservations[i]

			// Confirm on stock level (decrements reserved AND on_hand)
			sl, err := uc.stockRepo.FindBySKUAndWarehouse(txCtx, res.SKU, res.WarehouseID)
			if err != nil {
				return fmt.Errorf("ConfirmReservation - FindBySKUAndWarehouse sku=%s: %w", res.SKU, err)
			}

			if err := sl.ConfirmReservation(res.Quantity); err != nil {
				return fmt.Errorf("ConfirmReservation - ConfirmReservation sku=%s: %w", res.SKU, err)
			}

			if err := uc.stockRepo.Update(txCtx, sl); err != nil {
				return fmt.Errorf("ConfirmReservation - Update stock_level sku=%s: %w", res.SKU, err)
			}

			// Mark reservation as confirmed
			if err := res.Confirm(); err != nil {
				return fmt.Errorf("ConfirmReservation - Confirm reservation sku=%s: %w", res.SKU, err)
			}
			if err := uc.resRepo.Update(txCtx, res); err != nil {
				return fmt.Errorf("ConfirmReservation - Update reservation sku=%s: %w", res.SKU, err)
			}

			// Record movement
			movement := &domain.StockMovement{
				SKU:           res.SKU,
				WarehouseID:   res.WarehouseID,
				Type:          domain.MovementOutbound,
				Quantity:      res.Quantity,
				ReferenceType: "order",
				ReferenceID:   orderID,
				Note:          fmt.Sprintf("Confirmed %d units deducted for order %s", res.Quantity, orderID),
				CreatedAt:     time.Now().UTC(),
			}
			if err := uc.movRepo.Create(txCtx, movement); err != nil {
				return fmt.Errorf("ConfirmReservation - Create movement sku=%s: %w", res.SKU, err)
			}

			eventItems = append(eventItems, publisher.SKUQtyItem{
				SKU:         res.SKU,
				WarehouseID: res.WarehouseID,
				Quantity:    res.Quantity,
			})
		}

		return nil
	})
	if err != nil {
		return err
	}

	uc.publishEvent(ctx, publisher.EventStockConfirmed, publisher.StockConfirmedPayload{
		OrderID: orderID,
		Items:   eventItems,
	})

	return nil
}

// ReleaseReservation releases all pending reservations for an order,
// returning reserved stock to available.
func (uc *StockUC) ReleaseReservation(ctx context.Context, orderID string) error {
	var eventItems []publisher.SKUQtyItem

	err := uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		reservations, err := uc.resRepo.FindPendingByOrderID(txCtx, orderID)
		if err != nil {
			return fmt.Errorf("ReleaseReservation - FindPendingByOrderID: %w", err)
		}

		eventItems = make([]publisher.SKUQtyItem, 0, len(reservations))

		for i := range reservations {
			res := &reservations[i]

			// Release on stock level (decrements reserved)
			sl, err := uc.stockRepo.FindBySKUAndWarehouse(txCtx, res.SKU, res.WarehouseID)
			if err != nil {
				return fmt.Errorf("ReleaseReservation - FindBySKUAndWarehouse sku=%s: %w", res.SKU, err)
			}

			if err := sl.Release(res.Quantity); err != nil {
				return fmt.Errorf("ReleaseReservation - Release sku=%s: %w", res.SKU, err)
			}

			if err := uc.stockRepo.Update(txCtx, sl); err != nil {
				return fmt.Errorf("ReleaseReservation - Update stock_level sku=%s: %w", res.SKU, err)
			}

			// Mark reservation as released
			if err := res.Release(); err != nil {
				return fmt.Errorf("ReleaseReservation - Release reservation sku=%s: %w", res.SKU, err)
			}
			if err := uc.resRepo.Update(txCtx, res); err != nil {
				return fmt.Errorf("ReleaseReservation - Update reservation sku=%s: %w", res.SKU, err)
			}

			// Record movement
			movement := &domain.StockMovement{
				SKU:           res.SKU,
				WarehouseID:   res.WarehouseID,
				Type:          domain.MovementReleased,
				Quantity:      res.Quantity,
				ReferenceType: "order",
				ReferenceID:   orderID,
				Note:          fmt.Sprintf("Released %d reserved units for order %s", res.Quantity, orderID),
				CreatedAt:     time.Now().UTC(),
			}
			if err := uc.movRepo.Create(txCtx, movement); err != nil {
				return fmt.Errorf("ReleaseReservation - Create movement sku=%s: %w", res.SKU, err)
			}

			eventItems = append(eventItems, publisher.SKUQtyItem{
				SKU:         res.SKU,
				WarehouseID: res.WarehouseID,
				Quantity:    res.Quantity,
			})
		}

		return nil
	})
	if err != nil {
		return err
	}

	uc.publishEvent(ctx, publisher.EventStockReleased, publisher.StockReleasedPayload{
		OrderID: orderID,
		Items:   eventItems,
	})

	return nil
}

// GetStockLevel returns stock levels for a SKU across all warehouses.
// GetSKUAvailability returns the total available quantity for a single SKU.
func (uc *StockUC) GetSKUAvailability(ctx context.Context, sku string) (int, error) {
	if sku == "" {
		return 0, domain.ErrEmptySKU
	}
	avail, err := uc.stockRepo.BulkCheckAvailability(ctx, []string{sku})
	if err != nil {
		return 0, fmt.Errorf("StockUC.GetSKUAvailability: %w", err)
	}
	return avail[sku], nil
}

// GetStockLevelsBySKU returns all raw warehouse stock levels for a given SKU.
func (uc *StockUC) GetStockLevelsBySKU(ctx context.Context, sku string) ([]domain.StockLevel, error) {
	levels, err := uc.stockRepo.FindBySKU(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("StockUC.GetStockLevelsBySKU: %w", err)
	}
	return levels, nil
}

// ListStockLevels returns a paginated list of detailed stock levels with joined warehouse and catalog data.
func (uc *StockUC) ListStockLevels(ctx context.Context, filter domain.StockLevelFilter, p pagination.Params) ([]DetailedStockLevel, int64, error) {
	levels, total, err := uc.stockRepo.List(ctx, filter, p)
	if err != nil {
		return nil, 0, fmt.Errorf("StockUC.ListStockLevels - repo.List: %w", err)
	}

	if len(levels) == 0 {
		return []DetailedStockLevel{}, total, nil
	}

	// In-memory join: fetch warehouses
	whMap := make(map[string]domain.Warehouse)
	if uc.warehouseRepo != nil {
		warehouses, err := uc.warehouseRepo.Find(ctx)
		if err == nil {
			for _, w := range warehouses {
				whMap[w.ID] = w
			}
		}
	}

	// In-memory join: fetch product variants from catalog-service
	skuSet := make(map[string]struct{}, len(levels))
	for _, l := range levels {
		skuSet[l.SKU] = struct{}{}
	}
	skus := make([]string, 0, len(skuSet))
	for s := range skuSet {
		skus = append(skus, s)
	}

	variantMap := make(map[string]client.VariantDTO)
	if uc.catalogClient != nil && len(skus) > 0 {
		variants, err := uc.catalogClient.GetVariantsBySKUs(ctx, skus)
		if err == nil {
			for _, v := range variants {
				variantMap[v.SKU] = v
			}
		}
	}

	result := make([]DetailedStockLevel, len(levels))
	for i, l := range levels {
		wh := whMap[l.WarehouseID]
		variant := variantMap[l.SKU]

		result[i] = DetailedStockLevel{
			StockLevel:    l,
			WarehouseCode: wh.Code,
			WarehouseName: wh.Name,
			ProductName:   variant.ProductName,
			ProductImage:  variant.Image,
			Available:     l.Available(),
			LowStock:      l.NeedsReorder(),
		}
	}

	return result, total, nil
}

// GetStockLevel retrieves a single stock level by ID, enriched with warehouse and catalog details.
func (uc *StockUC) GetStockLevel(ctx context.Context, id string) (*DetailedStockLevel, error) {
	l, err := uc.stockRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("StockUC.GetStockLevel - find: %w", err)
	}
	if l == nil {
		return nil, nil
	}

	var whCode, whName string
	if uc.warehouseRepo != nil {
		wh, err := uc.warehouseRepo.FindByID(ctx, l.WarehouseID)
		if err == nil && wh != nil {
			whCode = wh.Code
			whName = wh.Name
		}
	}

	var prodName, prodImage string
	if uc.catalogClient != nil {
		variants, err := uc.catalogClient.GetVariantsBySKUs(ctx, []string{l.SKU})
		if err == nil && len(variants) > 0 {
			prodName = variants[0].ProductName
			prodImage = variants[0].Image
		}
	}

	return &DetailedStockLevel{
		StockLevel:    *l,
		WarehouseCode: whCode,
		WarehouseName: whName,
		ProductName:   prodName,
		ProductImage:  prodImage,
		Available:     l.Available(),
		LowStock:      l.NeedsReorder(),
	}, nil
}

// AdjustStock applies a manual quantity adjustment to a stock level and records an audit movement.
func (uc *StockUC) AdjustStock(ctx context.Context, input AdjustStockInput) (*DetailedStockLevel, error) {
	if input.SKU == "" {
		return nil, domain.ErrEmptySKU
	}
	if input.WarehouseID == "" {
		return nil, domain.ErrEmptyWhID
	}
	if input.QuantityDelta == 0 {
		return nil, fmt.Errorf("%w: delta cannot be zero", domain.ErrNonPositiveQuantity)
	}

	var updatedLevelID string

	err := uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		sl, err := uc.stockRepo.FindBySKUAndWarehouse(txCtx, input.SKU, input.WarehouseID)
		if err != nil {
			return fmt.Errorf("AdjustStock - find stock level: %w", err)
		}

		if sl == nil {
			if input.QuantityDelta < 0 {
				return fmt.Errorf("cannot reduce stock on nonexistent stock level: %w", domain.ErrInsufficientStock)
			}
			newSL, err := domain.NewStockLevel(domain.NewStockLevelParams{
				SKU:         input.SKU,
				WarehouseID: input.WarehouseID,
			})
			if err != nil {
				return err
			}
			newSL.OnHand = input.QuantityDelta
			if err := uc.stockRepo.Create(txCtx, newSL); err != nil {
				return fmt.Errorf("AdjustStock - create stock level: %w", err)
			}
			sl = newSL
		} else {
			if err := sl.AdjustOnHand(input.QuantityDelta); err != nil {
				return err
			}
			if err := uc.stockRepo.Update(txCtx, sl); err != nil {
				return fmt.Errorf("AdjustStock - update stock level: %w", err)
			}
		}

		refType := input.Reason
		if refType == "" {
			refType = "manual_adjustment"
		}
		mov := &domain.StockMovement{
			SKU:           input.SKU,
			WarehouseID:   input.WarehouseID,
			Type:          domain.MovementAdjustment,
			Quantity:      input.QuantityDelta,
			ReferenceType: refType,
			Note:          input.Note,
			CreatedBy:     input.CreatedBy,
			CreatedAt:     time.Now().UTC(),
		}
		if err := uc.movRepo.Create(txCtx, mov); err != nil {
			return fmt.Errorf("AdjustStock - create movement: %w", err)
		}

		updatedLevelID = sl.ID
		return nil
	})

	if err != nil {
		return nil, err
	}

	return uc.GetStockLevel(ctx, updatedLevelID)
}

// UpdateThresholds updates the reorder threshold and reorder quantity for a stock level.
func (uc *StockUC) UpdateThresholds(ctx context.Context, id string, reorderThreshold, reorderQuantity int) (*DetailedStockLevel, error) {
	if reorderThreshold < 0 || reorderQuantity < 0 {
		return nil, domain.ErrNegativeQuantity
	}

	sl, err := uc.stockRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("StockUC.UpdateThresholds - find: %w", err)
	}
	if sl == nil {
		return nil, fmt.Errorf("%w: stock level id=%s", domain.ErrStockLevelNotFound, id)
	}

	sl.ReorderThreshold = reorderThreshold
	sl.ReorderQuantity = reorderQuantity
	sl.Version++
	sl.UpdatedAt = time.Now().UTC()

	if err := uc.stockRepo.Update(ctx, sl); err != nil {
		return nil, fmt.Errorf("StockUC.UpdateThresholds - update: %w", err)
	}

	return uc.GetStockLevel(ctx, id)
}

// GetStockSummary returns inventory KPI metrics across all or a specific warehouse.
func (uc *StockUC) GetStockSummary(ctx context.Context, warehouseID string) (*domain.StockSummary, error) {
	summary, err := uc.stockRepo.GetSummary(ctx, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("StockUC.GetStockSummary: %w", err)
	}
	return summary, nil
}

// ListMovements returns paginated stock movements audit logs enriched with warehouse and product names.
func (uc *StockUC) ListMovements(ctx context.Context, filter domain.StockMovementFilter, p pagination.Params) ([]DetailedStockMovement, int64, error) {
	movs, total, err := uc.movRepo.List(ctx, filter, p)
	if err != nil {
		return nil, 0, fmt.Errorf("StockUC.ListMovements - repo.List: %w", err)
	}

	if len(movs) == 0 {
		return []DetailedStockMovement{}, total, nil
	}

	whMap := make(map[string]domain.Warehouse)
	if uc.warehouseRepo != nil {
		warehouses, err := uc.warehouseRepo.Find(ctx)
		if err == nil {
			for _, w := range warehouses {
				whMap[w.ID] = w
			}
		}
	}

	skuSet := make(map[string]struct{}, len(movs))
	for _, m := range movs {
		skuSet[m.SKU] = struct{}{}
	}
	skus := make([]string, 0, len(skuSet))
	for s := range skuSet {
		skus = append(skus, s)
	}

	variantMap := make(map[string]client.VariantDTO)
	if uc.catalogClient != nil && len(skus) > 0 {
		variants, err := uc.catalogClient.GetVariantsBySKUs(ctx, skus)
		if err == nil {
			for _, v := range variants {
				variantMap[v.SKU] = v
			}
		}
	}

	result := make([]DetailedStockMovement, len(movs))
	for i, m := range movs {
		wh := whMap[m.WarehouseID]
		variant := variantMap[m.SKU]

		result[i] = DetailedStockMovement{
			StockMovement: m,
			WarehouseCode: wh.Code,
			WarehouseName: wh.Name,
			ProductName:   variant.ProductName,
		}
	}

	return result, total, nil
}

// publishEvent is a best-effort helper that publishes domain events without
// blocking the caller if the publish fails.
func (uc *StockUC) publishEvent(ctx context.Context, eventType string, payload any) {
	if uc.publisher == nil {
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		uc.logger.Warn("StockUC.publishEvent - marshal %s: %v", eventType, err)
		return
	}

	event := publisher.Event{
		Type:      eventType,
		Payload:   data,
		Timestamp: time.Now().UTC(),
	}

	if err := uc.publisher.Publish(ctx, event); err != nil {
		uc.logger.Warn("StockUC.publishEvent - publish %s: %v", eventType, err)
	}
}
