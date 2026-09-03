package domain

import (
	"strings"
	"time"
)

type StockReservation struct {
	ID          string            `json:"id"`
	SKU         string            `json:"sku"`
	WarehouseID string            `json:"warehouse_id"`
	Quantity    int               `json:"quantity"`
	OrderID     string            `json:"order_id"`
	Status      ReservationStatus `json:"status"`
	ExpiresAt   time.Time         `json:"expires_at"`
	CreatedAt   time.Time         `json:"created_at"`
	ConfirmedAt *time.Time        `json:"confirmed_at,omitempty"`
	ReleasedAt  *time.Time        `json:"released_at,omitempty"`
}

type NewStockReservationParams struct {
	SKU         string
	WarehouseID string
	Quantity    int
	OrderID     string
	TTL         time.Duration
}

func NewStockReservation(params NewStockReservationParams) (*StockReservation, error) {
	sku := strings.TrimSpace(params.SKU)
	wh := strings.TrimSpace(params.WarehouseID)

	if sku == "" {
		return nil, ErrEmptySKU
	}
	if wh == "" {
		return nil, ErrEmptyWhID
	}
	if params.Quantity <= 0 {
		return nil, ErrNonPositiveQuantity
	}

	now := time.Now().UTC()
	return &StockReservation{
		SKU:         sku,
		WarehouseID: wh,
		Quantity:    params.Quantity,
		OrderID:     params.OrderID,
		Status:      ReservationPending,
		ExpiresAt:   now.Add(params.TTL),
		CreatedAt:   now,
	}, nil
}

func (r *StockReservation) Confirm() error {
	if r.Status != ReservationPending {
		return ErrReservationNotPending
	}

	now := time.Now().UTC()
	r.Status = ReservationConfirmed
	r.ConfirmedAt = &now
	return nil
}

func (r *StockReservation) Release() error {
	if r.Status != ReservationPending {
		return ErrReservationNotPending
	}

	now := time.Now().UTC()
	r.Status = ReservationReleased
	r.ReleasedAt = &now
	return nil
}

func (r *StockReservation) Expire() error {
	if r.Status != ReservationPending {
		return ErrReservationNotPending
	}

	now := time.Now().UTC()
	r.Status = ReservationExpired
	r.ReleasedAt = &now
	return nil
}

func (r *StockReservation) IsExpired(now time.Time) bool {
	return now.After(r.ExpiresAt)
}

type ReservationStatus string

const (
	ReservationPending   ReservationStatus = "pending"
	ReservationConfirmed ReservationStatus = "confirmed"
	ReservationReleased  ReservationStatus = "released"
	ReservationExpired   ReservationStatus = "expired"
)
