package domain

import (
	"strings"
	"time"
)

type Supplier struct {
	ID        string
	Version   int
	Code      string
	Name      string
	Email     string
	Phone     string
	Address   SupplierAddress
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SupplierAddress struct {
	Line1    string
	Ward     string
	District string
	City     string
}

func NewSupplier(code, name, phone string, address SupplierAddress) (*Supplier, error) {
	if strings.TrimSpace(code) == "" {
		return nil, ErrEmptySuppCode
	}
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptySuppName
	}
	now := time.Now().UTC()
	return &Supplier{
		Code: code, Version: 1, Name: name, Phone: phone, Address: address,
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
