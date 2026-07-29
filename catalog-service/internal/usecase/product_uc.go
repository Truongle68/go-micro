package usecase

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/repo"
	"context"
	"fmt"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type ProductUC struct {
	repo     repo.ProductRepository
	cateRepo repo.CategoryRepository
}

var _ ProductUsecase = (*ProductUC)(nil)

func NewProductUC(repo repo.ProductRepository, cateRepo repo.CategoryRepository) *ProductUC {
	return &ProductUC{
		repo:     repo,
		cateRepo: cateRepo,
	}
}

func toDomainVariants(in []ProductVariantInput) []domain.ProductVariant {
	variants := make([]domain.ProductVariant, len(in))
	for i, v := range in {
		variants[i] = domain.ProductVariant{
			VariantLabel: v.VariantLabel,
			PriceDelta:   v.PriceDelta,
			Sku:          v.Sku,
		}
	}
	return variants
}

func toDomainImages(in []ProductImageInput) []domain.ProductImage {
	images := make([]domain.ProductImage, len(in))
	for i, img := range in {
		images[i] = domain.ProductImage{
			Url:       img.Url,
			SortOrder: img.SortOrder,
		}
	}
	return images
}

func (uc *ProductUC) Create(ctx context.Context, in CreateProductInput) (*domain.Product, error) {
	p, err := domain.NewProduct(domain.NewProductParams{
		CategoryID:    in.CategoryID,
		Sku:           in.Sku,
		NameVi:        in.NameVi,
		NameEn:        in.NameEn,
		DescriptionVi: in.DescriptionVi,
		DescriptionEn: in.DescriptionEn,
		Unit:          in.Unit,
		BasePrice:     in.BasePrice,
		SalePrice:     in.SalePrice,
		IsActive:      in.IsActive,
		Variants:      toDomainVariants(in.Variants),
		Images:        toDomainImages(in.Images),
	})

	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("creating product: %w", err)
	}

	return p, nil
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
	return in.CategoryID == nil &&
		in.Sku == nil &&
		in.NameVi == nil &&
		in.NameEn == nil &&
		in.DescriptionVi == nil &&
		in.DescriptionEn == nil &&
		in.Unit == nil &&
		in.BasePrice == nil &&
		in.SalePrice == nil &&
		in.IsActive == nil &&
		in.Variants == nil &&
		in.Images == nil
}

func (uc *ProductUC) Update(ctx context.Context, in UpdateProductInput) (*domain.Product, error) {
	if in.ID == "" {
		return nil, domain.ErrEmptyProductID
	}

	if in.isEmpty() {
		return nil, domain.ErrNoFieldsToUpdate
	}

	if in.CategoryID != nil && *in.CategoryID != "" {
		if _, err := uc.cateRepo.FindByID(ctx, *in.CategoryID); err != nil {
			return nil, fmt.Errorf("finding category by id: %w", err)
		}
	}

	p, err := uc.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, fmt.Errorf("finding by id: %w", err)
	}

	var variants []domain.ProductVariant
	if in.Variants != nil {
		variants = toDomainVariants(in.Variants)
	}

	var images []domain.ProductImage
	if in.Images != nil {
		images = toDomainImages(in.Images)
	}

	if err := p.ApplyUpdate(domain.UpdateProductParams{
		CategoryID:    in.CategoryID,
		Sku:           in.Sku,
		NameVi:        in.NameVi,
		NameEn:        in.NameEn,
		DescriptionVi: in.DescriptionVi,
		DescriptionEn: in.DescriptionEn,
		Unit:          in.Unit,
		BasePrice:     in.BasePrice,
		SalePrice:     in.SalePrice,
		IsActive:      in.IsActive,
		Variants:      variants,
		Images:        images,
	}); err != nil {
		return nil, err
	}

	updated, err := uc.repo.Update(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("updating product: %w", err)
	}

	return updated, nil
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

func (uc *ProductUC) List(ctx context.Context, p pagination.Params) (*ProductList, error) {
	result, err := uc.repo.List(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}

	return uc.loadDetailedProductList(ctx, result)
}

func (uc *ProductUC) Search(ctx context.Context, sParams domain.SearchProductParams, pParams pagination.Params) (*ProductList, error) {
	if err := sParams.Validate(); err != nil {
		return nil, err
	}

	result, err := uc.repo.Search(ctx, sParams, pParams)
	if err != nil {
		return nil, fmt.Errorf("searching product: %w", err)
	}

	return uc.loadDetailedProductList(ctx, result)
}
