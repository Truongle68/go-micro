package mongodb

import (
	"catalog-service/internal/entity"
	"catalog-service/internal/usecase"
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductRepo struct {
	collection *mongo.Collection
}

var _ usecase.ProductRepository = (*ProductRepo)(nil)

func NewProductRepo(db *mongo.Database) *ProductRepo {
	return &ProductRepo{
		collection: db.Collection("products"),
	}
}

func (r *ProductRepo) Create(ctx context.Context, p *entity.Product) error { return nil }
func (r *ProductRepo) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	return nil, nil
}
func (r *ProductRepo) FindByCategory(ctx context.Context, categoryID string) ([]entity.Product, error) {
	return []entity.Product{}, nil
}
func (r *ProductRepo) Update(ctx context.Context, id string, p *entity.Product) (*entity.Product, error) {
	return nil, nil
}
func (r *ProductRepo) Delete(ctx context.Context, id string) error { return nil }
