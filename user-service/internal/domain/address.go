package domain

import "time"

type AddressLabel string

const (
	AddressLabelHome AddressLabel = "home"
	AddressLabelWork AddressLabel = "work"
)

type Address struct {
	ID          string       `json:"id"`
	UserID      string       `json:"user_id"`
	Label       AddressLabel `json:"label"`
	AddressLine string       `json:"address_line"`
	Ward        string       `json:"ward"`
	District    string       `json:"district"`
	City        string       `json:"city"`
	Lat         float64      `json:"lat"`
	Lng         float64      `json:"lng"`
	IsDefault   bool         `json:"is_default"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func NewAddress(userID, line, city string, label AddressLabel, opts ...AddressOption) (*Address, error) {
	if userID == "" {
		return nil, ErrEmptyUserID
	}
	if line == "" {
		return nil, ErrEmptyAddressLine
	}
	if city == "" {
		return nil, ErrEmptyCity
	}

	now := time.Now()
	a := &Address{
		UserID:      userID,
		Label:       label,
		AddressLine: line,
		City:        city,
		IsDefault:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a, nil
}

type UpdateAddressParams struct {
	Label       *AddressLabel
	AddressLine *string
	Ward        *string
	District    *string
	City        *string
	Lat         *float64
	Lng         *float64
}

func (a *Address) ApplyUpdate(params UpdateAddressParams) error {
	change := false
	if params.Label != nil {
		a.Label = *params.Label
		change = true
	}
	if params.AddressLine != nil {
		if *params.AddressLine == "" {
			return ErrEmptyAddressLine
		}
		a.AddressLine = *params.AddressLine
		change = true
	}
	if params.Ward != nil {
		a.Ward = *params.Ward
		change = true
	}
	if params.District != nil {
		a.District = *params.District
		change = true
	}
	if params.City != nil {
		if *params.City == "" {
			return ErrEmptyCity
		}
		a.City = *params.City
		change = true
	}
	if params.Lat != nil {
		a.Lat = *params.Lat
		change = true
	}
	if params.Lng != nil {
		a.Lng = *params.Lng
		change = true
	}

	if !change {
		return ErrNoFieldsToUpdate
	}
	a.UpdatedAt = time.Now()
	return nil
}
