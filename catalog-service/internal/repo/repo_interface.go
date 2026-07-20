package repo

import (
	"catalog-service/internal/domain"
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, p *domain.Product) error
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	FindByCategory(ctx context.Context, categoryID string) ([]domain.Product, error)
	Update(ctx context.Context, id string, p *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Product, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) error
	FindByID(ctx context.Context, id string) (*domain.Category, error)
	FindChildren(ctx context.Context, parentID string) ([]domain.Category, error)
	Update(ctx context.Context, id string, c *domain.Category) (*domain.Category, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Category, error)
}
