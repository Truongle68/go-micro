package repo

import (
	"catalog-service/internal/domain"
	"context"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type ProductRepository interface {
	Create(ctx context.Context, p *domain.Product) error
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	FindByCategory(ctx context.Context, categoryID string, p pagination.Params) (*domain.ListProductResult, error)
	Update(ctx context.Context, p *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, p pagination.Params) (*domain.ListProductResult, error)
	Search(ctx context.Context, sParams domain.SearchProductParams, pParams pagination.Params) (*domain.ListProductResult, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) error
	FindByID(ctx context.Context, id string) (*domain.Category, error)
	FindChildren(ctx context.Context, parentID string, p pagination.Params) (*domain.ListCategoryResult, error)
	Update(ctx context.Context, c *domain.Category) (*domain.Category, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, p pagination.Params) (*domain.ListCategoryResult, error)
}
