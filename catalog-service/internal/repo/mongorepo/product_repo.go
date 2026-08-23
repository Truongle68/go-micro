package mongorepo

import (
	"catalog-service/internal/domain"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
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

func NewProductRepo(db *mongo.Database) *ProductRepo {
	return &ProductRepo{
		productColl: db.Collection("products"),
	}
}

type categoryRefModel struct {
	ID   bson.ObjectID `bson:"id"`
	Name string        `bson:"name"`
}

type optionTypeModel struct {
	Name   string   `bson:"name"`
	Values []string `bson:"values"`
}

type priceModel struct {
	Amount   int    `bson:"amount"`
	Currency string `bson:"currency"`
}

type variantModel struct {
	ID         bson.ObjectID     `bson:"_id,omitempty"`
	SKU        string            `bson:"sku"`
	Attributes map[string]string `bson:"attributes,omitempty"`
	Price      priceModel        `bson:"price"`
	Image      string            `bson:"images"`
	IsActive   bool              `bson:"is_active"`
	CreatedAt  time.Time         `bson:"created_at"`
}

type specItemModel struct {
	Label string `bson:"label"`
	Value string `bson:"value"`
}

type specGroupModel struct {
	Group string          `bson:"group"`
	Items []specItemModel `bson:"items"`
}

type productModel struct {
	ID              bson.ObjectID      `bson:"_id"`
	Version         int                `bson:"version"`
	Slug            string             `bson:"slug"`
	Name            string             `bson:"name"`
	NameTranslation map[string]string  `bson:"name_translation,omitempty"`
	CategoryID      bson.ObjectID      `bson:"category_id"`
	CategoryPath    []categoryRefModel `bson:"category_path,omitempty"`
	Description     string             `bson:"description"`
	DescriptionHTML string             `bson:"description_html,omitempty"`
	Highlights      []string           `bson:"highlights,omitempty"`
	Tags            []string           `bson:"tags,omitempty"`
	Images          []string           `bson:"images,omitempty"`
	OptionTypes     []optionTypeModel  `bson:"option_types,omitempty"`
	Variants        []variantModel     `bson:"variants,omitempty"`
	Specifications  []specGroupModel   `bson:"specifications,omitempty"`
	Status          string             `bson:"status"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
}

func (m productModel) toDomain() *domain.Product {
	categoryPath := make([]domain.CategoryRef, len(m.CategoryPath))
	for i, cp := range m.CategoryPath {
		categoryPath[i] = domain.CategoryRef{
			ID:   cp.ID.Hex(),
			Name: cp.Name,
		}
	}

	optionTypes := make([]domain.OptionType, len(m.OptionTypes))
	for i, ot := range m.OptionTypes {
		optionTypes[i] = domain.OptionType{
			Name:   ot.Name,
			Values: ot.Values,
		}
	}

	variants := make([]domain.Variant, len(m.Variants))
	for i, v := range m.Variants {
		var vID string
		if !v.ID.IsZero() {
			vID = v.ID.Hex()
		}

		variants[i] = domain.Variant{
			ID:         vID,
			SKU:        v.SKU,
			Attributes: v.Attributes,
			Price: domain.Price{
				Amount:   v.Price.Amount,
				Currency: v.Price.Currency,
			},
			Image:     v.Image,
			IsActive:  v.IsActive,
			CreatedAt: v.CreatedAt,
		}
	}

	specs := make([]domain.SpecGroup, len(m.Specifications))
	for i, sg := range m.Specifications {
		items := make([]domain.SpecItem, len(sg.Items))
		for j, item := range sg.Items {
			items[j] = domain.SpecItem{
				Label: item.Label,
				Value: item.Value,
			}
		}
		specs[i] = domain.SpecGroup{
			Group: sg.Group,
			Items: items,
		}
	}

	return &domain.Product{
		ID:              m.ID.Hex(),
		Version:         m.Version,
		Slug:            m.Slug,
		Name:            m.Name,
		NameTranslation: m.NameTranslation,
		CategoryID:      m.CategoryID.Hex(),
		CategoryPath:    categoryPath,
		Description:     m.Description,
		DescriptionHTML: m.DescriptionHTML,
		Highlights:      m.Highlights,
		Tags:            m.Tags,
		Images:          m.Images,
		OptionTypes:     optionTypes,
		Variants:        variants,
		Specifications:  specs,
		Status:          domain.ProductStatus(m.Status),
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func (vm variantModel) toDomain() domain.Variant {
	return domain.Variant{
		ID:         vm.ID.Hex(),
		SKU:        vm.SKU,
		Attributes: vm.Attributes,
		Price: domain.Price{
			Amount:   vm.Price.Amount,
			Currency: vm.Price.Currency,
		},
		Image:     vm.Image,
		IsActive:  vm.IsActive,
		CreatedAt: vm.CreatedAt,
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

	categoryPath := make([]categoryRefModel, len(p.CategoryPath))
	for i, cp := range p.CategoryPath {
		catOID, err := bson.ObjectIDFromHex(cp.ID)
		if err != nil {
			return nil, domain.ErrInvalidCategoryID
		}
		categoryPath[i] = categoryRefModel{
			ID:   catOID,
			Name: cp.Name,
		}
	}

	optionTypes := make([]optionTypeModel, len(p.OptionTypes))
	for i, ot := range p.OptionTypes {
		optionTypes[i] = optionTypeModel{
			Name:   ot.Name,
			Values: ot.Values,
		}
	}

	variants := make([]variantModel, len(p.Variants))
	for i, v := range p.Variants {
		var vID bson.ObjectID
		if v.ID != "" {
			if oid, err := bson.ObjectIDFromHex(v.ID); err == nil {
				vID = oid
			}
		}
		if vID.IsZero() {
			vID = bson.NewObjectID()
		}

		variants[i] = variantModel{
			ID:         vID,
			SKU:        v.SKU,
			Attributes: v.Attributes,
			Price: priceModel{
				Amount:   v.Price.Amount,
				Currency: v.Price.Currency,
			},
			Image:     v.Image,
			IsActive:  v.IsActive,
			CreatedAt: v.CreatedAt,
		}
	}

	specs := make([]specGroupModel, len(p.Specifications))
	for i, sg := range p.Specifications {
		items := make([]specItemModel, len(sg.Items))
		for j, item := range sg.Items {
			items[j] = specItemModel{
				Label: item.Label,
				Value: item.Value,
			}
		}
		specs[i] = specGroupModel{
			Group: sg.Group,
			Items: items,
		}
	}

	m := &productModel{
		Version:         p.Version,
		Slug:            p.Slug,
		Name:            p.Name,
		NameTranslation: p.NameTranslation,
		CategoryID:      categoryOID,
		CategoryPath:    categoryPath,
		Description:     p.Description,
		DescriptionHTML: p.DescriptionHTML,
		Highlights:      p.Highlights,
		Tags:            p.Tags,
		Images:          p.Images,
		OptionTypes:     optionTypes,
		Variants:        variants,
		Specifications:  specs,
		Status:          string(p.Status),
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}

	if p.ID != "" {
		oid, err := bson.ObjectIDFromHex(p.ID)
		if err != nil {
			return nil, domain.ErrInvalidProductID
		}
		m.ID = oid
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

func (r *ProductRepo) ExistSlug(ctx context.Context, name string) (bool, error) {
	count, err := r.productColl.CountDocuments(ctx, bson.M{"slug": name})
	if err != nil {
		return false, fmt.Errorf("checking slug existence: %w", err)
	}
	return count > 0, nil
}

func (r *ProductRepo) EnsureIndexes(ctx context.Context) error {
	models := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "category_id", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "variants.sku", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := r.productColl.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("creating product indexes: %w", err)
	}
	return nil
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
		if field, val, ok := extractDuplicateField(err); ok {
			if strings.Contains(field, "sku") {
				return fmt.Errorf("%w: %s = %v", domain.ErrDuplicateSKU, field, val)
			}
			if field == "slug" {
				return fmt.Errorf("%w: %s = %v", domain.ErrDuplicateSlug, field, val)
			}
			return fmt.Errorf("%w: %s = %v", domain.ErrDuplicateField, field, val)
		}
		return fmt.Errorf("inserting product: %w", err)
	}

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
	err = r.productColl.FindOne(ctx, bson.M{"_id": oid}).Decode(&m)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
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
		SetSort(bson.D{{Key: "_id", Value: -1}}).
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

func (r *ProductRepo) Update(ctx context.Context, p *domain.Product, expectedVersion int) (*domain.Product, error) {
	m, err := productModelFromDomain(p)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	m.UpdatedAt = now

	filter := bson.M{"_id": m.ID, "version": expectedVersion}
	update := bson.M{
		"$set": bson.M{
			"slug":             m.Slug,
			"name":             m.Name,
			"name_translation": m.NameTranslation,
			"category_id":      m.CategoryID,
			"category_path":    m.CategoryPath,
			"description":      m.Description,
			"description_html": m.DescriptionHTML,
			"highlights":       m.Highlights,
			"tags":             m.Tags,
			"images":           m.Images,
			"option_types":     m.OptionTypes,
			"variants":         m.Variants,
			"specifications":   m.Specifications,
			"status":           m.Status,
			"updated_at":       m.UpdatedAt,
		},
		"$inc": bson.M{
			"version": 1,
		},
	}
	res, err := r.productColl.UpdateOne(ctx, filter, update)
	if err != nil {
		if field, val, ok := extractDuplicateField(err); ok {
			if strings.Contains(field, "sku") {
				return nil, fmt.Errorf("%w: %s = %v", domain.ErrDuplicateSKU, field, val)
			}
			if field == "slug" {
				return nil, fmt.Errorf("%w: %s = %v", domain.ErrDuplicateSlug, field, val)
			}
			return nil, fmt.Errorf("%w: %s = %v", domain.ErrDuplicateField, field, val)
		}
		return nil, fmt.Errorf("updating product: %w", err)
	}
	if res.MatchedCount == 0 {
		count, _ := r.productColl.CountDocuments(ctx, bson.M{"_id": m.ID})
		if count == 0 {
			return nil, domain.ErrProductNotFound
		}
		return nil, domain.ErrConcurrentUpdate
	}
	p.Version++
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

func (r *ProductRepo) FindVariantsBySKUs(ctx context.Context, skus []string) ([]domain.Variant, error) {
	filter := bson.M{
		"variants.sku": bson.M{
			"$in": skus,
		},
	}
	total, err := r.productColl.CountDocuments(ctx, filter)
	if err != nil {
		return []domain.Variant{}, fmt.Errorf("counting variants: %w", err)
	}

	if total == 0 {
		return []domain.Variant{}, nil
	}

	if total < int64(len(skus)) {
		return []domain.Variant{}, fmt.Errorf("some not found: %w", err)
	}

	cursor, err := r.productColl.Find(ctx, filter)
	if err != nil {
		return []domain.Variant{}, fmt.Errorf("finding variants by skus: %w", err)
	}
	defer cursor.Close(ctx)

	var models []variantModel
	if err := cursor.All(ctx, &models); err != nil {
		return []domain.Variant{}, fmt.Errorf("decoding variants: %w", err)
	}

	results := make([]domain.Variant, len(models))
	for i, m := range models {
		results[i] = m.toDomain()
	}

	return results, nil
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
	escaped := regexp.QuoteMeta(k)
	regexQuery := bson.M{
		"$regex":   escaped,
		"$options": "i",
	}

	return bson.M{
		"$or": bson.A{
			bson.M{"name": regexQuery},
			bson.M{"tags": regexQuery},
			bson.M{"description": regexQuery},
			bson.M{"variants.sku": regexQuery},
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

	if params.Status != nil {
		filter["status"] = string(*params.Status)
	}

	if params.MinPrice != nil || params.MaxPrice != nil {
		price := bson.M{}
		if params.MinPrice != nil {
			price["$gte"] = *params.MinPrice
		}
		if params.MaxPrice != nil {
			price["$lte"] = *params.MaxPrice
		}
		filter["variants.price.amount"] = price
	}
	return filter, nil
}
