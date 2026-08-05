package res

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/usecase"
	"time"
)

type CategoryRead struct {
	ID              string                `json:"id"`
	ParentID        *string               `json:"parent_id"`
	Name            string                `json:"name"`
	NameTranslation map[string]string     `json:"name_translation,omitempty"`
	Slug            string                `json:"slug"`
	Icon            string                `json:"icon,omitempty"`
	SortOrder       int64                 `json:"sort_order"`
	IsActive        bool                  `json:"is_active"`
	Ancestors       []domain.CategoryRef  `json:"ancestors,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

func ToCategoryRead(dto *usecase.CategoryDTO) CategoryRead {
	if dto == nil {
		return CategoryRead{}
	}
	return CategoryRead{
		ID:              dto.ID,
		ParentID:        dto.ParentID,
		Name:            dto.Name,
		NameTranslation: dto.NameTranslation,
		Slug:            dto.Slug,
		Icon:            dto.Icon,
		SortOrder:       dto.SortOrder,
		IsActive:        dto.IsActive,
		Ancestors:       dto.Ancestors,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}

func ToCategoryList(dtos []usecase.CategoryDTO) []CategoryRead {
	res := make([]CategoryRead, len(dtos))
	for i, dto := range dtos {
		res[i] = ToCategoryRead(&dto)
	}
	return res
}
