package res

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/usecase"
	"time"
)

type ProductVariant struct {
	VariantLabel string  `json:"variant_label"`
	PriceDelta   float64 `json:"price_delta"`
	Sku          string  `json:"sku"`
}

func toProductVariant(v domain.ProductVariant) ProductVariant {
	return ProductVariant{
		VariantLabel: v.VariantLabel,
		PriceDelta:   v.PriceDelta,
		Sku:          v.Sku,
	}
}

type DetailedCategory struct {
	ID        string `json:"id"`
	NameVi    string `json:"name_vi"`
	NameEn    string `json:"name_en"`
	Slug      string `json:"slug"`
	Icon      string `json:"icon"`
	SortOrder int64  `json:"sort_order"`
}

func toDetailedCategory(c domain.Category) *DetailedCategory {
	return &DetailedCategory{
		ID:        c.ID,
		NameVi:    c.NameVi,
		NameEn:    c.NameEn,
		Slug:      c.Slug,
		Icon:      c.Icon,
		SortOrder: c.SortOrder,
	}
}

type ProductRead struct {
	ID            string            `json:"id"`
	CategoryID    *string           `json:"category_id"`
	Category      *DetailedCategory `json:"category"`
	Sku           string            `json:"sku"`
	NameVi        string            `json:"name_vi"`
	NameEn        string            `json:"name_en"`
	DescriptionVi string            `json:"description_vi"`
	DescriptionEn string            `json:"description_en"`
	Unit          string            `json:"unit"`
	BasePrice     float64           `json:"base_price"`
	SalePrice     float64           `json:"sale_price"`
	RatingAvg     float64           `json:"rating_avg"`
	RatingCount   int64             `json:"rating_count"`
	IsActive      bool              `json:"is_active"`
	Variants      []ProductVariant  `json:"variants"`
	Images        []string          `json:"images"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

func ToProductRead(p *domain.Product) ProductRead {
	if p == nil {
		return ProductRead{}
	}

	variants := make([]ProductVariant, len(p.Variants))
	for i, v := range p.Variants {
		variants[i] = toProductVariant(v)
	}

	images := make([]string, len(p.Images))
	copy(images, p.Images)

	return ProductRead{
		ID:            p.ID,
		CategoryID:    &p.CategoryID,
		Sku:           p.Sku,
		NameVi:        p.NameVi,
		NameEn:        p.NameEn,
		DescriptionVi: p.DescriptionVi,
		DescriptionEn: p.DescriptionEn,
		Unit:          p.Unit,
		BasePrice:     p.BasePrice,
		SalePrice:     p.SalePrice,
		RatingAvg:     p.RatingAvg,
		RatingCount:   p.RatingCount,
		IsActive:      p.IsActive,
		Variants:      variants,
		Images:        images,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func ToDetailedProductRead(dp *usecase.DetailedProduct) ProductRead {
	if dp == nil {
		return ProductRead{}
	}

	variants := make([]ProductVariant, len(dp.Product.Variants))
	for i, v := range dp.Product.Variants {
		variants[i] = toProductVariant(v)
	}

	images := make([]string, len(dp.Product.Images))
	copy(images, dp.Product.Images)

	p := ProductRead{
		ID:            dp.Product.ID,
		Sku:           dp.Product.Sku,
		NameVi:        dp.Product.NameVi,
		NameEn:        dp.Product.NameEn,
		DescriptionVi: dp.Product.DescriptionVi,
		DescriptionEn: dp.Product.DescriptionEn,
		Unit:          dp.Product.Unit,
		BasePrice:     dp.Product.BasePrice,
		SalePrice:     dp.Product.SalePrice,
		RatingAvg:     dp.Product.RatingAvg,
		RatingCount:   dp.Product.RatingCount,
		IsActive:      dp.Product.IsActive,
		Variants:      variants,
		Images:        images,
		CreatedAt:     dp.Product.CreatedAt,
		UpdatedAt:     dp.Product.UpdatedAt,
	}

	if dp.Category != nil {
		p.Category = toDetailedCategory(*dp.Category)
	}

	return p
}

func ToProductList(dps []usecase.DetailedProduct) []ProductRead {
	products := make([]ProductRead, len(dps))
	for i, dp := range dps {
		products[i] = ToDetailedProductRead(&dp)
	}
	return products
}
