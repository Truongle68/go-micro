package publisher

// Event type constants for inventory domain events.
const (
	EventStockReserved  = "stock.reserved"
	EventStockConfirmed = "stock.confirmed"
	EventStockReleased  = "stock.released"
	EventStockAdjusted  = "stock.adjusted"
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
