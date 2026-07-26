package domain

import "errors"

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrEmptyProductID    = errors.New("product id cannot be empty")
	ErrEmptyCategoryID   = errors.New("category id cannot be empty")
	ErrEmptyName         = errors.New("name cannot be empty")
	ErrEmptySku          = errors.New("sku cannot be empty")
	ErrInvalidPrice      = errors.New("price cannot be negative")
	ErrInvalidCategoryID = errors.New("invalid category id")
)
