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
	Stock int    `json:"stock"`
}

type CategoryRefRead struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ImageRead struct {
	URL       string `json:"url"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
	AltText   string `json:"alt_text"`
}

type OptionTypeRead struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type PriceRead struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type InventoryRead struct {
	TotalAvailable int `json:"total_available"`
	Reserved       int `json:"reserved"`
}

type VariantRead struct {
	ID          string            `json:"id,omitempty"`
	SKU         string            `json:"sku"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	Price       PriceRead         `json:"price"`
	Inventory   InventoryRead     `json:"inventory"`
	WeightGrams int               `json:"weight_grams,omitempty"`
	Images      []ImageRead       `json:"images,omitempty"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
}

type SpecItemRead struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type SpecGroupRead struct {
	Group string         `json:"group"`
	Items []SpecItemRead `json:"items"`
}

type RatingSummaryRead struct {
	Average   float64          `json:"average"`
	Count     int              `json:"count"`
	Breakdown map[string]int64 `json:"breakdown,omitempty"`
}

type ShipsFromRead struct {
	WarehouseID string `json:"warehouse_id,omitempty"`
	Region      string `json:"region,omitempty"`
}

type ShippingInfoRead struct {
	IsFreeShipping bool          `json:"is_free_shipping"`
	ShipsFrom      ShipsFromRead `json:"ships_from,omitempty"`
	Fragile        bool          `json:"fragile"`
	ShippingClass  string        `json:"shipping_class,omitempty"`
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
	Images          []ImageRead          `json:"images,omitempty"`
	OptionTypes     []OptionTypeRead     `json:"option_types,omitempty"`
	Variants        []VariantRead        `json:"variants,omitempty"`
	Specifications  []SpecGroupRead      `json:"specifications,omitempty"`
	RatingSummary   RatingSummaryRead    `json:"rating_summary"`
	SalesCount      int64                `json:"sales_count"`
	Shipping        ShippingInfoRead     `json:"shipping"`
	Status          string               `json:"status"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}
