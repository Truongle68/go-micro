package domain

import "errors"

type ErrorCode string

const (
	CodeNotFound        ErrorCode = "NOT_FOUND"
	CodeCartNotFound    ErrorCode = "CART_NOT_FOUND"
	CodeItemNotFound    ErrorCode = "ITEM_NOT_FOUND"
	CodeInvalidQuantity ErrorCode = "INVALID_QUANTITY"
	CodeInvalidUserID   ErrorCode = "INVALID_USER_ID"
	CodeInvalidSKU      ErrorCode = "INVALID_SKU"
	CodeEmptyCart       ErrorCode = "EMPTY_CART"
	CodeEmptySKU        ErrorCode = "SKU_EMPTY"
	CodeSKUNotFound     ErrorCode = "SKU_NOT_FOUND"
	CodeInactiveVariant ErrorCode = "INACTIVE_VARIANT"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

func NewAppError(code ErrorCode, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: err.Error(),
		Err:     err,
	}
}

// Domain & HTTP Error Messages / Sentinels
var (
	ErrCartNotFound         = errors.New("cart not found")
	ErrItemNotFound         = errors.New("item not found in cart")
	ErrInvalidQuantity      = errors.New("quantity must be greater than zero")
	ErrInvalidUserID        = errors.New("user ID cannot be empty")
	ErrInvalidSKU           = errors.New("SKU cannot be empty")
	ErrEmptyCart            = errors.New("cart is empty")
	ErrSKURequired          = errors.New("sku parameter is required")
	ErrFailedToRetrieveCart = errors.New("failed to retrieve cart")
	ErrFailedToAddToCart    = errors.New("failed to add item to cart")
	ErrFailedToUpdateItem   = errors.New("failed to update item quantity")
	ErrFailedToRemoveItem   = errors.New("failed to remove item from cart")
	ErrFailedToClearCart    = errors.New("failed to clear cart")
	ErrEmptySKU             = errors.New("sku cannot be empty")
	ErrSKUNotFound          = errors.New("sku not found")
	ErrInactiveVariant      = errors.New("inactive variant")
)

var sentinelToCodeMap = map[error]ErrorCode{
	ErrCartNotFound:         CodeCartNotFound,
	ErrItemNotFound:         CodeItemNotFound,
	ErrInvalidQuantity:      CodeInvalidQuantity,
	ErrInvalidUserID:        CodeInvalidUserID,
	ErrInvalidSKU:           CodeInvalidSKU,
	ErrEmptyCart:            CodeEmptyCart,
	ErrSKURequired:          CodeInvalidSKU,
	ErrFailedToRetrieveCart: CodeCartNotFound,
	ErrEmptySKU:             CodeEmptySKU,
	ErrSKUNotFound:          CodeSKUNotFound,
	ErrInactiveVariant:      CodeInactiveVariant,
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
