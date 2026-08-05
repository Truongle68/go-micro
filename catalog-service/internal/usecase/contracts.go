package usecase

import (
	"catalog-service/internal/domain"
	"context"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type CategoryUsecase interface {
	Create(ctx context.Context, in CreateCategoryInput) (*CategoryDTO, error)
	GetByID(ctx context.Context, id string) (*CategoryDTO, error)
	GetChildren(ctx context.Context, parentID string, params pagination.Params) (*CategoryList, error)
	Update(ctx context.Context, in UpdateCategoryInput) (*CategoryDTO, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params pagination.Params) (*CategoryList, error)
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
