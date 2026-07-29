package req

import "catalog-service/internal/usecase"

type ProductVariant struct {
	VariantLabel string  `json:"variant_label" validate:"required"`
	PriceDelta   float64 `json:"price_delta"`
	Sku          string  `json:"sku" validate:"required"`
}

type CreateProduct struct {
	CategoryID    string           `json:"category_id" validate:"required"`
	Sku           string           `json:"sku" validate:"required"`
	NameVi        string           `json:"name_vi" validate:"required"`
	NameEn        string           `json:"name_en" validate:"required"`
	DescriptionVi string           `json:"description_vi"`
	DescriptionEn string           `json:"description_en"`
	Unit          string           `json:"unit" validate:"required"`
	BasePrice     float64          `json:"base_price" validate:"required,gte=0"`
	SalePrice     float64          `json:"sale_price" validate:"required,gte=0"`
	IsActive      bool             `json:"is_active"`
	Variants      []ProductVariant `json:"variants"`
	Images        []string         `json:"images"`
}

func (req *CreateProduct) ToCreateProductInput() usecase.CreateProductInput {
	var variants []usecase.ProductVariantInput
	for _, v := range req.Variants {
		variants = append(variants, usecase.ProductVariantInput{
			VariantLabel: v.VariantLabel,
			PriceDelta:   v.PriceDelta,
			Sku:          v.Sku,
		})
	}

	return usecase.CreateProductInput{
		CategoryID:    req.CategoryID,
		Sku:           req.Sku,
		NameVi:        req.NameVi,
		NameEn:        req.NameEn,
		DescriptionVi: req.DescriptionVi,
		DescriptionEn: req.DescriptionEn,
		Unit:          req.Unit,
		BasePrice:     req.BasePrice,
		SalePrice:     req.SalePrice,
		IsActive:      req.IsActive,
		Variants:      variants,
		Images:        req.Images,
	}
}

type UpdateProduct struct {
	CategoryID    *string          `json:"category_id"`
	Sku           *string          `json:"sku"`
	NameVi        *string          `json:"name_vi"`
	NameEn        *string          `json:"name_en"`
	DescriptionVi *string          `json:"description_vi"`
	DescriptionEn *string          `json:"description_en"`
	Unit          *string          `json:"unit"`
	BasePrice     *float64         `json:"base_price" validate:"omitempty,gte=0"`
	SalePrice     *float64         `json:"sale_price" validate:"omitempty,gte=0"`
	IsActive      *bool            `json:"is_active"`
	Variants      []ProductVariant `json:"variants"`
	Images        []string         `json:"images"`
}

func (req *UpdateProduct) ToUpdateProductInput(id string) usecase.UpdateProductInput {
	var variants []usecase.ProductVariantInput
	if req.Variants != nil {
		for _, v := range req.Variants {
			variants = append(variants, usecase.ProductVariantInput{
				VariantLabel: v.VariantLabel,
				PriceDelta:   v.PriceDelta,
				Sku:          v.Sku,
			})
		}
	}

	return usecase.UpdateProductInput{
		ID:            id,
		CategoryID:    req.CategoryID,
		Sku:           req.Sku,
		NameVi:        req.NameVi,
		NameEn:        req.NameEn,
		DescriptionVi: req.DescriptionVi,
		DescriptionEn: req.DescriptionEn,
		Unit:          req.Unit,
		BasePrice:     req.BasePrice,
		SalePrice:     req.SalePrice,
		IsActive:      req.IsActive,
		Variants:      variants,
		Images:        req.Images,
	}
}

type SearchProduct struct {
	Query      string   `json:"query"`
	CategoryID string   `json:"category_id"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	IsActive   *bool    `json:"is_active"`
	Page       int64    `json:"page"`
	Limit      int64    `json:"limit"`
}
