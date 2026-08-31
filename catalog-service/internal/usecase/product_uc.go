package usecase

import (
	"catalog-service/internal/domain"
	"catalog-service/pkg/sliceutil"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/gosimple/slug"
	"github.com/microcosm-cc/bluemonday"
)

type ProductUC struct {
	repo      ProductRepository
	cateRepo  CategoryRepository
	sanitizer *bluemonday.Policy
}

func NewProductUC(repo ProductRepository, cateRepo CategoryRepository) *ProductUC {
	return &ProductUC{
		repo:      repo,
		cateRepo:  cateRepo,
		sanitizer: bluemonday.UGCPolicy(),
	}
}

func (uc *ProductUC) Create(ctx context.Context, in CreateProductInput) (*domain.Product, error) {
	if err := validateVariantOptions(in.Variants, in.OptionTypes); err != nil {
		return nil, err
	}

	category, err := uc.cateRepo.FindByID(ctx, in.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("find category: %w", err)
	}
	if category == nil {
		return nil, domain.ErrCategoryNotFound
	}
	categoryPath, err := uc.cateRepo.BuildBreadcrumb(ctx, category.ID)
	if err != nil {
		return nil, fmt.Errorf("build category path: %w", err)
	}

	safeHTML := uc.sanitizer.Sanitize(in.DescriptionHTML)

	status := domain.ProductStatus(in.Status)
	if !status.IsValid() {
		return nil, domain.ErrInvalidProductStatus
	}

	images := in.Images
	if len(images) == 0 && len(in.Variants) > 0 {
		images = sliceutil.DedupeString(in.Variants, func(v CreateVariantInput) string {
			return v.Image
		})
	}

	now := time.Now()
	product := &domain.Product{
		Slug:            generateUniqueSlug(in.Name),
		Version:         1,
		Name:            in.Name,
		NameTranslation: in.NameTranslation,
		CategoryID:      category.ID,
		CategoryPath:    categoryPath,
		Tags:            in.Tags,
		Description:     in.Description,
		DescriptionHTML: safeHTML,
		Highlights:      in.Highlights,
		Images:          images,
		OptionTypes:     toDomainOptionTypes(in.OptionTypes),
		Variants:        toDomainVariants(in.Variants, now),
		Specifications:  toDomainSpecifications(in.Specifications),
		Status:          status,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	maxRetries := 3
	for retry := 0; retry <= maxRetries; retry++ {
		err = uc.repo.Create(ctx, product)
		if err == nil {
			return product, nil
		}
		if errors.Is(err, domain.ErrDuplicateSlug) && retry < maxRetries {
			product.Slug = fmt.Sprintf("%s-%d", slug.Make(in.Name), retry+1)
			continue
		}
		if errors.Is(err, domain.ErrDuplicateSKU) {
			return nil, domain.ErrDuplicateSKU
		}
		if errors.Is(err, domain.ErrDuplicateSlug) {
			return nil, domain.ErrDuplicateSlug
		}
		if domain.IsDuplicateKeyError(err) {
			return nil, err
		}
		return nil, fmt.Errorf("creating product: %w", err)
	}

	return product, nil
}

func validateVariantOptions(variants []CreateVariantInput, optionTypes []OptionTypeInput) error {
	if len(variants) == 0 {
		return fmt.Errorf("%w: product must have at least 1 variant", domain.ErrInvalidSimpleVariant)
	}
	if len(optionTypes) == 0 {
		if len(variants) != 1 {
			return fmt.Errorf("%w: product with no option types must have exactly 1 variant (got %d)", domain.ErrInvalidSimpleVariant, len(variants))
		}
		if len(variants[0].Attributes) > 0 {
			return fmt.Errorf("%w: product has no option types, but variant %s contains attributes", domain.ErrInvalidVariantAttribute, variants[0].SKU)
		}
		return nil
	}

	expectedVariantCount := 1
	for _, ot := range optionTypes {
		expectedVariantCount *= len(ot.Values)
	}

	if len(variants) > expectedVariantCount {
		return fmt.Errorf("%w: got %d, max possible is %d", domain.ErrExceedExpectedVariantCount, len(variants), expectedVariantCount)
	}

	return validateVariantAttributesMatchOptions(variants, optionTypes)
}

func validateVariantAttributesMatchOptions(variants []CreateVariantInput, optionTypes []OptionTypeInput) error {
	totalPairs := 0
	for _, ot := range optionTypes {
		totalPairs += len(ot.Values)
	}
	validPairs := make(map[string]struct{}, totalPairs)
	validKeys := make(map[string]struct{}, len(optionTypes))

	for _, ot := range optionTypes {
		key := strings.ToLower(strings.TrimSpace(ot.Name))
		if _, duplicate := validKeys[key]; duplicate {
			return fmt.Errorf("%w: duplicate option type name %q", domain.ErrInvalidOptionType, ot.Name)
		}
		validKeys[key] = struct{}{}
		for _, v := range ot.Values {
			pairKey := key + ":" + strings.ToLower(strings.TrimSpace(v))
			validPairs[pairKey] = struct{}{}
		}
	}
	seenSKUs := make(map[string]struct{}, len(variants))
	seenCombination := make(map[string]string, len(variants))
	for _, v := range variants {
		// check SKU uniqueness
		normalizedSKU := strings.ToUpper(strings.TrimSpace(v.SKU))
		if _, exists := seenSKUs[normalizedSKU]; exists {
			return fmt.Errorf("%w: %s", domain.ErrDuplicateSKU, v.SKU)
		}
		seenSKUs[normalizedSKU] = struct{}{}

		// check attribute count match
		if len(v.Attributes) != len(optionTypes) {
			return fmt.Errorf("%w: variants %s attributes count (%d) does not match option types count (%d)", domain.ErrInvalidVariantAttribute, v.SKU, len(v.Attributes), len(optionTypes))
		}

		// validate key/value pairs
		attrPairs := make([]string, 0, len(v.Attributes))
		for key, val := range v.Attributes {
			lowerKey := strings.ToLower(strings.TrimSpace(key))
			lowerVal := strings.ToLower(strings.TrimSpace(val))
			if _, exists := validKeys[lowerKey]; !exists {
				return fmt.Errorf("%w: variant %s has unknown attribute %q", domain.ErrInvalidVariantAttribute, v.SKU, key)
			}

			pairKey := lowerKey + ":" + lowerVal
			if _, valid := validPairs[pairKey]; !valid {
				return fmt.Errorf("%w: variant %s has invalid value %q for attribute %q", domain.ErrInvalidVariantAttribute, v.SKU, val, key)
			}

			attrPairs = append(attrPairs, pairKey)
		}

		sort.Strings(attrPairs)
		comboKey := strings.Join(attrPairs, "|")

		// check duplicate option combinations
		if existingSKU, duplicate := seenCombination[comboKey]; duplicate {
			return fmt.Errorf("%w: duplicate variant attribute combination between SKU %s and SKU %s",
				domain.ErrInvalidVariantAttribute, existingSKU, v.SKU)
		}
		seenCombination[comboKey] = v.SKU
	}
	return nil
}

func generateUniqueSlug(name string) string {
	base := slug.Make(name)
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%s-%s", base, hex.EncodeToString(b))
}

func toDomainVariants(in []CreateVariantInput, now time.Time) []domain.Variant {
	out := make([]domain.Variant, len(in))
	for i, v := range in {
		out[i] = domain.Variant{
			ID:         v.ID,
			SKU:        v.SKU,
			Attributes: v.Attributes,
			Price: domain.Price{
				Amount:   v.Price.Amount,
				Currency: v.Price.Currency,
			},
			Image:     v.Image,
			IsActive:  true,
			CreatedAt: now,
		}
	}
	return out
}

func toDomainOptionTypes(in []OptionTypeInput) []domain.OptionType {
	out := make([]domain.OptionType, len(in))
	for i, ot := range in {
		out[i] = domain.OptionType{
			Name:   ot.Name,
			Values: ot.Values,
		}
	}
	return out
}

func toDomainSpecifications(in []SpecGroupInput) []domain.SpecGroup {
	out := make([]domain.SpecGroup, len(in))
	for i, sg := range in {
		items := make([]domain.SpecItem, len(sg.Items))
		for j, item := range sg.Items {
			items[j] = domain.SpecItem{
				Label: item.Label,
				Value: item.Value,
			}
		}
		out[i] = domain.SpecGroup{
			Group: sg.Group,
			Items: items,
		}
	}
	return out
}

func (uc *ProductUC) GetByID(ctx context.Context, id string) (*DetailedProduct, error) {
	if id == "" {
		return nil, domain.ErrEmptyProductID
	}

	p, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding by id: %w", err)
	}

	c, err := uc.cateRepo.FindByID(ctx, p.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("finding by id: %w", err)
	}

	return &DetailedProduct{
		Product:  *p,
		Category: c,
	}, nil
}

