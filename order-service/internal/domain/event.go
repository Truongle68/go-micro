package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	AggregateOrder   = "order"
	EventOrderPlaced = "order.placed"
)

var (
	ErrEmptyAggregateType = errors.New("aggregate type cannot be empty")
	ErrEmptyAggregateID   = errors.New("aggregate id cannot be empty")
	ErrEmptyEventType     = errors.New("event type cannot be empty")
	ErrNilEventPayload    = errors.New("event payload cannot be nil")
)

// OrderPlacedEvent represents a domain event emitted when an order is placed
type OrderPlacedEvent struct {
	OrderID       string            `json:"order_id"`
	UserID        string            `json:"user_id"`
	CustomerEmail string            `json:"customer_email"`
	Items         []OrderPlacedItem `json:"items"`
	TotalAmount   int64             `json:"total_amount_cents"`
	PlacedAt      time.Time         `json:"placed_at"`
}

type OrderPlacedItem struct {
	SKU            string `json:"sku"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int64  `json:"unit_price_cents"`
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
