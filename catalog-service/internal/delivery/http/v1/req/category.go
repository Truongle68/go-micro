package req

import "catalog-service/internal/usecase"

type CreateCategory struct {
	ParentID        *string           `json:"parent_id"`
	Name            string            `json:"name" validate:"required"`
	NameTranslation map[string]string `json:"name_translation,omitempty"`
	Slug            string            `json:"slug" validate:"required"`
	Icon            string            `json:"icon"`
	SortOrder       int64             `json:"sort_order"`
	IsActive        *bool             `json:"is_active,omitempty"`
}

func (req *CreateCategory) ToCreateCategoryInput() usecase.CreateCategoryInput {
	return usecase.CreateCategoryInput{
		ParentID:        req.ParentID,
		Name:            req.Name,
		NameTranslation: req.NameTranslation,
		Slug:            req.Slug,
		Icon:            req.Icon,
		SortOrder:       req.SortOrder,
		IsActive:        req.IsActive,
	}
}

type UpdateCategory struct {
	ParentID        *string           `json:"parent_id"`
	Name            *string           `json:"name,omitempty"`
	NameTranslation map[string]string `json:"name_translation,omitempty"`
	Slug            *string           `json:"slug,omitempty"`
	Icon            *string           `json:"icon,omitempty"`
	SortOrder       *int64            `json:"sort_order,omitempty"`
	IsActive        *bool             `json:"is_active,omitempty"`
}

func (req *UpdateCategory) ToUpdateCategoryInput(id string) usecase.UpdateCategoryInput {
	return usecase.UpdateCategoryInput{
		ID:              id,
		ParentID:        req.ParentID,
		Name:            req.Name,
		NameTranslation: req.NameTranslation,
		Slug:            req.Slug,
		Icon:            req.Icon,
		SortOrder:       req.SortOrder,
		IsActive:        req.IsActive,
	}
}
