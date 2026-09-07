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
	CodeInactiveVariant             ErrorCode = "INACTIVE_VARIANT"
	CodeInvalidQuantity             ErrorCode = "INVALID_QUANTITY"
	CodeCartItemNotFound            ErrorCode = "CART_ITEM_NOT_FOUND"
	CodeCartItemQtyExceeded         ErrorCode = "CART_ITEM_QUANTITY_EXCEEDED"
	CodeInsufficientStock           ErrorCode = "INSUFFICIENT_STOCK"
	CodeReservationNotFound         ErrorCode = "RESERVATION_NOT_FOUND"
	CodeReservationExpired          ErrorCode = "RESERVATION_EXPIRED"
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
	ErrInactiveVariant             = errors.New("inactive variant")
	ErrInvalidQuantity             = errors.New("invalid quantity: must be positive")
	ErrCartItemNotFound            = errors.New("item not found in cart")
	ErrCartItemQtyExceeded         = errors.New("exceed item quantity in cart")
	ErrInsufficientStock           = errors.New("insufficient stock")
	ErrReservationNotFound         = errors.New("reservation not fount")
	ErrReservationExpired          = errors.New("reservation expired")
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
	ErrInactiveVariant:             CodeInactiveVariant,
	ErrInvalidQuantity:             CodeInvalidQuantity,
	ErrCartItemNotFound:            CodeCartItemNotFound,
	ErrCartItemQtyExceeded:         CodeCartItemQtyExceeded,
	ErrInsufficientStock:           CodeInsufficientStock,
	ErrReservationNotFound:         CodeReservationNotFound,
	ErrReservationExpired:          CodeReservationExpired,
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
