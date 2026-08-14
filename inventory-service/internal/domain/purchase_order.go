package domain

import "time"

type PurchaseOrder struct {
	ID          string              `json:"id"`
	SupplierID  string              `json:"supplier_id"`
	WarehouseID string              `json:"warehouse_id"`
	Lines       []PurchaseOrderLine `json:"lines"`
	Status      PurchaseOrderStatus `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	ExpectedAt  *time.Time          `json:"expected_at,omitempty"`
	ReceivedAt  *time.Time          `json:"received_at,omitempty"`
}

type PurchaseOrderLine struct {
	SKU              string `json:"sku"`
	QuantityOrdered  int    `json:"quantity_ordered"`
	QuantityReceived int    `json:"quantity_received"`
	UnitCost         int64  `json:"unit_cost"`
}

type PurchaseOrderStatus string

const (
	POStatusDraft             PurchaseOrderStatus = "draft"
	POStatusOrdered           PurchaseOrderStatus = "ordered"
	POStatusPartiallyReceived PurchaseOrderStatus = "partially_received"
	POStatusReceived          PurchaseOrderStatus = "received"
	POStatusCancelled         PurchaseOrderStatus = "cancelled"
)