func categoryIDInProducts(products []domain.Product) []string {
	seen := make(map[string]struct{}, len(products))
	ids := make([]string, 0, len(products))
	for _, p := range products {
		if p.CategoryID == "" {
			continue
		}
		if _, ok := seen[p.CategoryID]; ok {
			continue
		}
		seen[p.CategoryID] = struct{}{}
		ids = append(ids, p.CategoryID)
	}
	return ids
}

func toDetailedProducts(products []domain.Product, categories []domain.Category) []DetailedProduct {
	categoryByID := make(map[string]domain.Category, len(categories))
	for _, c := range categories {
		categoryByID[c.ID] = c
	}

	out := make([]DetailedProduct, len(products))
	for i, p := range products {
		dp := DetailedProduct{Product: p}
		if cat, ok := categoryByID[p.CategoryID]; ok {
			c := cat
			dp.Category = &c
		}
		out[i] = dp
	}
	return out
}

func (uc *ProductUC) loadDetailedProductList(
	ctx context.Context,
	productResult *domain.ProductListResult,
) (*ProductList, error) {
	categoryIDs := categoryIDInProducts(productResult.Products)

	var categories []domain.Category
	if len(categoryIDs) > 0 {
		categoryResult, err := uc.cateRepo.FindByIDs(ctx, categoryIDs)
		if err != nil {
			return nil, fmt.Errorf("finding categories by ids: %w", err)
		}
		categories = categoryResult.Categories
	}

	return &ProductList{
		Products:   toDetailedProducts(productResult.Products, categories),
		TotalCount: productResult.TotalCount,
	}, nil
}

