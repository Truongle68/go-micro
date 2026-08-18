package res

import (
	"time"
)

type ProductResponse struct {
	ID        string            `json:"id"`
	Slug      string            `json:"slug"`
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	Variants  []VariantResponse `json:"variants"`
	CreatedAt time.Time         `json:"created_at"`
}

type VariantResponse struct {
	ID    string `json:"id"`
	SKU   string `json:"sku"`
	Price int64  `json:"price"`
}

type CategoryRefRead struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type OptionTypeRead struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type PriceRead struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type VariantRead struct {
	ID         string            `json:"id,omitempty"`
	SKU        string            `json:"sku"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Price      PriceRead         `json:"price"`
	Image      string            `json:"image,omitempty"`
	IsActive   bool              `json:"is_active"`
	CreatedAt  time.Time         `json:"created_at"`
}

type SpecItemRead struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type SpecGroupRead struct {
	Group string         `json:"group"`
	Items []SpecItemRead `json:"items"`
}

type ProductCategoryRead struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	NameTranslation map[string]string `json:"name_translation,omitempty"`
	Slug            string            `json:"slug"`
	Icon            string            `json:"icon,omitempty"`
	SortOrder       int64             `json:"sort_order"`
	IsActive        bool              `json:"is_active"`
}

type ProductRead struct {
	ID              string               `json:"id"`
	Version         int                  `json:"version"`
	Slug            string               `json:"slug"`
	Name            string               `json:"name"`
	NameTranslation map[string]string    `json:"name_translation,omitempty"`
	CategoryID      string               `json:"category_id"`
	CategoryPath    []CategoryRefRead    `json:"category_path,omitempty"`
	Category        *ProductCategoryRead `json:"category,omitempty"`
	Description     string               `json:"description"`
	DescriptionHTML string               `json:"description_html,omitempty"`
	Highlights      []string             `json:"highlights,omitempty"`
	Tags            []string             `json:"tags,omitempty"`
	Images          []string             `json:"images,omitempty"`
	OptionTypes     []OptionTypeRead     `json:"option_types,omitempty"`
	Variants        []VariantRead        `json:"variants,omitempty"`
	Specifications  []SpecGroupRead      `json:"specifications,omitempty"`
	Status          string               `json:"status"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}
