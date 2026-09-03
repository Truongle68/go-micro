package domain

import "time"

type Warehouse struct {
	ID        string           `json:"id"`
	Code      string           `json:"code"`
	Name      string           `json:"name"`
	Region    string           `json:"region"`
	Address   WarehouseAddress `json:"address"`
	IsActive  bool             `json:"is_active"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type WarehouseAddress struct {
	Line1    string  `json:"line1"`
	Ward     string  `json:"ward,omitempty"`
	District string  `json:"district,omitempty"`
	City     string  `json:"city"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

type NewWarehouseParams struct {
	Code    string
	Name    string
	Region  string
	Address WarehouseAddress
}

func NewWarehouse(params NewWarehouseParams) (*Warehouse, error) {
	if params.Code == "" {
		return nil, ErrEmptyWhCode
	}

	if params.Name == "" {
		return nil, ErrEmptyWhName
	}

	if params.Address.City == "" {
		return nil, ErrEmptyWhAddressCity
	}

	now := time.Now().UTC()
	return &Warehouse{
		Code:      params.Code,
		Name:      params.Name,
		Region:    params.Region,
		Address:   params.Address,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
