package usecase

import (
	"catalog-service/internal/entity"
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, p *entity.Product) error
	FindByID(ctx context.Context, id string) (*entity.Product, error)
	FindByCategory(ctx context.Context, categoryID string) ([]entity.Product, error)
	Update(ctx context.Context, id string, p *entity.Product) (*entity.Product, error)
	Delete(ctx context.Context, id string) error
}

type CategoryRepository interface {
	Create(ctx context.Context, c *entity.Category) error
	FindByID(ctx context.Context, id string) (*entity.Category, error)
	FindChildren(ctx context.Context, parentID string) ([]entity.Category, error)
}
