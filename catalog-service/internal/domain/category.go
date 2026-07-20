package domain

import "time"

type Category struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	ParentID  *string   `json:"parent_id" bson:"parent_id,omitempty"`
	NameVi    string    `json:"name_vi" bson:"name_vi"`
	NameEn    string    `json:"name_en" bson:"name_en"`
	Slug      string    `json:"slug" bson:"slug"`
	Icon      string    `json:"icon" bson:"icon"`
	SortOrder int       `json:"sort_order" bson:"sort_order"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
