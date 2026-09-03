package req

type PurchaseOrderLineReq struct {
	SKU             string `json:"sku" binding:"required"`
	QuantityOrdered int    `json:"quantity_ordered" binding:"required,gt=0"`
	UnitCost        int64  `json:"unit_cost" binding:"min=0"`
}

type CreatePurchaseOrderReq struct {
	Code        string                 `json:"code" binding:"required"`
	SupplierID  string                 `json:"supplier_id" binding:"required"`
	WarehouseID string                 `json:"warehouse_id" binding:"required"`
	Lines       []PurchaseOrderLineReq `json:"lines" binding:"required,min=1,dive"`
}

type ReceiveLineReq struct {
	SKU      string `json:"sku" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}
