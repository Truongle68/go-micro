package domain

import (
	"encoding/json"
	"strings"
	"time"
)

type Supplier struct {
	ID        string          `json:"id"`
	Version   int             `json:"version"`
	Code      string          `json:"code"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Phone     string          `json:"phone"`
	Address   SupplierAddress `json:"address"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type SupplierAddress struct {
	Line1    string `json:"line1"`
	Ward     string `json:"ward"`
	District string `json:"district"`
	City     string `json:"city"`
}

func NewSupplier(code, name, phone, email string, address SupplierAddress) (*Supplier, error) {
	if strings.TrimSpace(code) == "" {
		return nil, ErrEmptySuppCode
	}
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptySuppName
	}
	if strings.TrimSpace(phone) == "" {
		return nil, ErrEmptySuppPhone
	}
	if strings.TrimSpace(email) == "" {
		return nil, ErrEmptySuppEmail
	}
	now := time.Now().UTC()
	return &Supplier{
		Code: code, Version: 1, Name: name, Phone: phone, Email: email, Address: address,
		IsActive:  true,
		CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (s *Supplier) Deactivate() error {
	if !s.IsActive {
		return ErrSuppAlreadyInactive
	}
	s.IsActive = false
	s.Version++
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Supplier) Reactivate() error {
	if s.IsActive {
		return ErrSuppAlreadyActive
	}
	s.IsActive = true
	s.Version++
	s.UpdatedAt = time.Now().UTC()
	return nil
}

type Supply struct {
	ID          string    `json:"id"`
	SupplierID  string    `json:"supplier_id"`
	SKU         string    `json:"sku"`
	Price       int64     `json:"price"`
	LeadTime    LeadTime  `json:"lead_time"`
	MinOrderQty int       `json:"min_order_qty"`
	IsPreferred bool      `json:"is_preferred"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LeadTime struct {
	Days int `json:"days"`
}

func NewLeadTime(days int) (LeadTime, error) {
	if days < 0 {
		return LeadTime{}, ErrNegativeLeadDays
	}
	return LeadTime{Days: days}, nil
}

func (r SupplyAssignmentResult) MarshalJSON() ([]byte, error) {
	type Alias SupplyAssignmentResult
	var errStr *string
	if r.Err != nil {
		s := r.Err.Error()
		errStr = &s
	}

	return json.Marshal(&struct {
		Alias
		Err *string `json:"err"`
	}{
		Alias: (Alias)(r),
		Err:   errStr,
	})
}

func (l LeadTime) IsFast() bool {
	return l.Days <= 3
}

func (l LeadTime) AsDuration() time.Duration {
	return time.Duration(l.Days) * 24 * time.Hour
}

func NewSupply(supplierID, sku string, price int64, leadTime LeadTime, moq int) (*Supply, error) {
	if strings.TrimSpace(sku) == "" {
		return nil, ErrEmptySKU
	}

	return &Supply{
		SupplierID:  supplierID,
		SKU:         sku,
		Price:       price,
		LeadTime:    leadTime,
		MinOrderQty: moq,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s *Supply) MarkPreferred() {
	s.IsPreferred = true
}

func (s *Supply) UpdatePrice(price int64) error {
	if price < 0 {
		return ErrInvalidPrice
	}
	s.Price = price
	return nil
}

type SupplyInput struct {
	SKU         string
	Price       int64
	LeadTime    LeadTime
	MinOrderQty int
}

type SupplyAssignmentResult struct {
	SKU     string  `json:"sku"`
	Success bool    `json:"success"`
	Supply  *Supply `json:"supply"`
	Err     error   `json:"err"`
}
