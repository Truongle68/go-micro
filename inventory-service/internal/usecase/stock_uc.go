package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"inventory-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/rabbitmq/publisher"
)

const _defaultReservationTTL = 15 * time.Minute

type SKUQty struct {
	SKU      string
	Quantity int
}

// StockUC handles stock checking, reservation, confirmation, and release.
type StockUC struct {
	stockRepo   StockLevelRepository
	resRepo     StockReservationRepository
	movRepo     StockMovementRepository
	transactor  Transactor
	publisher   EventPublisher
	logger      logger.Interface
}

// NewStockUC creates a new StockUC.
func NewStockUC(
	stockRepo StockLevelRepository,
	resRepo StockReservationRepository,
	movRepo StockMovementRepository,
	transactor Transactor,
	pub EventPublisher,
	l logger.Interface,
) *StockUC {
	return &StockUC{
		stockRepo:  stockRepo,
		resRepo:    resRepo,
		movRepo:    movRepo,
		transactor: transactor,
		publisher:  pub,
		logger:     l,
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
// warehouse with sufficient stock (Option A), then atomically reserves the quantity.
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
func (uc *StockUC) GetStockLevel(ctx context.Context, sku string) ([]domain.StockLevel, error) {
	levels, err := uc.stockRepo.FindBySKU(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("StockUC.GetStockLevel: %w", err)
	}
	return levels, nil
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
