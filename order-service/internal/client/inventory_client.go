package client

import "context"

// SKUQty represents a SKU and its quantity for stock operations.
type SKUQty struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

// InventoryClient defines the interface for communicating with the inventory service.
type InventoryClient interface {
	// CheckStock returns available quantities for the given items.
	CheckStock(ctx context.Context, items []SKUQty) (map[string]int, error)

	// ReserveStock reserves inventory for an order.
	ReserveStock(ctx context.Context, orderID string, items []SKUQty) error

	// ConfirmReservation confirms the reservation, deducting on-hand stock.
	ConfirmReservation(ctx context.Context, orderID string) error

	// ReleaseReservation releases reserved stock back to available.
	ReleaseReservation(ctx context.Context, orderID string) error
}
