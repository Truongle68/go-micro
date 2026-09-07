package domain

import (
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
