package req

type AdjustStockReq struct {
	SKU           string `json:"sku" binding:"required"`
	WarehouseID   string `json:"warehouse_id" binding:"required"`
	QuantityDelta int    `json:"quantity_delta" binding:"required"`
	Reason        string `json:"reason"`
	Note          string `json:"note"`
}

type UpdateThresholdsReq struct {
	ReorderThreshold int `json:"reorder_threshold" binding:"min=0"`
	ReorderQuantity  int `json:"reorder_quantity" binding:"min=0"`
}
