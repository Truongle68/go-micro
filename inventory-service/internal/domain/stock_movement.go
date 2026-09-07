package domain

import "time"

type StockMovement struct {
	ID            string       `json:"id"`
	SKU           string       `json:"sku"`
	WarehouseID   string       `json:"warehouse_id"`
	Type          MovementType `json:"type"`
	Quantity      int          `json:"quantity"`
	ReferenceType string       `json:"reference_type"`
	ReferenceID   string       `json:"reference_id"`
	Note          string       `json:"note"`
	CreatedBy     string       `json:"created_by"`
	CreatedAt     time.Time    `json:"created_at"`
}

type MovementType string

const (
	MovementInbound     MovementType = "inbound"
	MovementOutbound    MovementType = "outbound"
	MovementReserved    MovementType = "reserved"
	MovementReleased    MovementType = "released"
	MovementAdjustment  MovementType = "adjustment"
	MovementTransferIn  MovementType = "transfer_in"
	MovementTransferOut MovementType = "transfer_out"
)

type StockMovementFilter struct {
	SKU         string
	WarehouseID string
	Type        MovementType
}
