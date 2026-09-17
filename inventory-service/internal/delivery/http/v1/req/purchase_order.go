package req

import (
	"inventory-service/internal/domain"
)

type PurchaseOrderLineReq struct {
	SKU             string `json:"sku" binding:"required"`
	QuantityOrdered int    `json:"quantity_ordered" binding:"required,gt=0"`
	UnitCost        int64  `json:"unit_cost" binding:"min=0"`
}

type CreatePurchaseOrderReq struct {
	Code          string                 `json:"code" binding:"required"`
	SupplierID    string                 `json:"supplier_id" binding:"required"`
	WarehouseID   string                 `json:"warehouse_id" binding:"required"`
	CreatedByName string                 `json:"created_by_name" binding:"required"`
	Lines         []PurchaseOrderLineReq `json:"lines" binding:"required,min=1,dive"`
}

type ReceiveLineReq struct {
	SKU      string `json:"sku" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

type ReceiveGoodsReq struct {
	Lines []ReceiveLineInput `json:"lines" binding:"required,min=1,dive"`
}

type ReceiveLineInput struct {
	SKU      string `json:"sku" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

func (req ReceiveGoodsReq) ToReceiptLines() []domain.ReceiveLine {
	lines := make([]domain.ReceiveLine, len(req.Lines))
	for i, l := range req.Lines {
		lines[i] = domain.ReceiveLine{
			SKU:      l.SKU,
			Quantity: l.Quantity,
		}
	}
	return lines
}
