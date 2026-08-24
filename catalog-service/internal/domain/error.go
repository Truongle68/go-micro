package domain

import "errors"

type ErrorCode string

const (
	CodeProductNotFound            ErrorCode = "PRODUCT_NOT_FOUND"
	CodeCategoryNotFound           ErrorCode = "CATEGORY_NOT_FOUND"
	CodeEmptyProductID             ErrorCode = "EMPTY_PRODUCT_ID"
	CodeEmptyCategoryID            ErrorCode = "EMPTY_CATEGORY_ID"
	CodeEmptyName                  ErrorCode = "EMPTY_NAME"
	CodeEmptySKU                   ErrorCode = "EMPTY_SKU"
	CodeSKUNotFound                ErrorCode = "SKU_NOT_FOUND"
	CodeInvalidPrice               ErrorCode = "INVALID_PRICE"
	CodeInvalidPriceRange          ErrorCode = "INVALID_PRICE_RANGE"
	CodeInvalidCategoryID          ErrorCode = "INVALID_CATEGORY_ID"
	CodeInvalidProductID           ErrorCode = "INVALID_PRODUCT_ID"
	CodeInvalidProductStatus       ErrorCode = "INVALID_PRODUCT_STATUS"
	CodeNoFieldsToUpdate           ErrorCode = "NO_FIELDS_TO_UPDATE"
	CodeDuplicateField             ErrorCode = "DUPLICATE_FIELD"
	CodeDuplicateSKU               ErrorCode = "DUPLICATE_SKU"
	CodeDuplicateSlug              ErrorCode = "DUPLICATE_SLUG"
	CodeEmptySlug                  ErrorCode = "EMPTY_SLUG"
	CodeCircularCategoryParent     ErrorCode = "CIRCULAR_CATEGORY_PARENT"
	CodeInvalidVariantAttribute    ErrorCode = "INVALID_VARIANT_ATTRIBUTE"
	CodeInvalidOptionType          ErrorCode = "INVALID_OPTION_TYPE"
	CodeInvalidSimpleVariant       ErrorCode = "INVALID_SIMPLE_VARIANT"
	CodeExceedExpectedVariantCount ErrorCode = "EXCEED_EXPECTED_VARIANT_COUNT"
	CodeProductRequiresVariant     ErrorCode = "PRODUCT_REQUIRES_VARIANT"
	CodeInvalidVersion             ErrorCode = "INVALID_VERSION"
	CodeConcurrentUpdate           ErrorCode = "CONCURRENT_UPDATE"
	CodeInactiveVariant            ErrorCode = "INACTIVE_VARIANT"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

func NewAppError(code ErrorCode, err error) *AppError {
	return &AppError{Code: code, Message: err.Error(), Err: err}
}

var (
	ErrProductNotFound            = errors.New("product not found")
	ErrCategoryNotFound           = errors.New("category not found")
	ErrEmptyProductID             = errors.New("product id cannot be empty")
	ErrEmptyCategoryID            = errors.New("category id cannot be empty")
	ErrEmptyName                  = errors.New("name cannot be empty")
	ErrEmptySKU                   = errors.New("sku cannot be empty")
	ErrSKUNotFound                = errors.New("sku not found")
	ErrInvalidPrice               = errors.New("price cannot be negative")
	ErrInvalidPriceRange          = errors.New("min_price cannot be greater than max_price")
	ErrInvalidCategoryID          = errors.New("invalid category id")
	ErrInvalidProductID           = errors.New("invalid product id")
	ErrInvalidProductStatus       = errors.New("invalid product status")
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
	ErrInactiveVariant            = errors.New("inactive variant")
)

var sentinelToCodeMap = map[error]ErrorCode{
	ErrProductNotFound:            CodeProductNotFound,
	ErrCategoryNotFound:           CodeCategoryNotFound,
	ErrEmptyProductID:             CodeEmptyProductID,
	ErrEmptyCategoryID:            CodeEmptyCategoryID,
	ErrEmptyName:                  CodeEmptyName,
	ErrEmptySKU:                   CodeEmptySKU,
	ErrSKUNotFound:                CodeSKUNotFound,
	ErrInvalidPrice:               CodeInvalidPrice,
	ErrInvalidPriceRange:          CodeInvalidPriceRange,
	ErrInvalidCategoryID:          CodeInvalidCategoryID,
	ErrInvalidProductID:           CodeInvalidProductID,
	ErrInvalidProductStatus:       CodeInvalidProductStatus,
	ErrNoFieldsToUpdate:           CodeNoFieldsToUpdate,
	ErrDuplicateField:             CodeDuplicateField,
	ErrDuplicateSKU:               CodeDuplicateSKU,
	ErrDuplicateSlug:              CodeDuplicateSlug,
	ErrEmptySlug:                  CodeEmptySlug,
	ErrCircularCategoryParent:     CodeCircularCategoryParent,
	ErrInvalidVariantAttribute:    CodeInvalidVariantAttribute,
	ErrInvalidOptionType:          CodeInvalidOptionType,
	ErrInvalidSimpleVariant:       CodeInvalidSimpleVariant,
	ErrExceedExpectedVariantCount: CodeExceedExpectedVariantCount,
	ErrProductRequiresVariant:     CodeProductRequiresVariant,
	ErrInvalidVersion:             CodeInvalidVersion,
	ErrConcurrentUpdate:           CodeConcurrentUpdate,
	ErrInactiveVariant:            CodeInactiveVariant,
}

func ToAppError(err error) *AppError {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	for sentinel, code := range sentinelToCodeMap {
		if errors.Is(err, sentinel) {
			return NewAppError(code, sentinel)
		}
	}
	return nil
}

func IsDuplicateKeyError(err error) bool {
	return errors.Is(err, ErrDuplicateField) || errors.Is(err, ErrDuplicateSKU) || errors.Is(err, ErrDuplicateSlug)
}
