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
	collection *mongo.Collection
}

var _ repo.ProductRepository = (*ProductRepo)(nil)

func NewProductRepo(db *mongo.Database) *ProductRepo {
	return &ProductRepo{
		collection: db.Collection("products"),
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

	res, err := r.collection.InsertOne(ctx, m)
	if err != nil {
		return fmt.Errorf("inserting product: %w", err)
	}

	// If MongoDB generated the ID, assign it back to struct
	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		p.ID = oid.Hex()
	}

	p.CreatedAt, p.UpdatedAt = now, now
	return nil
}

func (r *ProductRepo) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	oid, err := idFromHex(id)
	if err != nil {
		return nil, err
	}

	var m productModel
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&m)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrProductNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding product by id: %w", err)
	}
	return m.toDomain(), nil
}

func (r *ProductRepo) FindByCategory(ctx context.Context, categoryID string, p pagination.Params) (*domain.ListProductResult, error) {
	oid, err := bson.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, domain.ErrInvalidCategoryID
	}

	filter := bson.M{"category_id": oid}
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count products by category: %w", err)
	}

	if total == 0 {
		return &domain.ListProductResult{
			Products:   []domain.Product{},
			TotalCount: 0,
		}, nil
	}

	findOpts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetSkip(p.Skip()).
		SetLimit(p.Limit)

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, fmt.Errorf("finding products by category: %w", err)
	}
	defer cursor.Close(ctx)

	var models []productModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, fmt.Errorf("decoding products by category: %w", err)
	}
	products := make([]domain.Product, len(models))
	for i, m := range models {
		if p := m.toDomain(); p != nil {
			products[i] = *p
		}
	}
	return &domain.ListProductResult{
		Products:   products,
		TotalCount: total,
	}, nil
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
	res, err := r.collection.UpdateOne(ctx, filter, update)
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
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return fmt.Errorf("deleting product: %w", err)
	}
	if res.DeletedCount == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *ProductRepo) List(ctx context.Context, p pagination.Params) (*domain.ListProductResult, error) {
	filter := bson.M{}
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	if total == 0 {
		return &domain.ListProductResult{
			Products:   []domain.Product{},
			TotalCount: 0,
		}, nil
	}

	findOpts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetSkip(p.Skip()).
		SetLimit(p.Limit)

	cursor, err := r.collection.Find(ctx, bson.M{}, findOpts)
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}
	defer cursor.Close(ctx)

	var models []productModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, fmt.Errorf("decoding product list: %w", err)
	}

	products := make([]domain.Product, len(models))
	for i, m := range models {
		if p := m.toDomain(); p != nil {
			products[i] = *p
		}
	}

	return &domain.ListProductResult{
		Products:   products,
		TotalCount: total,
	}, nil
}

func (r *ProductRepo) Search(ctx context.Context, sParams domain.SearchProductParams, pParams pagination.Params) (*domain.ListProductResult, error) {

	filterClauses, err := buildFilterClauses(sParams)
	if err != nil {
		return nil, err
	}

	if sParams.Query == "" && len(filterClauses) == 0 {
		return r.List(ctx, pParams)
	}

	pipeline := buildSearchPipeline(sParams, filterClauses, pParams.Skip(), pParams.Limit)

	opts := options.Aggregate().SetAllowDiskUse(true)
	cursor, err := r.collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, fmt.Errorf("aggregating product search: %w", err)
	}
	defer cursor.Close(ctx)

	return decodeFacetResult(ctx, cursor)
}

func buildFilterClauses(params domain.SearchProductParams) ([]bson.M, error) {
	var clauses []bson.M

	if params.CategoryID != "" {
		oid, err := bson.ObjectIDFromHex(params.CategoryID)
		if err != nil {
			return []bson.M{}, domain.ErrInvalidCategoryID
		}
		clauses = append(clauses, bson.M{
			"equals": bson.M{
				"path":  "category_id",
				"value": oid,
			},
		})
	}

	if params.IsActive != nil {
		clauses = append(clauses, bson.M{
			"equals": bson.M{
				"path":  "is_active",
				"value": *params.IsActive,
			},
		})
	}

	if params.MinPrice != nil || params.MaxPrice != nil {
		rangeFilter := bson.M{"path": "base_price"}
		if params.MinPrice != nil {
			rangeFilter["gte"] = *params.MinPrice
		}
		if params.MaxPrice != nil {
			rangeFilter["lte"] = *params.MaxPrice
		}
		clauses = append(clauses, bson.M{
			"range": rangeFilter,
		})
	}
	return clauses, nil
}

func buildSearchPipeline(sParams domain.SearchProductParams, filterClauses []bson.M, skip, limit int64) mongo.Pipeline {
	compoundStage := bson.M{}

	if len(filterClauses) > 0 {
		compoundStage["filter"] = filterClauses
	}

	if sParams.Query != "" {
		compoundStage["should"] = []bson.M{
			{
				"autocomplete": bson.M{
					"query": sParams.Query,
					"path":  "name_en",
				},
			},
			{
				"autocomplete": bson.M{
					"query": sParams.Query,
					"path":  "name_vi",
				},
			},
			{
				"text": bson.M{
					"query": sParams.Query,
					"path":  []string{"name_vi", "name_en", "description_vi", "description_en", "sku"},
					"fuzzy": bson.M{"maxEdits": 1, "prefixLength": 1},
				},
			},
		}
		compoundStage["minimumShouldMatch"] = 1
	}

	searchStage := bson.M{
		"index":    "products_fts_idx",
		"compound": compoundStage,
	}

	addScoreStage := bson.D{
		{Key: "$addFields", Value: bson.M{
			"score": bson.M{"$meta": "searchScore"},
		}},
	}

	sortStage := bson.D{
		{Key: "score", Value: -1},
		{Key: "_id", Value: 1},
	}

	facetStage := bson.D{
		{Key: "results", Value: bson.A{
			bson.D{{Key: "$sort", Value: sortStage}},
			bson.D{{Key: "$skip", Value: skip}},
			bson.D{{Key: "$limit", Value: limit}},
		}},
		{Key: "totalCount", Value: bson.A{
			bson.D{{Key: "$count", Value: "count"}},
		}},
	}

	return mongo.Pipeline{
		bson.D{{Key: "$search", Value: searchStage}},
		addScoreStage,
		bson.D{{Key: "$facet", Value: facetStage}},
	}
}

type facetResult struct {
	Results    []productModel `bson:"results"`
	TotalCount []struct {
		Count int64 `bson:"count"`
	} `bson:"totalCount"`
}

func decodeFacetResult(ctx context.Context, cursor *mongo.Cursor) (*domain.ListProductResult, error) {
	var rawResults []facetResult
	if err := cursor.All(ctx, &rawResults); err != nil {
		return nil, fmt.Errorf("decoding search results: %w", err)
	}

	res := &domain.ListProductResult{
		Products:   []domain.Product{},
		TotalCount: 0,
	}

	if len(rawResults) == 0 {
		return res, nil
	}

	products := make([]domain.Product, len(rawResults[0].Results))
	for i, m := range rawResults[0].Results {
		if p := m.toDomain(); p != nil {
			products[i] = *p
		}
	}
	res.Products = products

	if len(rawResults[0].TotalCount) > 0 {
		res.TotalCount = rawResults[0].TotalCount[0].Count
	}
	return res, nil
}
