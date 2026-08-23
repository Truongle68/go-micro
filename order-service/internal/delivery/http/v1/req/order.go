package req

import "order-service/internal/domain"

type CheckoutItemReq struct {
	ProductID    string            `json:"product_id"`
	VariantID    string            `json:"variant_id"`
	SKU          string            `json:"sku" binding:"required"`
	ProductName  string            `json:"product_name"`
	Image        string            `json:"image"`
	VariantAttrs map[string]string `json:"variant_attrs,omitempty"`
	UnitPrice    int64             `json:"unit_price" binding:"min=0"`
	Quantity     int               `json:"quantity" binding:"required,gt=0"`
}

type CheckoutReq struct {
	Items           []CheckoutItemReq      `json:"items,omitempty"`
	ShippingAddress domain.AddressSnapshot `json:"shipping_address" binding:"required"`
	ShippingFee     int64                  `json:"shipping_fee" binding:"min=0"`
	PaymentMethod   string                 `json:"payment_method"`
}

type ShipOrderReq struct {
	TrackingCode string `json:"tracking_code" binding:"required"`
}

type CancelOrderReq struct {
	Reason string `json:"reason"`
}
