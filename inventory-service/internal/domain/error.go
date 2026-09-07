package domain

import "errors"

type ErrorCode string

const (
	CodeInsufficientStock     ErrorCode = "INSUFFICIENT_STOCK"
	CodeNonPositiveQuantity   ErrorCode = "NON_POSITIVE_QUANTITY"
	CodeReservationNotPending ErrorCode = "RESERVATION_NOT_PENDING"
	// Warehouse error code
	CodeEmptyWhCode        ErrorCode = "EMPTY_WH_CODE"
	CodeEmptyWhName        ErrorCode = "EMPTY_WH_NAME"
	CodeEmptyWhAddressCity ErrorCode = "EMPTY_WH_ADDRESS_CITY"
	// Stock transfer error code
	CodeInvalidTransferRoute ErrorCode = "INVALID_TRANSFER_ROUTE"
	CodeEmptySKU             ErrorCode = "EMPTY_SKU"
	CodeEmptyWarehouseID     ErrorCode = "EMPTY_WAREHOUSE_ID"
	CodeEmptyFromWarehouseID ErrorCode = "EMPTY_FROM_WAREHOUSE_ID"
	CodeEmptyToWarehouseID   ErrorCode = "EMPTY_TO_WAREHOUSE_ID"
	CodeNegativeQuantity     ErrorCode = "NEGATIVE_QUANTITY"
	// Purchase order error code
	CodeEmptySuppID             ErrorCode = "EMPTY_SUPPLIER_ID"
	CodeEmptyPurchaseOrderLines ErrorCode = "EMPTY_PO_LINE"
	CodePODuplicateSKU          ErrorCode = "PO_DUPPLICATE_SKU"
	CodeInvalidPOTransition     ErrorCode = "INVALID_PO_TRANSITION"
	CodeReceivedExceedsOrdered  ErrorCode = "RECEIVED_EXCEEDS_ORDERED"
	// Supplier error code
	CodeEmptySuppName          ErrorCode = "EMPTY_SUPPLIER_NAME"
	CodeEmptySuppCode          ErrorCode = "EMPTY_SUPPLIER_CODE"
	CodeEmptySuppPhone         ErrorCode = "EMPTY_SUPPLIER_PHONE"
	CodeEmptySuppEmail         ErrorCode = "EMPTY_SUPPLIER_EMAIL"
	CodeSuppAlreadyActive      ErrorCode = "SUPPLIER_ALREADY_ACTIVE"
	CodeSuppAlreadyInactive    ErrorCode = "SUPPLIER_ALREADY_INACTIVE"
	CodeSuppNotFound           ErrorCode = "SUPPLIER_NOT_FOUND"
	CodePONotFound             ErrorCode = "PURCHASE_ORDER_NOT_FOUND"
	CodeConcurrentModification ErrorCode = "CONCURRENT_MODIFICATION"
	CodeWhNotFound             ErrorCode = "WAREHOUSE_NOT_FOUND"
	CodeWhInactive             ErrorCode = "WAREHOUSE_INACTIVE"
	CodeSKUNotFound            ErrorCode = "SKU_NOT_FOUND"
	CodeInactiveVariant        ErrorCode = "INACTIVE_VARIANT"
	CodeDuplicatePOCode        ErrorCode = "DUPLICATE_PURCHASE_ORDER_CODE"
	CodeDuplicateSuppCode      ErrorCode = "DUPLICATE_SUPPLIER_CODE"
	CodeStockLevelNotFound     ErrorCode = "STOCK_LEVEL_NOT_FOUND"
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
	ErrNonPositiveQuantity   = errors.New("quantity must be positive")
	ErrReservationNotPending = errors.New("reservation is not pending")
	// Warehouse error msg
	ErrEmptyWhCode        = errors.New("warehouse code cannot be empty")
	ErrEmptyWhName        = errors.New("warehouse name cannot be empty")
	ErrEmptyWhAddressCity = errors.New("warehouse address city cannot be empty")
	// Stock transfer error msg
	ErrInvalidTransferRoute = errors.New("source cannot be coicided with target")
	// Stock level error msg
	ErrEmptySKU      = errors.New("sku cannot be empty")
	ErrEmptyWhID     = errors.New("warehouse ID cannot be empty")
	ErrEmptyFromWhID = errors.New("from warehouse ID cannot be empty")
	ErrEmptyToWhID   = errors.New("to warehouse ID cannot be empty")
	// Purchase Order error msg
	ErrEmptySuppID             = errors.New("supplier id cannot be empty")
	ErrEmptyPurchaseOrderLines = errors.New("purchase order must have at least one line")
	ErrPODuplicateSKU          = errors.New("product with this SKU already exist in purchase order")
	ErrInvalidPOTransition     = errors.New("invalid purchase order status transition")
	ErrReceivedExceedsOrdered  = errors.New("received quantity exceeds ordered quantity")
	// Supplier error msg
	ErrEmptySuppName          = errors.New("supplier name cannot be empty")
	ErrEmptySuppCode          = errors.New("supplier code cannot be empty")
	ErrEmptySuppPhone         = errors.New("supplier phone cannot be empty")
	ErrEmptySuppEmail         = errors.New("supplier email cannot be empty")
	ErrSuppAlreadyActive      = errors.New("supplier already active")
	ErrSuppAlreadyInactive    = errors.New("supplier already inactive")
	ErrSuppNotFound           = errors.New("supplier not found")
	ErrPONotFound             = errors.New("purchase order not found")
	ErrConcurrentModification = errors.New("concurrent modification detected")
	ErrWhNotFound             = errors.New("warehouse not found")
	ErrWhAlreadyInactive      = errors.New("warehouse already inactive")
	ErrSKUNotFound            = errors.New("sku not found in catalog")
	ErrInactiveVariant        = errors.New("product variant is inactive")
	ErrDuplicatePOCode        = errors.New("purchase order with this code already exists")
	ErrDuplicateSuppCode      = errors.New("supplier with this code already exists")
	ErrStockLevelNotFound     = errors.New("stock level not found")

	// Business Rule / Invariant Errors
	ErrNegativeQuantity = errors.New("reorder threshold and quantity cannot be negative")
	ErrInvalidPrice     = errors.New("price cannot be negative")
)

