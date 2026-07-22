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

type ProductRepo struct {
	collection *mongo.Collection
}

var _ repo.ProductRepository = (*ProductRepo)(nil)

func NewProductRepo(db *mongo.Database) *ProductRepo {
	return &ProductRepo{
		collection: db.Collection("products"),
	}
}

func (r *ProductRepo) Create(ctx context.Context, p *domain.Product) error {
	if p.ID == "" {
		p.ID = bson.NewObjectID().Hex()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	// Ensure sub-variants and images have IDs
	for i := range p.Variants {
		if p.Variants[i].ID == "" {
			p.Variants[i].ID = bson.NewObjectID().Hex()
		}
	}
	for i := range p.Images {
		if p.Images[i].ID == "" {
			p.Images[i].ID = bson.NewObjectID().Hex()
		}
	}

	_, err := r.collection.InsertOne(ctx, p)
	return err
}

func (r *ProductRepo) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	var p domain.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepo) FindByCategory(ctx context.Context, categoryID string) ([]domain.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"category_id": categoryID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	if products == nil {
		return []domain.Product{}, nil
	}
	return products, nil
}

func (r *ProductRepo) Update(ctx context.Context, id string, p *domain.Product) (*domain.Product, error) {
	p.UpdatedAt = time.Now()
	p.ID = id

	// Ensure sub-variants and images have IDs and times
	for i := range p.Variants {
		if p.Variants[i].ID == "" {
			p.Variants[i].ID = bson.NewObjectID().Hex()
		}
	}
	for i := range p.Images {
		if p.Images[i].ID == "" {
			p.Images[i].ID = bson.NewObjectID().Hex()
		}
	}

	filter := bson.M{"_id": id}
	res, err := r.collection.ReplaceOne(ctx, filter, p)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, domain.ErrProductNotFound
	}
	return p, nil
}

func (r *ProductRepo) Delete(ctx context.Context, id string) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *ProductRepo) List(ctx context.Context) ([]domain.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	if products == nil {
		return []domain.Product{}, nil
	}
	return products, nil
}
