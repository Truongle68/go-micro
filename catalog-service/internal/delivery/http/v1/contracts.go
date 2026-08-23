package v1

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/usecase"
	"context"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type CategoryUsecase interface {
	Create(ctx context.Context, in usecase.CreateCategoryInput) (*usecase.CategoryDTO, error)
	GetByID(ctx context.Context, id string) (*usecase.CategoryDTO, error)
	GetChildren(ctx context.Context, parentID string, params pagination.Params) (*usecase.CategoryList, error)
	Update(ctx context.Context, in usecase.UpdateCategoryInput) (*usecase.CategoryDTO, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params pagination.Params) (*usecase.CategoryList, error)
}

var _ CategoryUsecase = (*usecase.CategoryUC)(nil)

type ProductUsecase interface {
	Create(ctx context.Context, in usecase.CreateProductInput) (*domain.Product, error)
	GetByID(ctx context.Context, id string) (*usecase.DetailedProduct, error)
	GetByCategory(ctx context.Context, categoryID string, params pagination.Params) (*usecase.ProductList, error)
	Update(ctx context.Context, in usecase.UpdateProductInput) (*domain.Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params pagination.Params) (*usecase.ProductList, error)
	Search(ctx context.Context, searchParams domain.SearchProductParams, paginatedParams pagination.Params) (*usecase.ProductList, error)
}

var _ ProductUsecase = (*usecase.ProductUC)(nil)