func (uc *ProductUC) GetByCategory(ctx context.Context, categoryID string, p pagination.Params) (*ProductList, error) {
	if categoryID == "" {
		return nil, domain.ErrEmptyCategoryID
	}
	productResult, err := uc.repo.FindByCategory(ctx, categoryID, p)
	if err != nil {
		return nil, fmt.Errorf("finding by category: %w", err)
	}

	return uc.loadDetailedProductList(ctx, productResult)
}

func (in UpdateProductInput) isEmpty() bool {
	return in.Name == nil &&
		in.NameTranslation == nil &&
		in.CategoryID == nil &&
		in.Description == nil &&
		in.DescriptionHTML == nil &&
		in.Highlights == nil &&
		in.Tags == nil &&
		in.Images == nil &&
		in.OptionTypes == nil &&
		in.Variants == nil &&
		in.Specifications == nil &&
		in.Status == nil
}

func (uc *ProductUC) Update(ctx context.Context, in UpdateProductInput) (*domain.Product, error) {
	if in.ID == "" {
		return nil, domain.ErrEmptyProductID
	}

	if in.isEmpty() {
		return nil, domain.ErrNoFieldsToUpdate
	}

	if in.Version <= 0 {
		return nil, domain.ErrInvalidVersion
	}

	p, err := uc.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, fmt.Errorf("finding by id: %w", err)
	}

	var categoryPath []domain.CategoryRef
	if in.CategoryID != nil && *in.CategoryID != "" {
		category, err := uc.cateRepo.FindByID(ctx, *in.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("finding category by id: %w", err)
		}
		if category == nil {
			return nil, domain.ErrCategoryNotFound
		}
		path, err := uc.cateRepo.BuildBreadcrumb(ctx, category.ID)
		if err != nil {
			return nil, fmt.Errorf("building category path: %w", err)
		}
		categoryPath = path
	}

	if err := uc.validateVariantState(p, in); err != nil {
		return nil, err
	}

	var descHTML *string
	if in.DescriptionHTML != nil {
		safeHTML := uc.sanitizer.Sanitize(*in.DescriptionHTML)
		descHTML = &safeHTML
	}

	now := time.Now()
	var optionTypes []domain.OptionType
	if in.OptionTypes != nil {
		optionTypes = toDomainOptionTypes(in.OptionTypes)
	}
	var variants []domain.Variant
	if in.Variants != nil {
		variants = toDomainVariants(in.Variants, now)
	}
	var specs []domain.SpecGroup
	if in.Specifications != nil {
		specs = toDomainSpecifications(in.Specifications)
	}

	if err := p.ApplyUpdate(domain.UpdateProductParams{
		Name:            in.Name,
		NameTranslation: in.NameTranslation,
		CategoryID:      in.CategoryID,
		CategoryPath:    categoryPath,
		Description:     in.Description,
		DescriptionHTML: descHTML,
		Highlights:      in.Highlights,
		Tags:            in.Tags,
		OptionTypes:     optionTypes,
		Variants:        variants,
		Specifications:  specs,
		Status:          in.Status,
	}); err != nil {
		return nil, err
	}

	updated, err := uc.repo.Update(ctx, p, in.Version)
	if err != nil {
		if domain.IsDuplicateKeyError(err) {
			return nil, err
		}
		return nil, fmt.Errorf("updating product: %w", err)
	}

	return updated, nil
}

