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

const (
	defaultSearchLimit = 10
	maxSearchLimit     = 100
)

type ProductRepo struct {
	productColl *mongo.Collection
}

var _ repo.ProductRepository = (*ProductRepo)(nil)

func NewProductRepo(db *mongo.Database) *ProductRepo {
	return &ProductRepo{
		productColl: db.Collection("products"),
	}
}

type variantModel struct {
	ID           bson.ObjectID `bson:"_id"`
	VariantLabel string        `bson:"variant_label"`
	PriceDelta   float64       `bson:"price_delta"`
	Sku          string        `bson:"sku"`
}

type imageModel struct {
	ID        bson.ObjectID `bson:"_id"`
	Url       string        `bson:"url"`
	SortOrder int           `bson:"sort_order"`
}

type populatedCategoryModel struct {
	ID        string `bson:"id"`
	NameVi    string `bson:"name_vi"`
	NameEn    string `bson:"name_en"`
	Slug      string `bson:"slug"`
	Icon      string `bson:"icon"`
	SortOrder int64  `bson:"sort_order"`
}

type productModel struct {
	ID            bson.ObjectID  `bson:"_id"`
	CategoryID    bson.ObjectID  `bson:"category_id"`
	Sku           string         `bson:"sku"`
	NameVi        string         `bson:"name_vi"`
	NameEn        string         `bson:"name_en"`
	DescriptionVi string         `bson:"description_vi"`
	DescriptionEn string         `bson:"description_en"`
	Unit          string         `bson:"unit"`
	BasePrice     float64        `bson:"base_price"`
	SalePrice     float64        `bson:"sale_price"`
	RatingAvg     float64        `bson:"rating_avg"`
	RatingCount   int64          `bson:"rating_count"`
	IsActive      bool           `bson:"is_active"`
	Variants      []variantModel `bson:"variants"`
	Images        []imageModel   `bson:"images"`
	CreatedAt     time.Time      `bson:"created_at"`
	UpdatedAt     time.Time      `bson:"updated_at"`
}

func (m productModel) toDomain() *domain.Product {
	variants := make([]domain.ProductVariant, len(m.Variants))
	for i, v := range m.Variants {
		variants[i] = domain.ProductVariant{
			ID:           v.ID.Hex(),
			VariantLabel: v.VariantLabel,
			PriceDelta:   v.PriceDelta,
			Sku:          v.Sku,
		}
	}

	images := make([]domain.ProductImage, len(m.Images))
	for i, v := range m.Images {
		images[i] = domain.ProductImage{
			ID:        v.ID.Hex(),
			Url:       v.Url,
			SortOrder: int64(v.SortOrder),
		}
	}

	return &domain.Product{
		ID:            m.ID.Hex(),
		CategoryID:    m.CategoryID.Hex(),
		Sku:           m.Sku,
		NameVi:        m.NameVi,
		NameEn:        m.NameEn,
		DescriptionVi: m.DescriptionVi,
		DescriptionEn: m.DescriptionEn,
		Unit:          m.Unit,
		BasePrice:     m.BasePrice,
		SalePrice:     m.SalePrice,
		RatingAvg:     m.RatingAvg,
		RatingCount:   m.RatingCount,
		IsActive:      m.IsActive,
		Variants:      variants,
		Images:        images,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func productModelFromDomain(p *domain.Product) (*productModel, error) {
	if p == nil {
		return nil, domain.ErrProductNotFound
	}
	categoryOID, err := bson.ObjectIDFromHex(p.CategoryID)
	if err != nil {
		return nil, domain.ErrInvalidCategoryID
	}

	m := &productModel{
		CategoryID:    categoryOID,
		Sku:           p.Sku,
		NameVi:        p.NameVi,
		NameEn:        p.NameEn,
		DescriptionVi: p.DescriptionVi,
		DescriptionEn: p.DescriptionEn,
		Unit:          p.Unit,
		BasePrice:     p.BasePrice,
		SalePrice:     p.SalePrice,
		RatingAvg:     p.RatingAvg,
		RatingCount:   p.RatingCount,
		IsActive:      p.IsActive,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}

	if p.ID != "" {
		oid, err := bson.ObjectIDFromHex(p.ID)
		if err != nil {
			return nil, domain.ErrInvalidProductID
		}
		m.ID = oid
	}

	m.Variants = make([]variantModel, len(p.Variants))
	for i, v := range p.Variants {
		vm := variantModel{
			VariantLabel: v.VariantLabel,
			PriceDelta:   v.PriceDelta,
			Sku:          v.Sku,
		}

		if v.ID != "" {
			oid, err := bson.ObjectIDFromHex(v.ID)
			if err != nil {
				return nil, fmt.Errorf("invalid variant id: %w", err)
			}
			vm.ID = oid
		} else {
			vm.ID = bson.NewObjectID()
		}
		m.Variants[i] = vm
	}

	m.Images = make([]imageModel, len(p.Images))
	for i, img := range p.Images {
		im := imageModel{
			Url:       img.Url,
			SortOrder: int(img.SortOrder),
		}

		if img.ID != "" {
			oid, err := bson.ObjectIDFromHex(img.ID)
			if err != nil {
				return nil, fmt.Errorf("invalid image id: %w", err)
			}
			im.ID = oid
		} else {
			im.ID = bson.NewObjectID()
		}
		m.Images[i] = im
	}

	return m, nil
}

func idFromHex(id string) (bson.ObjectID, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, domain.ErrInvalidProductID
	}
	return oid, nil
}

func (r *ProductRepo) Create(ctx context.Context, p *domain.Product) error {
	m, err := productModelFromDomain(p)
	if err != nil {
		return err
	}

	if m.ID.IsZero() {
		m.ID = bson.NewObjectID()
	}

	now := time.Now()
	m.CreatedAt, m.UpdatedAt = now, now

	res, err := r.productColl.InsertOne(ctx, m)
	if err != nil {
		return fmt.Errorf("inserting product: %w", err)
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		p.ID = oid.Hex()
	}

	p.CreatedAt, p.UpdatedAt = now, now
	return nil
}

func categoryLookupStage() mongo.Pipeline {
	return mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "categories",
			"localField":   "category_id",
			"foreignField": "_id",
			"as":           "category",
		}}},

		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$category",
			"preserveNullAndEmptyArrays": true,
		}}},
	}
}

