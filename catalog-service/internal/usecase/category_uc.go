package usecase

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/repo"
	"context"
	"fmt"
	"time"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type CategoryDTO struct {
	ID              string               `json:"id"`
	ParentID        *string              `json:"parent_id"`
	Name            string               `json:"name"`
	NameTranslation map[string]string    `json:"name_translation"`
	Slug            string               `json:"slug"`
	Icon            string               `json:"icon"`
	SortOrder       int64                `json:"sort_order"`
	IsActive        bool                 `json:"is_active"`
	Ancestors       []domain.CategoryRef `json:"ancestors"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

type CreateCategoryInput struct {
	ParentID        *string
	Name            string
	NameTranslation map[string]string
	Slug            string
	Icon            string
	SortOrder       int64
	IsActive        *bool
}

type UpdateCategoryInput struct {
	ID              string
	ParentID        *string
	Name            *string
	NameTranslation map[string]string
	Slug            *string
	Icon            *string
	SortOrder       *int64
	IsActive        *bool
}

type CategoryList struct {
	Categories []CategoryDTO `json:"categories"`
	TotalCount int64         `json:"total_count"`
}

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
		ID:              c.ID,
		ParentID:        c.ParentID,
		Name:            c.Name,
		NameTranslation: c.NameTranslation,
		Slug:            c.Slug,
		Icon:            c.Icon,
		SortOrder:       c.SortOrder,
		IsActive:        c.IsActive,
		Ancestors:       c.Ancestors,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
	}
}

func (uc *CategoryUC) buildAncestors(ctx context.Context, parentID *string) ([]domain.CategoryRef, error) {
	if parentID == nil || *parentID == "" {
		return nil, nil
	}
	parent, err := uc.repo.FindByID(ctx, *parentID)
	if err != nil {
		return nil, err
	}
	ancestors := make([]domain.CategoryRef, 0, len(parent.Ancestors)+1)
	ancestors = append(ancestors, parent.Ancestors...)
	ancestors = append(ancestors, domain.CategoryRef{
		ID:   parent.ID,
		Name: parent.Name,
	})
	return ancestors, nil
}

func (uc *CategoryUC) Create(ctx context.Context, in CreateCategoryInput) (*CategoryDTO, error) {
	ancestors, err := uc.buildAncestors(ctx, in.ParentID)
	if err != nil {
		return nil, fmt.Errorf("building ancestors: %w", err)
	}

	c, err := domain.NewCategory(domain.NewCategoryParams{
		ParentID:        in.ParentID,
		Name:            in.Name,
		NameTranslation: in.NameTranslation,
		Slug:            in.Slug,
		Icon:            in.Icon,
		SortOrder:       in.SortOrder,
		IsActive:        in.IsActive,
		Ancestors:       ancestors,
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

func (uc *CategoryUC) GetChildren(ctx context.Context, parentID string, p pagination.Params) (*CategoryList, error) {
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

	return &CategoryList{
		Categories: dtos,
		TotalCount: result.TotalCount,
	}, nil
}

func (in UpdateCategoryInput) isEmpty() bool {
	return in.ParentID == nil &&
		in.Name == nil &&
		in.NameTranslation == nil &&
		in.Slug == nil &&
		in.Icon == nil &&
		in.SortOrder == nil &&
		in.IsActive == nil
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

	var newAncestors []domain.CategoryRef
	if in.ParentID != nil {
		if *in.ParentID != "" {
			if *in.ParentID == in.ID {
				return nil, domain.ErrCircularCategoryParent
			}
			parent, err := uc.repo.FindByID(ctx, *in.ParentID)
			if err != nil {
				return nil, fmt.Errorf("finding parent category: %w", err)
			}
			for _, a := range parent.Ancestors {
				if a.ID == in.ID {
					return nil, domain.ErrCircularCategoryParent
				}
			}
			ancestors, err := uc.buildAncestors(ctx, in.ParentID)
			if err != nil {
				return nil, fmt.Errorf("building ancestors: %w", err)
			}
			newAncestors = ancestors
		}
	}

	c.ApplyUpdate(domain.UpdateCategoryParams{
		ParentID:        in.ParentID,
		Name:            in.Name,
		NameTranslation: in.NameTranslation,
		Slug:            in.Slug,
		Icon:            in.Icon,
		SortOrder:       in.SortOrder,
		IsActive:        in.IsActive,
		Ancestors:       newAncestors,
	})

	updated, err := uc.repo.Update(ctx, c)
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

func (uc *CategoryUC) List(ctx context.Context, p pagination.Params) (*CategoryList, error) {
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
	return &CategoryList{
		Categories: dtos,
		TotalCount: result.TotalCount,
	}, nil
}
