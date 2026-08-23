package domain

import "time"

type Category struct {
	ID              string
	ParentID        *string
	Name            string
	NameTranslation map[string]string
	Slug            string
	Icon            string
	SortOrder       int
	IsActive        bool
	Ancestors       []CategoryRef
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type NewCategoryParams struct {
	ParentID        *string
	Name            string
	NameTranslation map[string]string
	Slug            string
	Icon            string
	SortOrder       int
	IsActive        *bool
	Ancestors       []CategoryRef
}

func NewCategory(params NewCategoryParams) (*Category, error) {
	if params.Name == "" {
		return nil, ErrEmptyName
	}

	if params.Slug == "" {
		return nil, ErrEmptySlug
	}

	isActive := true
	if params.IsActive != nil {
		isActive = *params.IsActive
	}

	now := time.Now()
	return &Category{
		ParentID:        params.ParentID,
		Name:            params.Name,
		NameTranslation: params.NameTranslation,
		Slug:            params.Slug,
		Icon:            params.Icon,
		SortOrder:       params.SortOrder,
		IsActive:        isActive,
		Ancestors:       params.Ancestors,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

type UpdateCategoryParams struct {
	ID              string
	ParentID        *string
	Name            *string
	NameTranslation map[string]string
	Slug            *string
	Icon            *string
	SortOrder       *int
	IsActive        *bool
	Ancestors       []CategoryRef
}

func (c *Category) ApplyUpdate(params UpdateCategoryParams) {
	if params.ParentID != nil {
		c.ParentID = params.ParentID
	}
	if params.Name != nil {
		c.Name = *params.Name
	}
	if params.NameTranslation != nil {
		c.NameTranslation = params.NameTranslation
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
	if params.IsActive != nil {
		c.IsActive = *params.IsActive
	}
	if params.Ancestors != nil {
		c.Ancestors = params.Ancestors
	}

	c.UpdatedAt = time.Now()
}

type ListCategoryResult struct {
	Categories []Category
	TotalCount int64
}
