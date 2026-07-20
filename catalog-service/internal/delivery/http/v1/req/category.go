package req

import "catalog-service/internal/usecase"

type CreateCategory struct {
	ParentID  *string `json:"parent_id"`
	NameVi    string  `json:"name_vi" validate:"required"`
	NameEn    string  `json:"name_en" validate:"required"`
	Slug      string  `json:"slug" validate:"required"`
	Icon      string  `json:"icon"`
	SortOrder int     `json:"sort_order"`
}

func (req *CreateCategory) ToCreateCategoryInput() usecase.CreateCategoryInput {
	return usecase.CreateCategoryInput{
		ParentID:  req.ParentID,
		NameVi:    req.NameVi,
		NameEn:    req.NameEn,
		Slug:      req.Slug,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
	}
}

type UpdateCategory struct {
	ParentID  *string `json:"parent_id"`
	NameVi    *string `json:"name_vi"`
	NameEn    *string `json:"name_en"`
	Slug      *string `json:"slug"`
	Icon      *string `json:"icon"`
	SortOrder *int    `json:"sort_order"`
}

func (req *UpdateCategory) ToUpdateCategoryInput(id string) usecase.UpdateCategoryInput {
	return usecase.UpdateCategoryInput{
		ID:        id,
		ParentID:  req.ParentID,
		NameVi:    req.NameVi,
		NameEn:    req.NameEn,
		Slug:      req.Slug,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
	}
}