func appendCategoryLookup(pipeline mongo.Pipeline) mongo.Pipeline {
	return append(pipeline, categoryLookupStage()...)
}

func (r *ProductRepo) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	oid, err := idFromHex(id)
	if err != nil {
		return nil, err
	}

	var m productModel
	err = r.productColl.FindOne(ctx, bson.M{"_id": oid}).Decode(&m)

	if err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return m.toDomain(), nil
}

func (r *ProductRepo) findByFilter(ctx context.Context, filter bson.M, p pagination.Params) (*domain.ProductListResult, error) {
	total, err := r.productColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("counting products: %w", err)
	}

	if total == 0 {
		return &domain.ProductListResult{
			Products:   []domain.Product{},
			TotalCount: 0,
		}, nil
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetSkip(p.Skip()).
		SetLimit(p.Limit)

	cursor, err := r.productColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("finding products: %w", err)
	}
	defer cursor.Close(ctx)

	return decodeProductList(ctx, cursor, total)
}

func decodeProductList(ctx context.Context, cursor *mongo.Cursor, total int64) (*domain.ProductListResult, error) {
	var models []productModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, fmt.Errorf("decoding products: %w", err)
	}
	products := make([]domain.Product, len(models))
	for i, m := range models {
		if p := m.toDomain(); p != nil {
			products[i] = *p
		}
	}

	return &domain.ProductListResult{
		Products:   products,
		TotalCount: total,
	}, nil
}

func (r *ProductRepo) FindByCategory(ctx context.Context, categoryID string, p pagination.Params) (*domain.ProductListResult, error) {
	oid, err := bson.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, domain.ErrInvalidCategoryID
	}

	filter := bson.M{"category_id": oid}
	return r.findByFilter(ctx, filter, p)
}

func (r *ProductRepo) Update(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	m, err := productModelFromDomain(p)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	m.UpdatedAt = now

	filter := bson.M{"_id": m.ID}
	update := bson.M{
		"$set": bson.M{
			"category_id":    m.CategoryID,
			"sku":            m.Sku,
			"name_vi":        m.NameVi,
			"name_en":        m.NameEn,
			"description_vi": m.DescriptionVi,
			"description_en": m.DescriptionEn,
			"unit":           m.Unit,
			"base_price":     m.BasePrice,
			"sale_price":     m.SalePrice,
			"rating_avg":     m.RatingAvg,
			"rating_count":   m.RatingCount,
			"is_active":      m.IsActive,
			"variants":       m.Variants,
			"images":         m.Images,
			"updated_at":     m.UpdatedAt,
		},
	}
	res, err := r.productColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("updating product: %w", err)
	}
	if res.MatchedCount == 0 {
		return nil, domain.ErrProductNotFound
	}
	p.UpdatedAt = now

	return p, nil
}

func (r *ProductRepo) Delete(ctx context.Context, id string) error {
	oid, err := idFromHex(id)
	if err != nil {
		return err
	}
	res, err := r.productColl.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return fmt.Errorf("deleting product: %w", err)
	}
	if res.DeletedCount == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *ProductRepo) List(ctx context.Context, p pagination.Params) (*domain.ProductListResult, error) {
	return r.findByFilter(ctx, bson.M{}, p)
}

func (r *ProductRepo) Search(ctx context.Context, sParams domain.SearchProductParams, pParams pagination.Params) (*domain.ProductListResult, error) {

	filter, err := buildQueryFilter(sParams)
	if err != nil {
		return nil, err
	}

	if len(filter) == 0 {
		return r.List(ctx, pParams)
	}

	return r.findByFilter(ctx, filter, pParams)
}

func buildQueryFilter(params domain.SearchProductParams) (bson.M, error) {
	filter, err := buildFilter(params)
	if err != nil {
		return nil, err
	}

	if params.Query == "" {
		return filter, nil
	}

	search := buildSearchFilter(params.Query)
	if len(filter) == 0 {
		return search, nil
	}
	return bson.M{
		"$and": bson.A{filter, search},
	}, nil
}

func buildSearchFilter(k string) bson.M {
	regexQuery := bson.M{
		"$regex":   k,
		"$options": "i",
	}

	return bson.M{
		"$or": bson.A{
			bson.M{"name_vi": regexQuery},
			bson.M{"name_en": regexQuery},
			bson.M{"sku": regexQuery},
			bson.M{"description_vi": regexQuery},
			bson.M{"description_en": regexQuery},
		},
	}
}

func buildFilter(params domain.SearchProductParams) (bson.M, error) {
	filter := bson.M{}

	if params.CategoryID != "" {
		oid, err := bson.ObjectIDFromHex(params.CategoryID)
		if err != nil {
			return bson.M{}, domain.ErrInvalidCategoryID
		}
		filter["category_id"] = oid
	}

	if params.IsActive != nil {
		filter["is_active"] = *params.IsActive
	}

	if params.MinPrice != nil || params.MaxPrice != nil {
		price := bson.M{}
		if params.MinPrice != nil {
			price["$gte"] = *params.MinPrice
		}
		if params.MaxPrice != nil {
			price["$lte"] = *params.MaxPrice
		}
		filter["base_price"] = price
	}
	return filter, nil
}
