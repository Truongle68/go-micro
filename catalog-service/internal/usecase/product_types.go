package usecase

import (
	"catalog-service/internal/domain"
	"time"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type DetailedProduct struct {
	Product  domain.Product
	Category *domain.Category
}

type OptionTypeInput struct {
	Name   string
	Values []string
}

type CreateVariantInput struct {
	ID         string
	SKU        string
	Attributes map[string]string
	Price      PriceInput
	Image      string
}

type PriceInput struct {
	Amount   int
	Currency string
}

type SpecGroupInput struct {
	Group string
	Items []SpecItemInput
}

type SpecItemInput struct {
	Label string
	Value string
}

type CreateProductInput struct {
	Name            string
	NameTranslation map[string]string
	CategoryID      string
	Description     string
	DescriptionHTML string
	Highlights      []string
	Tags            []string
	Images          []string
	OptionTypes     []OptionTypeInput
	Variants        []CreateVariantInput
	Specifications  []SpecGroupInput
	Status          string
}

type UpdateProductInput struct {
	ID              string
	Version         int
	Name            *string
	NameTranslation map[string]string
	CategoryID      *string
	Description     *string
	DescriptionHTML *string
	Highlights      []string
	Tags            []string
	Images          []string
	OptionTypes     []OptionTypeInput
	Variants        []CreateVariantInput
	Specifications  []SpecGroupInput
	Status          *domain.ProductStatus
}

type ProductList struct {
	Products   []DetailedProduct `json:"products"`
	TotalCount int64             `json:"total_count"`
}

type VariantView struct {
	domain.Variant
	ProductID   string
	ProductName string
}

type PublicListInput struct {
	CategoryID string
	Query      string
	MinPrice   *int
	MaxPrice   *int
	Sort       []domain.SortKey
	Page       pagination.Params
}

type AdminListInput struct {
	Statuses    []domain.ProductStatus
	CategoryID  string
	Query       string
	MinPrice    *int
	MaxPrice    *int
	SKU         string
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	Sort        []domain.SortKey
	Page        pagination.Params
}
