package domain

import (
	"strings"
	"time"
)

type StockLevel struct {
	ID               string    `json:"id"`
	SKU              string    `json:"sku"`
	WarehouseID      string    `json:"warehouse_id"`
	OnHand           int       `json:"on_hand"`
	Reserved         int       `json:"reserved"`
	ReorderThreshold int       `json:"reorder_threshold"`
	ReorderQuantity  int       `json:"reorder_quantity"`
	Version          int       `json:"version"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type NewStockLevelParams struct {
	SKU              string
	WarehouseID      string
	ReorderThreshold int
	ReorderQuantity  int
}

func NewStockLevel(params NewStockLevelParams) (*StockLevel, error) {
	if strings.TrimSpace(params.SKU) == "" {
		return nil, ErrEmptySKU
	}
	if strings.TrimSpace(params.WarehouseID) == "" {
		return nil, ErrEmptyWhID
	}

	if params.ReorderThreshold < 0 || params.ReorderQuantity < 0 {
		return nil, ErrNegativeQuantity
	}

	now := time.Now().UTC()
	return &StockLevel{
		SKU:              params.SKU,
		WarehouseID:      params.WarehouseID,
		OnHand:           0,
		Reserved:         0,
		ReorderThreshold: params.ReorderThreshold,
		ReorderQuantity:  params.ReorderQuantity,
		Version:          1,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

func (s *StockLevel) Available() int {
	avai := s.OnHand - s.Reserved
	if avai < 0 {
		return 0
	}
	return avai
}

func (s *StockLevel) NeedsReorder() bool {
	return s.Available() <= s.ReorderThreshold
}

func (s *StockLevel) Reserve(qty int) error {
	if qty <= 0 {
		return ErrNonPositiveQuantity
	}

	if s.Available() < qty {
		return ErrInsufficientStock
	}

	s.Reserved += qty
	s.Version++
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *StockLevel) ConfirmReservation(qty int) error {
	if qty <= 0 {
		return ErrNonPositiveQuantity
	}

	if s.Reserved < qty {
		return ErrInsufficientStock
	}

	s.Reserved -= qty
	s.OnHand -= qty
	s.Version++
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *StockLevel) Release(qty int) error {
	if qty <= 0 {
		return ErrNonPositiveQuantity
	}

	if s.Reserved < qty {
		return ErrInsufficientStock
	}

	s.Reserved -= qty
	s.Version++
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *StockLevel) AdjustOnHand(delta int) error {
	newOnHand := s.OnHand + delta
	if newOnHand < 0 {
		return ErrInsufficientStock
	}
	s.OnHand = newOnHand
	s.Version++
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *StockLevel) ApplyTransferOut(qty int) error {
	if qty <= 0 {
		return ErrNonPositiveQuantity
	}
	if s.Available() < qty {
		return ErrInsufficientStock
	}
	s.OnHand -= qty
	s.Version++
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *StockLevel) ApplyTransferIn(qty int) {
	s.OnHand += qty
	s.Version++
	s.UpdatedAt = time.Now().UTC()
}

type StockLevelFilter struct {
	WarehouseID string
	SKU         string
	LowStock    bool
}

type StockSummary struct {
	TotalOnHand     int `json:"total_on_hand"`
	TotalReserved   int `json:"total_reserved"`
	TotalAvailable  int `json:"total_available"`
	TotalSKUs       int `json:"total_skus"`
	LowStockCount   int `json:"low_stock_count"`
	OutOfStockCount int `json:"out_of_stock_count"`
}
