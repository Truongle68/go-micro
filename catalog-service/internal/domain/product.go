package domain

import (
	"catalog-service/pkg/sliceutil"
	"time"
)

type Product struct {
	ID              string
	Version         int
	Slug            string
	Name            string
	NameTranslation map[string]string
	CategoryID      string
	CategoryPath    []CategoryRef
	Description     string
	DescriptionHTML string
	Highlights      []string
	Tags            []string
	Images          []string
	OptionTypes     []OptionType
	Variants        []Variant
	Specifications  []SpecGroup
	Status          ProductStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Variant struct {
	ID         string
	SKU        string
	Attributes map[string]string
	Price      Price
	Image      string
	IsActive   bool
	CreatedAt  time.Time
}

type OptionType struct {
	Name   string
	Values []string
}

type Price struct {
	Amount   int64
	Currency string
}

type SpecGroup struct {
	Group string
	Items []SpecItem
}

type SpecItem struct {
	Label string
	Value string
}

type CategoryRef struct {
	ID   string
	Name string
}

type UpdateProductParams struct {
	Name            *string
	NameTranslation map[string]string
	CategoryID      *string
	CategoryPath    []CategoryRef
	Description     *string
	DescriptionHTML *string
	Highlights      []string
	Tags            []string
	OptionTypes     []OptionType
	Variants        []Variant
	Specifications  []SpecGroup
	Status          *ProductStatus
}

func (p *Product) ApplyUpdate(params UpdateProductParams) error {
	if params.Name != nil {
		if *params.Name == "" {
			return ErrEmptyName
		}
		p.Name = *params.Name
	}
	if params.NameTranslation != nil {
		p.NameTranslation = params.NameTranslation
	}
	if params.CategoryID != nil {
		if *params.CategoryID == "" {
			return ErrEmptyCategoryID
		}
		p.CategoryID = *params.CategoryID
	}
	if params.CategoryPath != nil {
		p.CategoryPath = params.CategoryPath
	}
	if params.Description != nil {
		p.Description = *params.Description
	}
	if params.DescriptionHTML != nil {
		p.DescriptionHTML = *params.DescriptionHTML
	}
	if params.Highlights != nil {
		p.Highlights = params.Highlights
	}
	if params.Tags != nil {
		p.Tags = params.Tags
	}
	if params.OptionTypes != nil {
		p.OptionTypes = params.OptionTypes
	}
	if params.Variants != nil {
		if len(params.Variants) == 0 {
			return ErrProductRequiresVariant
		}
		p.Variants = params.Variants
		p.Images = sliceutil.DedupeString(params.Variants, func(v Variant) string {
			return v.Image
		})
	}
	if params.Specifications != nil {
		p.Specifications = params.Specifications
	}
	if params.Status != nil {
		if !params.Status.IsValid() {
			return ErrInvalidProductStatus
		}
		p.Status = *params.Status
	}
	p.UpdatedAt = time.Now()
	return nil
}

type SearchProductParams struct {
	Query      string
	CategoryID string
	MinPrice   *int64
	MaxPrice   *int64
	Status     *ProductStatus
}

type SearchProductsQuery struct {
	Query      string         `form:"q"`
	CategoryID string         `form:"category_id"`
	MinPrice   *int64         `form:"min_price"`
	MaxPrice   *int64         `form:"max_price"`
	Status     *ProductStatus `form:"status"`
}

func (q SearchProductsQuery) ToDomainParams() SearchProductParams {
	return SearchProductParams{
		Query:      q.Query,
		CategoryID: q.CategoryID,
		MinPrice:   q.MinPrice,
		MaxPrice:   q.MaxPrice,
		Status:     q.Status,
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

type ProductListResult struct {
	Products   []Product
	TotalCount int64
}

type ProductStatus string

const (
	ProductStatusDraft    ProductStatus = "draft"
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusArchived ProductStatus = "archived"
)

func (s ProductStatus) IsValid() bool {
	switch s {
	case ProductStatusDraft, ProductStatusActive, ProductStatusInactive, ProductStatusArchived:
		return true
	default:
		return false
	}
}
