package entity

import "time"

type Category struct {
	ID        string    `json:"id"`
	ParentID  *string   `json:"parent_id"`
	NameVi    string    `json:"name_vi"`
	NameEn    string    `json:"name_en"`
	Slug      string    `json:"slug"`
	Icon      string    `json:"icon"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
