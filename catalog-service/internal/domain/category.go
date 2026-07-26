package domain

import "time"

type Category struct {
	ID        string
	ParentID  *string
	NameVi    string
	NameEn    string
	Slug      string
	Icon      string
	SortOrder int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewCategoryParams struct {
	ParentID  *string
	NameVi    string
	NameEn    string
	Slug      string
	Icon      string
	SortOrder int64
}

func NewCategory(params NewCategoryParams) (*Category, error) {
	if params.NameVi == "" && params.NameEn == "" {
		return nil, ErrEmptyName
	}

	if params.Slug == "" {
		return nil, ErrEmptySku
	}

	now := time.Now()
	return &Category{
		ParentID:  params.ParentID,
		NameVi:    params.NameVi,
		NameEn:    params.NameEn,
		Slug:      params.Slug,
		Icon:      params.Icon,
		SortOrder: params.SortOrder,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

type UpdateCategoryParams struct {
	ID        string
	ParentID  *string
	NameVi    *string
	NameEn    *string
	Slug      *string
	Icon      *string
	SortOrder *int64
}

func (c *Category) ApplyUpdate(params UpdateCategoryParams) {
	if params.ParentID != nil {
		c.ParentID = params.ParentID
	}
	if params.NameVi != nil {
		c.NameVi = *params.NameVi
	}
	if params.NameEn != nil {
		c.NameEn = *params.NameEn
	}
	if params.Slug != nil {
		c.Slug = *params.Slug
	}
	if params.Icon != nil {
		c.Icon = *params.Icon
	}
	if params.SortOrder != nil {
		c.SortOrder = *params.SortOrder
	}

	c.UpdatedAt = time.Now()
}
