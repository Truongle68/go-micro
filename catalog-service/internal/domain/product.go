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

type PopulatedCategory struct {
	ID        string
	NameVi    string
	NameEn    string
	Slug      string
	Icon      string
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
	if params.CategoryID != nil {
		if *params.CategoryID == "" {
			return ErrEmptyCategoryID
		}
		p.CategoryID = *params.CategoryID
	}
	if params.Sku != nil {
		if *params.Sku == "" {
			return ErrEmptySku
		}
		p.Sku = *params.Sku
	}
	if params.NameVi != nil {
		p.NameVi = *params.NameVi
	}
	if params.NameEn != nil {
		p.NameEn = *params.NameEn
	}
	if p.NameVi == "" && p.NameEn == "" {
		return ErrEmptyName
	}
	if params.DescriptionVi != nil {
		p.DescriptionVi = *params.DescriptionVi
	}
	if params.DescriptionEn != nil {
		p.DescriptionEn = *params.DescriptionEn
	}
	if params.Unit != nil {
		p.Unit = *params.Unit
	}
	if params.BasePrice != nil {
		if *params.BasePrice < 0 {
			return ErrInvalidPrice
		}
		p.BasePrice = *params.BasePrice
	}
	if params.SalePrice != nil {
		if *params.SalePrice < 0 {
			return ErrInvalidPrice
		}
		p.SalePrice = *params.SalePrice
	}
	if params.IsActive != nil {
		p.IsActive = *params.IsActive
	}
	if params.Variants != nil {
		p.Variants = params.Variants
	}
	if params.Images != nil {
		p.Images = params.Images
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

type SearchProductsQuery struct {
	Query      string   `form:"q"`
	CategoryID string   `form:"category_id"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	IsActive   *bool    `form:"is_active"`
}

func (q SearchProductsQuery) ToDomainParams() SearchProductParams {
	return SearchProductParams{
		Query:      q.Query,
		CategoryID: q.CategoryID,
		MinPrice:   q.MinPrice,
		MaxPrice:   q.MaxPrice,
		IsActive:   q.IsActive,
	}
}

func (p SearchProductParams) Validate() error {
	if (p.MinPrice != nil && *p.MinPrice < 0) || (p.MaxPrice != nil && *p.MaxPrice < 0) {
		return ErrInvalidPrice
	}
	if p.MinPrice != nil && p.MaxPrice != nil && *p.MinPrice > *p.MaxPrice {
		return ErrInvalidPriceRange
	}
	return nil
}

var ErrInvalidPriceRange = errors.New("min_price cannot be greater than max_price")

type ProductListResult struct {
	Products   []Product
	TotalCount int64
}
