package mongodb

import (
	"catalog-service/internal/entity"
	"catalog-service/internal/usecase"
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CategoryRepo struct {
	collection *mongo.Collection
}

var _ usecase.CategoryRepository = (*CategoryRepo)(nil)

func NewCategoryRepo(db *mongo.Database) *CategoryRepo {
	return &CategoryRepo{
		collection: db.Collection("categories"),
	}
}

func (r *CategoryRepo) Create(ctx context.Context, c *entity.Category) error { return nil }
func (r *CategoryRepo) FindByID(ctx context.Context, id string) (*entity.Category, error) {
	return nil, nil
}
func (r *CategoryRepo) FindChildren(ctx context.Context, parentID string) ([]entity.Category, error) {
	return []entity.Category{}, nil
}

