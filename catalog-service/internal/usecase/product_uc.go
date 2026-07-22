package usecase

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/repo"
	"context"
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

func (uc *ProductUC) Create(ctx context.Context, in CreateProductInput) (*ProductDTO, error) {
	if in.CategoryID == "" {
		return nil, domain.ErrEmptyCategoryID
	}
	if in.NameVi == "" && in.NameEn == "" {
		return nil, domain.ErrEmptyName
	}
	if in.Sku == "" {
		return nil, domain.ErrEmptySku
	}
	if in.BasePrice < 0 || in.SalePrice < 0 {
		return nil, domain.ErrInvalidPrice
	}

	variants := make([]domain.ProductVariant, len(in.Variants))
	for i, v := range in.Variants {
		variants[i] = domain.ProductVariant{
			VariantLabel: v.VariantLabel,
			PriceDelta:   v.PriceDelta,
			Sku:          v.Sku,
		}
	}

	images := make([]domain.ProductImage, len(in.Images))
	for i, img := range in.Images {
		images[i] = domain.ProductImage{
			Url:       img.Url,
			SortOrder: img.SortOrder,
		}
	}

	p := &domain.Product{
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
	}

	err := uc.repo.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	return toProductDTO(p), nil
}

func (uc *ProductUC) GetByID(ctx context.Context, id string) (*ProductDTO, error) {
	if id == "" {
		return nil, domain.ErrEmptyProductID
	}
	p, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toProductDTO(p), nil
}

func (uc *ProductUC) GetByCategory(ctx context.Context, categoryID string) ([]*ProductDTO, error) {
	products, err := uc.repo.FindByCategory(ctx, categoryID)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	if in.CategoryID != nil {
		p.CategoryID = *in.CategoryID
	}
	if in.Sku != nil {
		p.Sku = *in.Sku
	}
	if in.NameVi != nil {
		p.NameVi = *in.NameVi
	}
	if in.NameEn != nil {
		p.NameEn = *in.NameEn
	}
	if in.DescriptionVi != nil {
		p.DescriptionVi = *in.DescriptionVi
	}
	if in.DescriptionEn != nil {
		p.DescriptionEn = *in.DescriptionEn
	}
	if in.Unit != nil {
		p.Unit = *in.Unit
	}
	if in.BasePrice != nil {
		p.BasePrice = *in.BasePrice
	}
	if in.SalePrice != nil {
		p.SalePrice = *in.SalePrice
	}
	if in.IsActive != nil {
		p.IsActive = *in.IsActive
	}

	if in.Variants != nil {
		variants := make([]domain.ProductVariant, len(in.Variants))
		for i, v := range in.Variants {
			variants[i] = domain.ProductVariant{
				VariantLabel: v.VariantLabel,
				PriceDelta:   v.PriceDelta,
				Sku:          v.Sku,
			}
		}
		p.Variants = variants
	}

	if in.Images != nil {
		images := make([]domain.ProductImage, len(in.Images))
		for i, img := range in.Images {
			images[i] = domain.ProductImage{
				Url:       img.Url,
				SortOrder: img.SortOrder,
			}
		}
		p.Images = images
	}

	updated, err := uc.repo.Update(ctx, in.ID, p)
	if err != nil {
		return nil, err
	}

	return toProductDTO(updated), nil
}

func (uc *ProductUC) Delete(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrEmptyProductID
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *ProductUC) List(ctx context.Context) ([]*ProductDTO, error) {
	products, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]*ProductDTO, len(products))
	for i := range products {
		dtos[i] = toProductDTO(&products[i])
	}
	return dtos, nil
}
