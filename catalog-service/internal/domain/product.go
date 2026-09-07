package domain

import (
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
	Amount   int
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
	Images          []string
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
	}
	p.Images = params.Images
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

type ProductListFilter struct {
	Statuses    []ProductStatus
	CategoryID  string
	Query       string
	MinPrice    *int
	MaxPrice    *int
	SKU         string
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	Sort        []SortKey
}

type ProductListResult struct {
	Products   []Product
	TotalCount int64
}

type SearchProductParams struct {
	Query      string
	CategoryID string
	MinPrice   *int
	MaxPrice   *int
	Status     *ProductStatus
}

type SearchProductsQuery struct {
	Query      string         `form:"q"`
	CategoryID string         `form:"category_id"`
	MinPrice   *int           `form:"min_price"`
	MaxPrice   *int           `form:"max_price"`
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

type ProductStatusCounts struct {
	All      int `json:"all"`
	Draft    int `json:"draft"`
	Active   int `json:"active"`
	Inactive int `json:"inactive"`
	Archived int `json:"archived"`
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

type SortDir string

const (
	SortAsc  SortDir = "asc"
	SortDesc SortDir = "desc"
)

func (d SortDir) IsValid() bool {
	switch d {
	case SortAsc, SortDesc:
		return true
	default:
		return false
	}
}

type SortField string

const (
	SortByName  SortField = "name"
	SortByPrice SortField = "price"
	SortByDate  SortField = "date"
)

func (f SortField) IsValid() bool {
	switch f {
	case SortByName, SortByDate, SortByPrice:
		return true
	default:
		return false
	}
}

type SortKey struct {
	Field SortField `json:"field"`
	Dir   SortDir   `json:"dir"`
}