// Helper to isolate variant & option type validation logic
func (uc *ProductUC) validateVariantState(p *domain.Product, in UpdateProductInput) error {
	if in.Variants != nil {
		opts := in.OptionTypes
		if opts == nil {
			opts = make([]OptionTypeInput, 0, len(p.OptionTypes))
			for _, ot := range p.OptionTypes {
				opts = append(opts, OptionTypeInput{Name: ot.Name, Values: ot.Values})
			}
		}
		return validateVariantOptions(in.Variants, opts)
	}

	if in.OptionTypes != nil {
		existingVariants := make([]CreateVariantInput, len(p.Variants))
		for i, v := range p.Variants {
			existingVariants[i] = CreateVariantInput{
				SKU:        v.SKU,
				Attributes: v.Attributes,
			}
		}
		return validateVariantOptions(existingVariants, in.OptionTypes)
	}

	return nil
}

func (uc *ProductUC) Delete(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrEmptyProductID
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting product: %w", err)
	}
	return nil
}

func (uc *ProductUC) ListPublic(ctx context.Context, in PublicListInput) (*ProductList, error) {
	in.Query = strings.TrimSpace(in.Query)
	in.CategoryID = strings.TrimSpace(in.CategoryID)
	p := in.Page.Normalize()

	filter := applyPublicListPolicy(in)
	result, err := uc.repo.List(ctx, filter, p)
	if err != nil {
		return nil, fmt.Errorf("ProductUC.ListPublic: %w", err)
	}

	return uc.loadDetailedProductList(ctx, result)
}

func (uc *ProductUC) ListAdmin(ctx context.Context, in AdminListInput) (*ProductList, error) {
	in.Query = strings.TrimSpace(in.Query)
	in.CategoryID = strings.TrimSpace(in.CategoryID)
	in.SKU = strings.TrimSpace(in.SKU)
	p := in.Page.Normalize()

	filter := applyAdminListPolicy(in)
	result, err := uc.repo.List(ctx, filter, p)
	if err != nil {
		return nil, fmt.Errorf("ProductUC.ListAdmin: %w", err)
	}

	return uc.loadDetailedProductList(ctx, result)
}

func (uc *ProductUC) Search(ctx context.Context, sParams domain.SearchProductParams, pParams pagination.Params) (*ProductList, error) {
	if err := sParams.Validate(); err != nil {
		return nil, err
	}

	result, err := uc.repo.Search(ctx, sParams, pParams)
	if err != nil {
		return nil, fmt.Errorf("ProductUC.Search: %w", err)
	}

	return uc.loadDetailedProductList(ctx, result)
}

func (uc *ProductUC) GetVariantsBySKUs(ctx context.Context, skus []string) ([]VariantView, error) {
	products, err := uc.repo.FindByVariantSKUs(ctx, skus)
	if err != nil {
		return nil, err
	}

	want := make(map[string]struct{}, len(skus))
	for _, sku := range skus {
		want[sku] = struct{}{}
	}

	out := make([]VariantView, 0, len(skus))
	found := make(map[string]struct{})
	for _, p := range products {
		for _, v := range p.Variants {
			if _, ok := want[v.SKU]; !ok {
				continue
			}
			found[v.SKU] = struct{}{}
			out = append(out, VariantView{
				Variant:     v,
				ProductID:   p.ID,
				ProductName: p.Name,
			})
		}
	}

	for _, sku := range skus {
		if _, ok := found[sku]; !ok {
			return nil, fmt.Errorf("%w: %s", domain.ErrSKUNotFound, sku)
		}
	}

	return out, nil
}

func (uc *ProductUC) CountPerStatus(ctx context.Context, in AdminListInput) (*domain.ProductStatusCounts, error) {
	in.Query = strings.TrimSpace(in.Query)
	in.CategoryID = strings.TrimSpace(in.CategoryID)
	in.SKU = strings.TrimSpace(in.SKU)

	filter := applyAdminListPolicy(in)
	return uc.repo.CountPerStatus(ctx, filter)
}