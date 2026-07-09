package entity

import "time"

type Product struct {
	ID            string    `json:"id"`
	CategoryID    string    `json:"category_id"`
	Sku           string    `json:"sku"`
	NameVi        string    `json:"name_vi"`
	NameEn        string    `json:"name_en"`
	DescriptionVi string    `json:"description_vi"`
	DescriptionEn string    `json:"description_en"`
	Unit          string    `json:"unit"`
	BasePrice     float64   `json:"base_price"`
	SalePrice     float64   `json:"sale_price"`
	RatingAvg     float64   `json:"rating_avg"`
	RatingCount   int32     `json:"rating_count"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProductVariants struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	VariantLabel string    `json:"variant_label"`
	PriceDelta   float64   `json:"price_delta"`
	Sku          string    `json:"sku"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductImage struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Url       string    `json:"url"`
	SortOrder int32     `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
