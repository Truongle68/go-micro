package usecase

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/repo"
	"context"
)

type CategoryUC struct {
	repo repo.CategoryRepository
}

var _ CategoryUsecase = (*CategoryUC)(nil)

func NewCategoryUC(repo repo.CategoryRepository) *CategoryUC {
	return &CategoryUC{
		repo: repo,
	}
}

func toCategoryDTO(c *domain.Category) *CategoryDTO {
	if c == nil {
		return nil
	}
	return &CategoryDTO{
		ID:        c.ID,
		ParentID:  c.ParentID,
		NameVi:    c.NameVi,
		NameEn:    c.NameEn,
		Slug:      c.Slug,
		Icon:      c.Icon,
		SortOrder: c.SortOrder,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (uc *CategoryUC) Create(ctx context.Context, in CreateCategoryInput) (*CategoryDTO, error) {
	if in.NameVi == "" && in.NameEn == "" {
		return nil, domain.ErrEmptyName
	}

	c := &domain.Category{
		ParentID:  in.ParentID,
		NameVi:    in.NameVi,
		NameEn:    in.NameEn,
		Slug:      in.Slug,
		Icon:      in.Icon,
		SortOrder: in.SortOrder,
	}

	err := uc.repo.Create(ctx, c)
	if err != nil {
		return nil, err
	}

	return toCategoryDTO(c), nil
}

func (uc *CategoryUC) GetByID(ctx context.Context, id string) (*CategoryDTO, error) {
	if id == "" {
		return nil, domain.ErrEmptyCategoryID
	}
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toCategoryDTO(c), nil
}

func (uc *CategoryUC) GetChildren(ctx context.Context, parentID string) ([]*CategoryDTO, error) {
	categories, err := uc.repo.FindChildren(ctx, parentID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*CategoryDTO, len(categories))
	for i := range categories {
		dtos[i] = toCategoryDTO(&categories[i])
	}
	return dtos, nil
}

func (uc *CategoryUC) Update(ctx context.Context, in UpdateCategoryInput) (*CategoryDTO, error) {
	if in.ID == "" {
		return nil, domain.ErrEmptyCategoryID
	}

	c, err := uc.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if in.ParentID != nil {
		c.ParentID = in.ParentID
	}
	if in.NameVi != nil {
		c.NameVi = *in.NameVi
	}
	if in.NameEn != nil {
		c.NameEn = *in.NameEn
	}
	if in.Slug != nil {
		c.Slug = *in.Slug
	}
	if in.Icon != nil {
		c.Icon = *in.Icon
	}
	if in.SortOrder != nil {
		c.SortOrder = *in.SortOrder
	}

	updated, err := uc.repo.Update(ctx, in.ID, c)
	if err != nil {
		return nil, err
	}

	return toCategoryDTO(updated), nil
}

func (uc *CategoryUC) Delete(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrEmptyCategoryID
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *CategoryUC) List(ctx context.Context) ([]*CategoryDTO, error) {
	categories, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]*CategoryDTO, len(categories))
	for i := range categories {
		dtos[i] = toCategoryDTO(&categories[i])
	}
	return dtos, nil
}
