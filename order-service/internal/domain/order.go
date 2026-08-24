package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "pending_payment"
	OrderStatusFailedStock    OrderStatus = "failed_stock"
	OrderStatusFailedPayment  OrderStatus = "failed_payment"
	OrderStatusConfirmed      OrderStatus = "confirmed"
	OrderStatusPreparing      OrderStatus = "preparing"
	OrderStatusShipped        OrderStatus = "shipped"
	OrderStatusDelivered      OrderStatus = "delivered"
	OrderStatusCancelled      OrderStatus = "cancelled"
	OrderStatusRefunded       OrderStatus = "refunded"
)

type AddressSnapshot struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Street   string `json:"street"`
	Ward     string `json:"ward,omitempty"`
	District string `json:"district,omitempty"`
	City     string `json:"city"`
	Country  string `json:"country,omitempty"`
}

type OrderItem struct {
	ID           string            `json:"id"`
	OrderID      string            `json:"order_id"`
	ProductID    string            `json:"product_id"`
	VariantID    string            `json:"variant_id"`
	SKU          string            `json:"sku"`
	ProductName  string            `json:"product_name"`
	Image        string            `json:"image"`
	VariantAttrs map[string]string `json:"variant_attrs,omitempty"`
	UnitPrice    int64             `json:"unit_price"`
	Quantity     int               `json:"quantity"`
}

type OrderStatusHistory struct {
	ID         string      `json:"id"`
	OrderID    string      `json:"order_id"`
	FromStatus OrderStatus `json:"from_status"`
	ToStatus   OrderStatus `json:"to_status"`
	Note       string      `json:"note,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
}

type Order struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	Status          OrderStatus     `json:"status"`
	Items           []OrderItem     `json:"items"`
	ShippingAddress AddressSnapshot `json:"shipping_address"`
	Subtotal        int64           `json:"subtotal"`
	ShippingFee     int64           `json:"shipping_fee"`
	Total           int64           `json:"total"`
	PaymentMethod   string          `json:"payment_method"`
	PaymentID       string          `json:"payment_id,omitempty"`
	TrackingCode    string          `json:"tracking_code,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func NewOrder(userID string, items []OrderItem, address AddressSnapshot, shippingFee int64, paymentMethod string) (*Order, *OrderStatusHistory, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, nil, ErrInvalidUserID
	}
	if len(items) == 0 {
		return nil, nil, ErrEmptyOrderItems
	}
	if strings.TrimSpace(address.FullName) == "" || strings.TrimSpace(address.Phone) == "" || strings.TrimSpace(address.Street) == "" || strings.TrimSpace(address.City) == "" {
		return nil, nil, ErrInvalidAddress
	}

	orderID := uuid.New().String()
	now := time.Now().UTC()

	var subtotal int64
	processedItems := make([]OrderItem, len(items))
	for i, item := range items {
		itemID := item.ID
		if itemID == "" {
			itemID = uuid.New().String()
		}
		subtotal += item.UnitPrice * int64(item.Quantity)
		processedItems[i] = OrderItem{
			ID:           itemID,
			OrderID:      orderID,
			ProductID:    item.ProductID,
			VariantID:    item.VariantID,
			SKU:          item.SKU,
			ProductName:  item.ProductName,
			Image:        item.Image,
			VariantAttrs: item.VariantAttrs,
			UnitPrice:    item.UnitPrice,
			Quantity:     item.Quantity,
		}
	}

	if paymentMethod == "" {
		paymentMethod = "cod"
	}
	if address.Country == "" {
		address.Country = "Vietnam"
	}

	initialStatus := OrderStatusPendingPayment
	order := &Order{
		ID:              orderID,
		UserID:          userID,
		Status:          initialStatus,
		Items:           processedItems,
		ShippingAddress: address,
		Subtotal:        subtotal,
		ShippingFee:     shippingFee,
		Total:           subtotal + shippingFee,
		PaymentMethod:   paymentMethod,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	history := &OrderStatusHistory{
		ID:         uuid.New().String(),
		OrderID:    orderID,
		FromStatus: "",
		ToStatus:   initialStatus,
		Note:       "Order created",
		CreatedAt:  now,
	}

	return order, history, nil
}

