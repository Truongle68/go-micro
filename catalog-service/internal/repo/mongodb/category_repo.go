package repo

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/repo"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CategoryRepo struct {
	collection *mongo.Collection
}

var _ repo.CategoryRepository = (*CategoryRepo)(nil)

func NewCategoryRepo(db *mongo.Database) *CategoryRepo {
	return &CategoryRepo{
		collection: db.Collection("categories"),
	}
}

func (r *CategoryRepo) Create(ctx context.Context, c *domain.Category) error {
	if c.ID == "" {
		c.ID = bson.NewObjectID().Hex()
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, c)
	return err
}

func (r *CategoryRepo) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	var c domain.Category
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrCategoryNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepo) FindChildren(ctx context.Context, parentID string) ([]domain.Category, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"parent_id": parentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []domain.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	if categories == nil {
		return []domain.Category{}, nil
	}
	return categories, nil
}

func (r *CategoryRepo) Update(ctx context.Context, id string, c *domain.Category) (*domain.Category, error) {
	c.UpdatedAt = time.Now()
	filter := bson.M{"_id": id}
	c.ID = id
	res, err := r.collection.ReplaceOne(ctx, filter, c)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, domain.ErrCategoryNotFound
	}
	return c, nil
}

func (r *CategoryRepo) Delete(ctx context.Context, id string) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrCategoryNotFound
	}
	return nil
}

func (r *CategoryRepo) List(ctx context.Context) ([]domain.Category, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []domain.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	if categories == nil {
		return []domain.Category{}, nil
	}
	return categories, nil
}
