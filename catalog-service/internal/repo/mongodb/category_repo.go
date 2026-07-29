package repo

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/repo"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

type categoryModel struct {
	ID        bson.ObjectID  `bson:"_id,omitempty"`
	ParentID  *bson.ObjectID `bson:"parent_id,omitempty"`
	NameVi    string         `bson:"name_vi"`
	NameEn    string         `bson:"name_en"`
	Slug      string         `bson:"slug"`
	Icon      string         `bson:"icon"`
	SortOrder int64          `bson:"sort_order"`
	CreatedAt time.Time      `bson:"created_at"`
	UpdatedAt time.Time      `bson:"updated_at"`
}

func (m categoryModel) toDomain() *domain.Category {
	var parentID *string
	if m.ParentID != nil && !m.ParentID.IsZero() {
		pid := m.ParentID.Hex()
		parentID = &pid
	}

	return &domain.Category{
		ID:        m.ID.Hex(),
		ParentID:  parentID,
		NameVi:    m.NameVi,
		NameEn:    m.NameEn,
		Slug:      m.Slug,
		Icon:      m.Icon,
		SortOrder: m.SortOrder,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func categoryModelFromDomain(c *domain.Category) (*categoryModel, error) {
	if c == nil {
		return nil, domain.ErrCategoryNotFound
	}
	now := time.Now()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	c.UpdatedAt = now

	m := &categoryModel{
		NameVi:    c.NameVi,
		NameEn:    c.NameEn,
		Slug:      c.Slug,
		Icon:      c.Icon,
		SortOrder: c.SortOrder,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if c.ID != "" {
		oid, err := bson.ObjectIDFromHex(c.ID)
		if err != nil {
			return nil, domain.ErrInvalidCategoryID
		}
		m.ID = oid
	}

	if c.ParentID != nil && *c.ParentID != "" {
		poid, err := bson.ObjectIDFromHex(*c.ParentID)
		if err != nil {
			return nil, domain.ErrInvalidCategoryID
		}
		m.ParentID = &poid
	}

	return m, nil
}

func (r *CategoryRepo) Create(ctx context.Context, c *domain.Category) error {
	m, err := categoryModelFromDomain(c)
	if err != nil {
		return err
	}
	if m.ID.IsZero() {
		m.ID = bson.NewObjectID()
		c.ID = m.ID.Hex()
	}

	res, err := r.collection.InsertOne(ctx, m)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		c.ID = oid.Hex()
	}
	return nil
}

func (r *CategoryRepo) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidCategoryID
	}

	var m categoryModel
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&m)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrCategoryNotFound
	}
	if err != nil {
		return nil, err
	}
	c := m.toDomain()
	return c, nil
}

func (r *CategoryRepo) findWithFilter(ctx context.Context, filter bson.M, p *pagination.Params) (*domain.ListCategoryResult, error) {
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count child categories: %w", err)
	}

	if total == 0 {
		return &domain.ListCategoryResult{
			Categories: []domain.Category{},
			TotalCount: 0,
		}, nil
	}

	findOpts := options.Find().
		SetSort(bson.D{{Key: "sort_order", Value: 1}, {Key: "_id", Value: 1}})

	if p != nil {
		findOpts.SetSkip(p.Skip()).SetLimit(p.Limit)
	}

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	return decodeCategoryList(ctx, cursor, total)
}

func (r *CategoryRepo) FindByIDs(ctx context.Context, ids []string) (*domain.ListCategoryResult, error) {
	if len(ids) == 0 {
		return &domain.ListCategoryResult{
			Categories: []domain.Category{},
			TotalCount: 0,
		}, nil
	}

	oids := make([]bson.ObjectID, len(ids))
	for _, id := range ids {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, domain.ErrInvalidCategoryID
		}
		oids = append(oids, oid)
	}

	filter := bson.M{
		"_id": bson.M{"$in": oids},
	}
	return r.findWithFilter(ctx, filter, nil)
}

func (r *CategoryRepo) FindChildren(ctx context.Context, parentID string, p pagination.Params) (*domain.ListCategoryResult, error) {
	poid, err := bson.ObjectIDFromHex(parentID)
	if err != nil {
		return nil, domain.ErrInvalidCategoryID
	}

	filter := bson.M{"parent_id": poid}
	return r.findWithFilter(ctx, filter, &p)
}

func (r *CategoryRepo) Update(ctx context.Context, c *domain.Category) (*domain.Category, error) {
	m, err := categoryModelFromDomain(c)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	m.UpdatedAt = now

	filter := bson.M{"_id": m.ID}
	update := bson.M{
		"$set": bson.M{
			"parent_id":  m.ParentID,
			"name_vi":    m.NameVi,
			"name_en":    m.NameEn,
			"slug":       m.Slug,
			"icon":       m.Icon,
			"sort_order": m.SortOrder,
			"updated_at": m.UpdatedAt,
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("updating category: %w", err)
	}
	if res.MatchedCount == 0 {
		return nil, domain.ErrCategoryNotFound
	}
	c.UpdatedAt = now
	return c, nil
}

func (r *CategoryRepo) Delete(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidCategoryID
	}

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrCategoryNotFound
	}
	return nil
}

func (r *CategoryRepo) List(ctx context.Context, p pagination.Params) (*domain.ListCategoryResult, error) {
	return r.findWithFilter(ctx, bson.M{}, &p)
}

func decodeCategoryList(ctx context.Context, cursor *mongo.Cursor, total int64) (*domain.ListCategoryResult, error) {
	var models []categoryModel
	if err := cursor.All(ctx, &models); err != nil {
		return &domain.ListCategoryResult{
			Categories: []domain.Category{},
			TotalCount: 0,
		}, fmt.Errorf("decoding category list: %w", err)
	}

	categories := make([]domain.Category, len(models))
	for i, m := range models {
		if c := m.toDomain(); c != nil {
			categories[i] = *c
		}
	}
	return &domain.ListCategoryResult{
		Categories: categories,
		TotalCount: total,
	}, nil
}
