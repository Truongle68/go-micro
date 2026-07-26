package domain

import (
	"errors"
	"time"
)

type ProductVariant struct {
	ID           string
	VariantLabel string
	PriceDelta   float64
	Sku          string
}

type ProductImage struct {
	ID        string
	Url       string
	SortOrder int64
}

type Product struct {
	ID            string
	CategoryID    string
	Sku           string
	NameVi        string
	NameEn        string
	DescriptionVi string
	DescriptionEn string
	Unit          string
	BasePrice     float64
	SalePrice     float64
	RatingAvg     float64
	RatingCount   int64
	IsActive      bool
	Variants      []ProductVariant
	Images        []ProductImage
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type NewProductParams struct {
	CategoryID    string
	Sku           string
	NameVi        string
	NameEn        string
	DescriptionVi string
	DescriptionEn string
	Unit          string
	BasePrice     float64
	SalePrice     float64
	IsActive      bool
	Variants      []ProductVariant
	Images        []ProductImage
}

func NewProduct(params NewProductParams) (*Product, error) {
	if params.CategoryID == "" {
		return nil, ErrEmptyCategoryID
	}
	if params.NameVi == "" && params.NameEn == "" {
		return nil, ErrEmptyName
	}
	if params.Sku == "" {
		return nil, ErrEmptySku
	}
	if params.BasePrice < 0 || params.SalePrice < 0 {
		return nil, ErrInvalidPrice
	}

	now := time.Now()
	return &Product{
		CategoryID:    params.CategoryID,
		Sku:           params.Sku,
		NameVi:        params.NameVi,
		NameEn:        params.NameEn,
		DescriptionVi: params.DescriptionVi,
		DescriptionEn: params.DescriptionEn,
		Unit:          params.Unit,
		BasePrice:     params.BasePrice,
		SalePrice:     params.SalePrice,
		IsActive:      params.IsActive,
		Variants:      params.Variants,
		Images:        params.Images,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

type UpdateProductParams struct {
	CategoryID    *string
	Sku           *string
	NameVi        *string
	NameEn        *string
	DescriptionVi *string
	DescriptionEn *string
	Unit          *string
	BasePrice     *float64
	SalePrice     *float64
	IsActive      *bool
	Variants      []ProductVariant
	Images        []ProductImage
}

func (p *Product) ApplyUpdate(params UpdateProductParams) error {
	changed := false

	if params.CategoryID != nil {
		if *params.CategoryID == "" {
			return ErrEmptyCategoryID
		}
		p.CategoryID = *params.CategoryID
		changed = true
	}
	if params.Sku != nil {
		if *params.Sku == "" {
			return ErrEmptySku
		}
		p.Sku = *params.Sku
		changed = true
	}
	if params.NameVi != nil {
		p.NameVi = *params.NameVi
		changed = true
	}
	if params.NameEn != nil {
		p.NameEn = *params.NameEn
		changed = true
	}
	if p.NameVi == "" && p.NameEn == "" {
		return ErrEmptyName
	}
	if params.DescriptionVi != nil {
		p.DescriptionVi = *params.DescriptionVi
		changed = true
	}
	if params.DescriptionEn != nil {
		p.DescriptionEn = *params.DescriptionEn
		changed = true
	}
	if params.Unit != nil {
		p.Unit = *params.Unit
		changed = true
	}
	if params.BasePrice != nil {
		if *params.BasePrice < 0 {
			return ErrInvalidPrice
		}
		p.BasePrice = *params.BasePrice
		changed = true
	}
	if params.SalePrice != nil {
		if *params.SalePrice < 0 {
			return ErrInvalidPrice
		}
		p.SalePrice = *params.SalePrice
		changed = true
	}
	if params.IsActive != nil {
		p.IsActive = *params.IsActive
		changed = true
	}
	if params.Variants != nil {
		p.Variants = params.Variants
		changed = true
	}
	if params.Images != nil {
		p.Images = params.Images
		changed = true
	}

	if !changed {
		return ErrNoFieldsToUpdate
	}
	p.UpdatedAt = time.Now()
	return nil
}

type SearchProductParams struct {
	Query      string
	CategoryID string
	MinPrice   *float64
	MaxPrice   *float64
	IsActive   *bool
}

func (p SearchProductParams) Validate() error {
	if p.MinPrice != nil && p.MaxPrice != nil && *p.MinPrice > *p.MaxPrice {
		return ErrInvalidPriceRange
	}
	return nil
}

var ErrInvalidPriceRange = errors.New("min_price cannot be greater than max_price")

type ProductSearchResult struct {
	Products   []Product
	TotalCount int64
}
