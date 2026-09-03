package domain

import (
	"strings"
	"time"
)

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

type NewStockTransferParams struct {
	SKU             string
	Quantity        int
	FromWarehouseID string
	ToWarehouseID   string
}

func NewStockTransfer(params NewStockTransferParams) (*StockTransfer, error) {
	sku := strings.TrimSpace(params.SKU)
	fromWH := strings.TrimSpace(params.FromWarehouseID)
	toWH := strings.TrimSpace(params.ToWarehouseID)

	if sku == "" {
		return nil, ErrEmptySKU
	}
	if fromWH == "" {
		return nil, ErrEmptyFromWhID
	}
	if toWH == "" {
		return nil, ErrEmptyToWhID
	}
	if fromWH == toWH {
		return nil, ErrInvalidTransferRoute
	}

	if params.Quantity <= 0 {
		return nil, ErrNonPositiveQuantity
	}

	now := time.Now().UTC()
	return &StockTransfer{
		SKU:             sku,
		Quantity:        params.Quantity,
		FromWarehouseID: fromWH,
		ToWarehouseID:   toWH,
		Status:          TransferPending,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

type TransferStatus string

const (
	TransferPending   TransferStatus = "pending"
	TransferInTransit TransferStatus = "in_transit"
	TransferReceived  TransferStatus = "received"
	TransferCancelled TransferStatus = "cancelled"
)
