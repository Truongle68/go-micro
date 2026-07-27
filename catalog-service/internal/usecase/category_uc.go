package usecase

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/repo"
	"context"
	"fmt"

	"github.com/TruongLe68/go-micro/pkg/pagination"
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
	if in.ParentID != nil && *in.ParentID != "" {
		if _, err := uc.repo.FindByID(ctx, *in.ParentID); err != nil {
			return nil, err
		}
	}
	c, err := domain.NewCategory(domain.NewCategoryParams{
		ParentID:  in.ParentID,
		NameVi:    in.NameVi,
		NameEn:    in.NameVi,
		Slug:      in.Slug,
		Icon:      in.Icon,
		SortOrder: in.SortOrder,
	})

	if err != nil {
		return nil, err
	}

	if err = uc.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("creating category: %w", err)
	}

	return toCategoryDTO(c), nil
}

func (uc *CategoryUC) GetByID(ctx context.Context, id string) (*CategoryDTO, error) {
	if id == "" {
		return nil, domain.ErrEmptyCategoryID
	}
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding by id: %w", err)
	}
	return toCategoryDTO(c), nil
}

func (uc *CategoryUC) GetChildren(ctx context.Context, parentID string, p pagination.Params) (*CategoryListResultDTO, error) {
	result, err := uc.repo.FindChildren(ctx, parentID, p)
	if err != nil {
		return nil, fmt.Errorf("finding children: %w", err)
	}

	dtos := make([]CategoryDTO, len(result.Categories))
	for i, c := range result.Categories {
		dto := toCategoryDTO(&c)
		if dto != nil {
			dtos[i] = *dto
		}
	}

	return &CategoryListResultDTO{
		Categories: dtos,
		TotalCount: result.TotalCount,
	}, nil
}

func (in UpdateCategoryInput) isEmpty() bool {
	return in.ParentID == nil &&
		in.NameVi == nil &&
		in.NameEn == nil &&
		in.Slug == nil &&
		in.Icon == nil &&
		in.SortOrder == nil
}

func (uc *CategoryUC) Update(ctx context.Context, in UpdateCategoryInput) (*CategoryDTO, error) {
	if in.ID == "" {
		return nil, domain.ErrEmptyCategoryID
	}

	if in.isEmpty() {
		return nil, domain.ErrNoFieldsToUpdate
	}

	c, err := uc.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, fmt.Errorf("finding by id: %w", err)
	}

	c.ApplyUpdate(domain.UpdateCategoryParams{
		ParentID:  in.ParentID,
		NameVi:    in.NameVi,
		NameEn:    in.NameEn,
		Slug:      in.Slug,
		Icon:      in.Icon,
		SortOrder: in.SortOrder,
	})

	updated, err := uc.repo.Update(ctx, in.ID, c)
	if err != nil {
		return nil, fmt.Errorf("updating category: %w", err)
	}

	return toCategoryDTO(updated), nil
}

func (uc *CategoryUC) Delete(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrEmptyCategoryID
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *CategoryUC) List(ctx context.Context, p pagination.Params) (*CategoryListResultDTO, error) {
	result, err := uc.repo.List(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("listing categories: %w", err)
	}

	dtos := make([]CategoryDTO, len(result.Categories))
	for i, c := range result.Categories {
		dto := toCategoryDTO(&c)
		if dto != nil {
			dtos[i] = *dto
		}
	}
	return &CategoryListResultDTO{
		Categories: dtos,
		TotalCount: result.TotalCount,
	}, nil
}
