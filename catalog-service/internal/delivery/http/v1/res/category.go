package res

import (
	"catalog-service/internal/usecase"
	"time"
)

type CategoryRead struct {
	ID        string    `json:"id"`
	ParentID  *string   `json:"parent_id"`
	NameVi    string    `json:"name_vi"`
	NameEn    string    `json:"name_en"`
	Slug      string    `json:"slug"`
	Icon      string    `json:"icon"`
	SortOrder int64     `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToCategoryRead(dto *usecase.CategoryDTO) CategoryRead {
	if dto == nil {
		return CategoryRead{}
	}
	return CategoryRead{
		ID:        dto.ID,
		ParentID:  dto.ParentID,
		NameVi:    dto.NameVi,
		NameEn:    dto.NameEn,
		Slug:      dto.Slug,
		Icon:      dto.Icon,
		SortOrder: dto.SortOrder,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func ToCategoryList(dtos []usecase.CategoryDTO) []CategoryRead {
	res := make([]CategoryRead, len(dtos))
	for i, dto := range dtos {
		res[i] = ToCategoryRead(&dto)
	}
	return res
}
