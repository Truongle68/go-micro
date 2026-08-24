package domain

import "errors"

type ErrorCode string

const (
	CodeOrderNotFound               ErrorCode = "ORDER_NOT_FOUND"
	CodeInvalidOrderTransition      ErrorCode = "INVALID_ORDER_TRANSITION"
	CodeEmptyOrderItems             ErrorCode = "EMPTY_ORDER_ITEMS"
	CodeInvalidUserID               ErrorCode = "INVALID_USER_ID"
	CodeInvalidAddress              ErrorCode = "INVALID_ADDRESS"
	CodeOrderAlreadyCancelled       ErrorCode = "ORDER_ALREADY_CANCELLED"
	CodeCannotCancelShippedOrder    ErrorCode = "CANNOT_CANCEL_SHIPPED_ORDER"
	CodeCannotCancelDeliveriedOrder ErrorCode = "CANNOT_CANCEL_DELIVERIED_ORDER"
	CodeOrderAlreadyRefunded        ErrorCode = "ORDER_ALREADY_REFUNDED"
	CodeCartFetchFailed             ErrorCode = "CART_FETCH_FAILED"
	CodeCartEmpty                   ErrorCode = "CART_EMPTY"
	CodeEmptySKU                    ErrorCode = "SKU_EMPTY"
	CodeSKUNotFound                 ErrorCode = "SKU_NOT_FOUND"
	CodeInvalidQuantity             ErrorCode = "INVALID_QUANTITY"
	CodeInactiveVariant             ErrorCode = "INACTIVE_VARIANT"
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

var (
	ErrOrderNotFound               = errors.New("order not found")
	ErrInvalidOrderTransition      = errors.New("invalid order status transition")
	ErrEmptyOrderItems             = errors.New("order items cannot be empty")
	ErrInvalidUserID               = errors.New("user ID cannot be empty")
	ErrInvalidAddress              = errors.New("invalid shipping address")
	ErrOrderAlreadyCancelled       = errors.New("order is already cancelled")
	ErrOrderAlreadyRefunded        = errors.New("order is already refunded")
	ErrCannotCancelShippedOrder    = errors.New("cannot cancel order that is already shipped")
	ErrCannotCancelDeliveriedOrder = errors.New("cannot cancel order that is already deliveried")
	ErrCartFetchFailed             = errors.New("failed to fetch cart from cart service")
	ErrCartEmpty                   = errors.New("cart cannot be empty")
	ErrEmptySKU                    = errors.New("sku cannot be empty")
	ErrSKUNotFound                 = errors.New("sku not found")
	ErrInvalidQuantity             = errors.New("invalid quantity: must be positive")
	ErrInactiveVariant             = errors.New("inactive variant")
)

var sentinelToCodeMap = map[error]ErrorCode{
	ErrOrderNotFound:               CodeOrderNotFound,
	ErrInvalidOrderTransition:      CodeInvalidOrderTransition,
	ErrEmptyOrderItems:             CodeEmptyOrderItems,
	ErrInvalidUserID:               CodeInvalidUserID,
	ErrInvalidAddress:              CodeInvalidAddress,
	ErrOrderAlreadyCancelled:       CodeOrderAlreadyCancelled,
	ErrCannotCancelDeliveriedOrder: CodeCannotCancelDeliveriedOrder,
	ErrCannotCancelShippedOrder:    CodeCannotCancelShippedOrder,
	ErrOrderAlreadyRefunded:        CodeOrderAlreadyRefunded,
	ErrCartFetchFailed:             CodeCartFetchFailed,
	ErrCartEmpty:                   CodeCartEmpty,
	ErrEmptySKU:                    CodeEmptySKU,
	ErrSKUNotFound:                 CodeSKUNotFound,
	ErrInvalidQuantity:             CodeInvalidQuantity,
	ErrInactiveVariant:             CodeInactiveVariant,
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
