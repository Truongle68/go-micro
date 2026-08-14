package domain

import "time"

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
		return ErrInvalidQuantity
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
		return ErrInvalidQuantity
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
		return ErrInvalidQuantity
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
		return ErrInvalidQuantity
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
