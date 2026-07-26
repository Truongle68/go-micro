package usecase

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/repo"
	"context"
	"fmt"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type ProductUC struct {
	repo repo.ProductRepository
}

var _ ProductUsecase = (*ProductUC)(nil)

func NewProductUC(repo repo.ProductRepository) *ProductUC {
	return &ProductUC{
		repo: repo,
	}
}

func toProductDTO(p *domain.Product) *ProductDTO {
	if p == nil {
		return nil
	}

	variants := make([]*ProductVariantDTO, len(p.Variants))
	for i, v := range p.Variants {
		variants[i] = &ProductVariantDTO{
			ID:           v.ID,
			VariantLabel: v.VariantLabel,
			PriceDelta:   v.PriceDelta,
			Sku:          v.Sku,
		}
	}

	images := make([]*ProductImageDTO, len(p.Images))
	for i, img := range p.Images {
		images[i] = &ProductImageDTO{
			ID:        img.ID,
			Url:       img.Url,
			SortOrder: img.SortOrder,
		}
	}

	return &ProductDTO{
		ID:            p.ID,
		CategoryID:    p.CategoryID,
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
		Variants:      variants,
		Images:        images,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
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

func (uc *ProductUC) Create(ctx context.Context, in CreateProductInput) (*ProductDTO, error) {
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

	return toProductDTO(p), nil
}

func (uc *ProductUC) GetByID(ctx context.Context, id string) (*ProductDTO, error) {
	if id == "" {
		return nil, domain.ErrEmptyProductID
	}
	p, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding by id: %w", err)
	}
	return toProductDTO(p), nil
}

func (uc *ProductUC) GetByCategory(ctx context.Context, categoryID string) ([]*ProductDTO, error) {
	if categoryID == "" {
		return nil, domain.ErrEmptyCategoryID
	}
	products, err := uc.repo.FindByCategory(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("finding by category: %w", err)
	}

	dtos := make([]*ProductDTO, len(products))
	for i := range products {
		dtos[i] = toProductDTO(&products[i])
	}
	return dtos, nil
}

func (uc *ProductUC) Update(ctx context.Context, in UpdateProductInput) (*ProductDTO, error) {
	if in.ID == "" {
		return nil, domain.ErrEmptyProductID
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

	return toProductDTO(updated), nil
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

func (uc *ProductUC) List(ctx context.Context, p pagination.Params) ([]*ProductDTO, error) {
	products, err := uc.repo.List(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}

	dtos := make([]*ProductDTO, len(products))
	for i := range products {
		dtos[i] = toProductDTO(&products[i])
	}
	return dtos, nil
}

func (uc *ProductUC) Search(ctx context.Context, sParams domain.SearchProductParams, pParams pagination.Params) (*ProductSearchResultDTO, error) {
	if err := sParams.Validate(); err != nil {
		return nil, err
	}

	result, err := uc.repo.Search(ctx, sParams, pParams)
	if err != nil {
		return nil, fmt.Errorf("searching product: %w", err)
	}

	dtos := make([]*ProductDTO, len(result.Products))
	for i := range result.Products {
		dtos[i] = toProductDTO(&result.Products[i])
	}

	return &ProductSearchResultDTO{
		Products:   dtos,
		TotalCount: result.TotalCount,
	}, nil
}
