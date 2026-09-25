package rabbitmq

const (
	// Exchange name constants
	ExchangeOrder       = "order.events"
	ExchangeOrderPlaced = "order.placed.fanout"
	ExchangeInventory   = "inventory.events"
	ExchangeCatalog     = "catalog.events"
	// Exchange type constants
	ExchangeTypeDirect  = "direct"
	ExchangeTypeTopic   = "topic"
	ExchangeTypeFanout  = "fanout"
	ExchangeTypeHeaders = "headers"
)

// SKUQtyItem represents a single SKU and its quantity in an event payload.
type SKUQtyItem struct {
	SKU         string `json:"sku"`
	WarehouseID string `json:"warehouse_id"`
	Quantity    int    `json:"quantity"`
}

// StockReservedPayload is the payload for stock.reserved events.
type StockReservedPayload struct {
	OrderID string       `json:"order_id"`
	Items   []SKUQtyItem `json:"items"`
}

// StockConfirmedPayload is the payload for stock.confirmed events.
type StockConfirmedPayload struct {
	OrderID string       `json:"order_id"`
	Items   []SKUQtyItem `json:"items"`
}

// StockReleasedPayload is the payload for stock.released events.
type StockReleasedPayload struct {
	OrderID string       `json:"order_id"`
	Items   []SKUQtyItem `json:"items"`
}

// StockAdjustedPayload is the payload for stock.adjusted events.
type StockAdjustedPayload struct {
	SKU         string `json:"sku"`
	WarehouseID string `json:"warehouse_id"`
	Delta       int    `json:"delta"`
	NewOnHand   int    `json:"new_on_hand"`
}
