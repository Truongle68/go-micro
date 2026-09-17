package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	EventGoodsReceived     = "goods.received"
	AggregatePurchaseOrder = "purchase_order"
)

var (
	ErrEmptyAggregateType = errors.New("aggregate type cannot be empty")
	ErrEmptyAggregateID   = errors.New("aggregate id cannot be empty")
	ErrEmptyEventType     = errors.New("event type cannot be empty")
	ErrNilEventPayload    = errors.New("event payload cannot be nil")
)

// GoodsReceivedEvent represents a domain event emitted when goods are received against a purchase order.
type GoodsReceivedEvent struct {
	EventType       string              `json:"event_type"`
	PurchaseOrderID string              `json:"purchase_order_id"`
	POCode          string              `json:"po_code"`
	SupplierID      string              `json:"supplier_id"`
	SupplierCode    string              `json:"supplier_code"`
	WarehouseID     string              `json:"warehouse_id"`
	WarehouseCode   string              `json:"warehouse_code"`
	ReceivedAt      time.Time           `json:"received_at"`
	Lines           []GoodsReceivedLine `json:"lines"`
}

// GoodsReceivedLine describes an individual SKU line received in a shipment.
type GoodsReceivedLine struct {
	SKU              string `json:"sku"`
	ProductName      string `json:"product_name"`
	QuantityOrdered  int    `json:"quantity_ordered"`
	QuantityReceived int    `json:"quantity_received"`
	UnitCost         int64  `json:"unit_cost"`
}

// OutboxEvent represents a transactional outbox table row.
type OutboxEvent struct {
	ID            string          `json:"id"`
	AggregateType string          `json:"aggregate_type"`
	AggregateID   string          `json:"aggregate_id"`
	EventType     string          `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
	CreatedAt     time.Time       `json:"created_at"`
	PublishedAt   *time.Time      `json:"published_at,omitempty"`
	RetryCount    int             `json:"retry_count"`
}

// NewOutboxEvent validates inputs and serializes the domain payload into an OutboxEvent.
func NewOutboxEvent(aggregateType, aggregateID, eventType string, payload any) (*OutboxEvent, error) {
	if strings.TrimSpace(aggregateType) == "" {
		return nil, ErrEmptyAggregateType
	}
	if strings.TrimSpace(aggregateID) == "" {
		return nil, ErrEmptyAggregateID
	}
	if strings.TrimSpace(eventType) == "" {
		return nil, ErrEmptyEventType
	}
	if payload == nil {
		return nil, ErrNilEventPayload
	}

	var rawPayload json.RawMessage
	switch v := payload.(type) {
	case json.RawMessage:
		rawPayload = v
	case []byte:
		rawPayload = json.RawMessage(v)
	default:
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("marshal outbox payload: %w", err)
		}
		rawPayload = json.RawMessage(data)
	}

	return &OutboxEvent{
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Payload:       rawPayload,
		CreatedAt:     time.Now().UTC(),
		RetryCount:    0,
	}, nil
}
