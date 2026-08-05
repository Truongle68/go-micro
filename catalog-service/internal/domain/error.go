package domain

import "errors"

var (
	ErrProductNotFound            = errors.New("product not found")
	ErrCategoryNotFound           = errors.New("category not found")
	ErrEmptyProductID             = errors.New("product id cannot be empty")
	ErrEmptyCategoryID            = errors.New("category id cannot be empty")
	ErrEmptyName                  = errors.New("name cannot be empty")
	ErrEmptySku                   = errors.New("sku cannot be empty")
	ErrInvalidPrice               = errors.New("price cannot be negative")
	ErrInvalidCategoryID          = errors.New("invalid category id")
	ErrInvalidProductID           = errors.New("invalid product id")
	ErrInvalidStatus              = errors.New("invalid product status")
	ErrNoFieldsToUpdate           = errors.New("no fields provided to update")
	ErrDuplicateField             = errors.New("field duplicated")
	ErrDuplicateSKU               = errors.New("duplicated SKU")
	ErrDuplicateSlug              = errors.New("duplicated slug")
	ErrEmptySlug                  = errors.New("slug cannot be empty")
	ErrCircularCategoryParent     = errors.New("category parent cannot be self or descendant")
	ErrInvalidVariantAttribute    = errors.New("invalid variant attribute")
	ErrInvalidOptionType          = errors.New("invalid option type")
	ErrInvalidSimpleVariant       = errors.New("invalid simple variant")
	ErrExceedExpectedVariantCount = errors.New("exceed variant quantity expectation")
	ErrProductRequiresVariant     = errors.New("variant is required")
	ErrInvalidVersion             = errors.New("invalid version")
	ErrConcurrentUpdate           = errors.New("concurrent update")
)

func IsDuplicateKeyError(err error) bool {
	return errors.Is(err, ErrDuplicateField) || errors.Is(err, ErrDuplicateSKU) || errors.Is(err, ErrDuplicateSlug)
}