func (o *Order) MarkConfirmed(paymentID string) (*OrderStatusHistory, error) {
	if o.Status != OrderStatusPendingPayment {
		return nil, ErrInvalidOrderTransition
	}

	oldStatus := o.Status
	o.Status = OrderStatusConfirmed
	if paymentID != "" {
		o.PaymentID = paymentID
	}
	o.UpdatedAt = time.Now().UTC()

	return &OrderStatusHistory{
		ID:         uuid.New().String(),
		OrderID:    o.ID,
		FromStatus: oldStatus,
		ToStatus:   o.Status,
		Note:       "Payment confirmed",
		CreatedAt:  o.UpdatedAt,
	}, nil
}

func (o *Order) MarkFailed(status OrderStatus, note string) (*OrderStatusHistory, error) {
	if o.Status != OrderStatusPendingPayment {
		return nil, ErrInvalidOrderTransition
	}
	if status != OrderStatusFailedStock && status != OrderStatusFailedPayment {
		return nil, ErrInvalidOrderTransition
	}

	oldStatus := o.Status
	o.Status = status
	o.UpdatedAt = time.Now().UTC()

	return &OrderStatusHistory{
		ID:         uuid.New().String(),
		OrderID:    o.ID,
		FromStatus: oldStatus,
		ToStatus:   o.Status,
		Note:       note,
		CreatedAt:  o.UpdatedAt,
	}, nil
}

func (o *Order) Prepare() (*OrderStatusHistory, error) {
	if o.Status != OrderStatusConfirmed {
		return nil, ErrInvalidOrderTransition
	}

	oldStatus := o.Status
	o.Status = OrderStatusPreparing
	o.UpdatedAt = time.Now().UTC()

	return &OrderStatusHistory{
		ID:         uuid.New().String(),
		OrderID:    o.ID,
		FromStatus: oldStatus,
		ToStatus:   o.Status,
		Note:       "Order is being prepared in warehouse",
		CreatedAt:  o.UpdatedAt,
	}, nil
}

func (o *Order) Ship(trackingCode string) (*OrderStatusHistory, error) {
	if o.Status != OrderStatusPreparing && o.Status != OrderStatusConfirmed {
		return nil, ErrInvalidOrderTransition
	}

	oldStatus := o.Status
	o.Status = OrderStatusShipped
	o.TrackingCode = trackingCode
	o.UpdatedAt = time.Now().UTC()

	return &OrderStatusHistory{
		ID:         uuid.New().String(),
		OrderID:    o.ID,
		FromStatus: oldStatus,
		ToStatus:   o.Status,
		Note:       "Order shipped with tracking code: " + trackingCode,
		CreatedAt:  o.UpdatedAt,
	}, nil
}

func (o *Order) Deliver() (*OrderStatusHistory, error) {
	if o.Status != OrderStatusShipped {
		return nil, ErrInvalidOrderTransition
	}

	oldStatus := o.Status
	o.Status = OrderStatusDelivered
	o.UpdatedAt = time.Now().UTC()

	return &OrderStatusHistory{
		ID:         uuid.New().String(),
		OrderID:    o.ID,
		FromStatus: oldStatus,
		ToStatus:   o.Status,
		Note:       "Order delivered successfully",
		CreatedAt:  o.UpdatedAt,
	}, nil
}

func (o *Order) Cancel(reason string) (*OrderStatusHistory, error) {
	switch o.Status {
	case OrderStatusCancelled:
		return nil, ErrOrderAlreadyCancelled
	case OrderStatusShipped:
		return nil, ErrCannotCancelShippedOrder
	case OrderStatusDelivered:
		return nil, ErrCannotCancelDeliveriedOrder
	case OrderStatusRefunded:
		return nil, ErrOrderAlreadyRefunded
	}

	oldStatus := o.Status
	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now().UTC()

	note := "Order cancelled"
	if reason != "" {
		note = "Order cancelled: " + reason
	}

	return &OrderStatusHistory{
		ID:         uuid.New().String(),
		OrderID:    o.ID,
		FromStatus: oldStatus,
		ToStatus:   o.Status,
		Note:       note,
		CreatedAt:  o.UpdatedAt,
	}, nil
}
