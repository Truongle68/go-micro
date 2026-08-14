package domain

import "errors"

type ErrorCode string

const (
	CodeInsufficientStock     ErrorCode = "INSUFFICIENT_STOCK"
	CodeInvalidQuantity       ErrorCode = "INVALID_QUANTITY"
	CodeReservationNotPending ErrorCode = "RESERVATION_NOT_PENDING"
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
	ErrInsufficientStock     = errors.New("insufficient available stock")
	ErrInvalidQuantity       = errors.New("quantity must be positive")
	ErrReservationNotPending = errors.New("reservation is not pending")
)

var sentinelToCodeMap = map[error]ErrorCode{
	ErrInsufficientStock:     CodeInsufficientStock,
	ErrInvalidQuantity:       CodeInvalidQuantity,
	ErrReservationNotPending: CodeReservationNotPending,
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