var sentinelToCodeMap = map[error]ErrorCode{
	ErrInsufficientStock:       CodeInsufficientStock,
	ErrNonPositiveQuantity:     CodeNonPositiveQuantity,
	ErrReservationNotPending:   CodeReservationNotPending,
	ErrEmptyWhCode:             CodeEmptyWhCode,
	ErrEmptyWhName:             CodeEmptyWhName,
	ErrEmptyWhAddressCity:      CodeEmptyWhAddressCity,
	ErrInvalidTransferRoute:    CodeInvalidTransferRoute,
	ErrEmptySKU:                CodeEmptySKU,
	ErrEmptyWhID:               CodeEmptyWarehouseID,
	ErrEmptyFromWhID:           CodeEmptyFromWarehouseID,
	ErrEmptyToWhID:             CodeEmptyToWarehouseID,
	ErrNegativeQuantity:        CodeNegativeQuantity,
	ErrEmptySuppName:           CodeEmptySuppName,
	ErrEmptySuppCode:           CodeEmptySuppCode,
	ErrEmptySuppEmail:          CodeEmptySuppEmail,
	ErrEmptySuppPhone:          CodeEmptySuppPhone,
	ErrSuppAlreadyActive:       CodeSuppAlreadyActive,
	ErrSuppAlreadyInactive:     CodeSuppAlreadyInactive,
	ErrEmptySuppID:             CodeEmptySuppID,
	ErrEmptyPurchaseOrderLines: CodeEmptyPurchaseOrderLines,
	ErrPODuplicateSKU:          CodePODuplicateSKU,
	ErrInvalidPOTransition:     CodeInvalidPOTransition,
	ErrReceivedExceedsOrdered:  CodeReceivedExceedsOrdered,
	ErrSuppNotFound:            CodeSuppNotFound,
	ErrPONotFound:              CodePONotFound,
	ErrConcurrentModification:  CodeConcurrentModification,
	ErrWhNotFound:              CodeWhNotFound,
	ErrWhAlreadyInactive:       CodeWhInactive,
	ErrSKUNotFound:             CodeSKUNotFound,
	ErrInactiveVariant:         CodeInactiveVariant,
	ErrDuplicatePOCode:         CodeDuplicatePOCode,
	ErrDuplicateSuppCode:       CodeDuplicateSuppCode,
	ErrStockLevelNotFound:      CodeStockLevelNotFound,
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
