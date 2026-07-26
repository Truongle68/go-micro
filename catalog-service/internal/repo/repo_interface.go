package repo

import (
	"catalog-service/internal/domain"
	"context"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type ProductRepository interface {
	Create(ctx context.Context, p *domain.Product) error
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	FindByCategory(ctx context.Context, categoryID string) ([]domain.Product, error)
	Update(ctx context.Context, p *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, p pagination.Params) ([]domain.Product, error)
	Search(ctx context.Context, sParams domain.SearchProductParams, pParams pagination.Params) (*domain.ProductSearchResult, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) error
	FindByID(ctx context.Context, id string) (*domain.Category, error)
	FindChildren(ctx context.Context, parentID string) ([]domain.Category, error)
	Update(ctx context.Context, id string, c *domain.Category) (*domain.Category, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, p pagination.Params) ([]domain.Category, error)
}
