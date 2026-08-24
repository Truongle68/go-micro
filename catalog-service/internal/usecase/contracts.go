package usecase

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/repo/mongorepo"
	"context"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type ProductRepository interface {
	EnsureIndexes(ctx context.Context) error
	Create(ctx context.Context, p *domain.Product) error
	ExistSlug(ctx context.Context, name string) (bool, error)
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	FindByCategory(ctx context.Context, categoryID string, p pagination.Params) (*domain.ProductListResult, error)
	Update(ctx context.Context, p *domain.Product, expectedVersion int) (*domain.Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, p pagination.Params) (*domain.ProductListResult, error)
	Search(ctx context.Context, sParams domain.SearchProductParams, pParams pagination.Params) (*domain.ProductListResult, error)
	FindByVariantSKUs(ctx context.Context, skus []string) ([]domain.Product, error)
}

var _ ProductRepository = (*mongorepo.ProductRepo)(nil)

type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) error
	BuildBreadcrumb(ctx context.Context, id string) ([]domain.CategoryRef, error)
	FindByID(ctx context.Context, id string) (*domain.Category, error)
	FindByIDs(ctx context.Context, ids []string) (*domain.ListCategoryResult, error)
	FindChildren(ctx context.Context, parentID string, p pagination.Params) (*domain.ListCategoryResult, error)
	Update(ctx context.Context, c *domain.Category) (*domain.Category, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, p pagination.Params) (*domain.ListCategoryResult, error)
}

var _ CategoryRepository = (*mongorepo.CategoryRepo)(nil)
