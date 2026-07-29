package usecase

import (
	"catalog-service/internal/domain"
	"context"
	"time"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type CategoryDTO struct {
	ID        string    `json:"id"`
	ParentID  *string   `json:"parent_id"`
	NameVi    string    `json:"name_vi"`
	NameEn    string    `json:"name_en"`
	Slug      string    `json:"slug"`
	Icon      string    `json:"icon"`
	SortOrder int64     `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCategoryInput struct {
	ParentID  *string
	NameVi    string
	NameEn    string
	Slug      string
	Icon      string
	SortOrder int64
}

type UpdateCategoryInput struct {
	ID        string
	ParentID  *string
	NameVi    *string
	NameEn    *string
	Slug      *string
	Icon      *string
	SortOrder *int64
}

type CategoryList struct {
	Categories []CategoryDTO `json:"categories"`
	TotalCount int64         `json:"total_count"`
}

type CategoryUsecase interface {
	Create(ctx context.Context, in CreateCategoryInput) (*CategoryDTO, error)
	GetByID(ctx context.Context, id string) (*CategoryDTO, error)
	GetChildren(ctx context.Context, parentID string, params pagination.Params) (*CategoryList, error)
	Update(ctx context.Context, in UpdateCategoryInput) (*CategoryDTO, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params pagination.Params) (*CategoryList, error)
}

type DetailedProduct struct {
	Product  domain.Product
	Category *domain.Category
}

type ProductVariantInput struct {
	VariantLabel string
	PriceDelta   float64
	Sku          string
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
	Images        []string
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
	Images        []string
}

type ProductList struct {
	Products   []DetailedProduct `json:"products"`
	TotalCount int64             `json:"total_count"`
}

type ProductUsecase interface {
	Create(ctx context.Context, in CreateProductInput) (*domain.Product, error)
	GetByID(ctx context.Context, id string) (*DetailedProduct, error)
	GetByCategory(ctx context.Context, categoryID string, params pagination.Params) (*ProductList, error)
	Update(ctx context.Context, in UpdateProductInput) (*domain.Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params pagination.Params) (*ProductList, error)
	Search(ctx context.Context, searchParams domain.SearchProductParams, paginatedParams pagination.Params) (*ProductList, error)
}
