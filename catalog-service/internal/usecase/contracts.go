package usecase

import (
	"context"
	"time"
)

type CategoryDTO struct {
	ID        string    `json:"id"`
	ParentID  *string   `json:"parent_id"`
	NameVi    string    `json:"name_vi"`
	NameEn    string    `json:"name_en"`
	Slug      string    `json:"slug"`
	Icon      string    `json:"icon"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCategoryInput struct {
	ParentID  *string
	NameVi    string
	NameEn    string
	Slug      string
	Icon      string
	SortOrder int
}

type UpdateCategoryInput struct {
	ID        string
	ParentID  *string
	NameVi    *string
	NameEn    *string
	Slug      *string
	Icon      *string
	SortOrder *int
}

type CategoryUsecase interface {
	Create(ctx context.Context, in CreateCategoryInput) (*CategoryDTO, error)
	GetByID(ctx context.Context, id string) (*CategoryDTO, error)
	GetChildren(ctx context.Context, parentID string) ([]*CategoryDTO, error)
	Update(ctx context.Context, in UpdateCategoryInput) (*CategoryDTO, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*CategoryDTO, error)
}

type ProductVariantDTO struct {
	ID           string    `json:"id"`
	VariantLabel string    `json:"variant_label"`
	PriceDelta   float64   `json:"price_delta"`
	Sku          string    `json:"sku"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductImageDTO struct {
	ID        string    `json:"id"`
	Url       string    `json:"url"`
	SortOrder int32     `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductDTO struct {
	ID            string               `json:"id"`
	CategoryID    string               `json:"category_id"`
	Sku           string               `json:"sku"`
	NameVi        string               `json:"name_vi"`
	NameEn        string               `json:"name_en"`
	DescriptionVi string               `json:"description_vi"`
	DescriptionEn string               `json:"description_en"`
	Unit          string               `json:"unit"`
	BasePrice     float64              `json:"base_price"`
	SalePrice     float64              `json:"sale_price"`
	RatingAvg     float64              `json:"rating_avg"`
	RatingCount   int32                `json:"rating_count"`
	IsActive      bool                 `json:"is_active"`
	Variants      []*ProductVariantDTO `json:"variants"`
	Images        []*ProductImageDTO   `json:"images"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

type ProductVariantInput struct {
	VariantLabel string
	PriceDelta   float64
	Sku          string
}

type ProductImageInput struct {
	Url       string
	SortOrder int32
}

type CreateProductInput struct {
	CategoryID    string
	Sku           string
	NameVi        string
	NameEn        string
	DescriptionVi string
	DescriptionEn string
	Unit          string
	BasePrice     float64
	SalePrice     float64
	IsActive      bool
	Variants      []ProductVariantInput
	Images        []ProductImageInput
}

type UpdateProductInput struct {
	ID            string
	CategoryID    *string
	Sku           *string
	NameVi        *string
	NameEn        *string
	DescriptionVi *string
	DescriptionEn *string
	Unit          *string
	BasePrice     *float64
	SalePrice     *float64
	IsActive      *bool
	Variants      []ProductVariantInput
	Images        []ProductImageInput
}

type ProductUsecase interface {
	Create(ctx context.Context, in CreateProductInput) (*ProductDTO, error)
	GetByID(ctx context.Context, id string) (*ProductDTO, error)
	GetByCategory(ctx context.Context, categoryID string) ([]*ProductDTO, error)
	Update(ctx context.Context, in UpdateProductInput) (*ProductDTO, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*ProductDTO, error)
}
