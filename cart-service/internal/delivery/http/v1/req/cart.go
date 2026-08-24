package req

type AddItemReq struct {
	SKU      string `json:"sku" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

type UpdateItemReq struct {
	Quantity int `json:"quantity" binding:"min=0"`
}

type RemoveItemsReq struct {
	SKUs []string `json:"skus" binding:"required,gt=0"`
}
