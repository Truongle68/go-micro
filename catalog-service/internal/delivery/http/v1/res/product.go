package res

import (
	"catalog-service/internal/usecase"
	"time"
)

type ProductVariantResponse struct {
	ID           string    `json:"id"`
	VariantLabel string    `json:"variant_label"`
	PriceDelta   float64   `json:"price_delta"`
	Sku          string    `json:"sku"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductImageResponse struct {
	ID        string    `json:"id"`
	Url       string    `json:"url"`
	SortOrder int32     `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductResponse struct {
	ID            string                    `json:"id"`
	CategoryID    string                    `json:"category_id"`
	Sku           string                    `json:"sku"`
	NameVi        string                    `json:"name_vi"`
	NameEn        string                    `json:"name_en"`
	DescriptionVi string                    `json:"description_vi"`
	DescriptionEn string                    `json:"description_en"`
	Unit          string                    `json:"unit"`
	BasePrice     float64                   `json:"base_price"`
	SalePrice     float64                   `json:"sale_price"`
	RatingAvg     float64                   `json:"rating_avg"`
	RatingCount   int32                     `json:"rating_count"`
	IsActive      bool                      `json:"is_active"`
	Variants      []*ProductVariantResponse `json:"variants"`
	Images        []*ProductImageResponse   `json:"images"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

func ToProductResponse(dto *usecase.ProductDTO) ProductResponse {
	if dto == nil {
		return ProductResponse{}
	}

	variants := make([]*ProductVariantResponse, len(dto.Variants))
	for i, v := range dto.Variants {
		variants[i] = &ProductVariantResponse{
			ID:           v.ID,
			VariantLabel: v.VariantLabel,
			PriceDelta:   v.PriceDelta,
			Sku:          v.Sku,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		}
	}

	images := make([]*ProductImageResponse, len(dto.Images))
	for i, img := range dto.Images {
		images[i] = &ProductImageResponse{
			ID:        img.ID,
			Url:       img.Url,
			SortOrder: img.SortOrder,
			CreatedAt: img.CreatedAt,
			UpdatedAt: img.UpdatedAt,
		}
	}

	return ProductResponse{
		ID:            dto.ID,
		CategoryID:    dto.CategoryID,
		Sku:           dto.Sku,
		NameVi:        dto.NameVi,
		NameEn:        dto.NameEn,
		DescriptionVi: dto.DescriptionVi,
		DescriptionEn: dto.DescriptionEn,
		Unit:          dto.Unit,
		BasePrice:     dto.BasePrice,
		SalePrice:     dto.SalePrice,
		RatingAvg:     dto.RatingAvg,
		RatingCount:   dto.RatingCount,
		IsActive:      dto.IsActive,
		Variants:      variants,
		Images:        images,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
	}
}

func ToProductListResponse(dtos []*usecase.ProductDTO) []ProductResponse {
	res := make([]ProductResponse, len(dtos))
	for i, dto := range dtos {
		res[i] = ToProductResponse(dto)
	}
	return res
}
