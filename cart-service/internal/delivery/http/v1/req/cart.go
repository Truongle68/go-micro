package req

type AddItemReq struct {
	SKU      string `json:"sku" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

type UpdateItemReq struct {
	Quantity int `json:"quantity" binding:"min=0"`
}
