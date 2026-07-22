package domain

import "time"

type Product struct {
	ID            string           `json:"id" bson:"_id,omitempty"`
	CategoryID    string           `json:"category_id" bson:"category_id"`
	Sku           string           `json:"sku" bson:"sku"`
	NameVi        string           `json:"name_vi" bson:"name_vi"`
	NameEn        string           `json:"name_en" bson:"name_en"`
	DescriptionVi string           `json:"description_vi" bson:"description_vi"`
	DescriptionEn string           `json:"description_en" bson:"description_en"`
	Unit          string           `json:"unit" bson:"unit"`
	BasePrice     float64          `json:"base_price" bson:"base_price"`
	SalePrice     float64          `json:"sale_price" bson:"sale_price"`
	RatingAvg     float64          `json:"rating_avg" bson:"rating_avg"`
	RatingCount   int32            `json:"rating_count" bson:"rating_count"`
	IsActive      bool             `json:"is_active" bson:"is_active"`
	Variants      []ProductVariant `json:"variants" bson:"variants,omitempty"`
	Images        []ProductImage   `json:"images" bson:"images,omitempty"`
	CreatedAt     time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" bson:"updated_at"`
}

type ProductVariant struct {
	ID           string  `json:"id" bson:"id"`
	VariantLabel string  `json:"variant_label" bson:"variant_label"`
	PriceDelta   float64 `json:"price_delta" bson:"price_delta"`
	Sku          string  `json:"sku" bson:"sku"`
}

type ProductImage struct {
	ID        string `json:"id" bson:"id"`
	Url       string `json:"url" bson:"url"`
	SortOrder int32  `json:"sort_order" bson:"sort_order"`
}
