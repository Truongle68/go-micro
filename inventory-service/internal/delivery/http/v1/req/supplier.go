package req

type SupplierAddressReq struct {
	Line1    string `json:"line1"`
	Ward     string `json:"ward"`
	District string `json:"district"`
	City     string `json:"city"`
}

type CreateSupplierReq struct {
	Code    string             `json:"code" binding:"required"`
	Name    string             `json:"name" binding:"required"`
	Phone   string             `json:"phone"`
	Email   string             `json:"email"`
	Address SupplierAddressReq `json:"address"`
}

type UpdateSupplierReq struct {
	Name    string             `json:"name" binding:"required"`
	Email   string             `json:"email"`
	Phone   string             `json:"phone"`
	Address SupplierAddressReq `json:"address"`
}
