package domain

import "time"

type StockTransfer struct {
	ID              string         `json:"id"`
	SKU             string         `json:"sku"`
	Quantity        int            `json:"quantity"`
	FromWarehouseID string         `json:"from_warehouse_id"`
	ToWarehouseID   string         `json:"to_warehouse_id"`
	Status          TransferStatus `json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	ShippedAt       *time.Time     `json:"shipped_at,omitempty"`
	ReceivedAt      *time.Time     `json:"received_at,omitempty"`
}

type TransferStatus string

const (
	TransferPending   TransferStatus = "pending"
	TransferInTransit TransferStatus = "in_transit"
	TransferReceived  TransferStatus = "received"
	TransferCancelled TransferStatus = "cancelled"
)
