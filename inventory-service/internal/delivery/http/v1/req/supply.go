package req

import "inventory-service/internal/domain"

type AssignProductsReq struct {
	SupplierID string        `json:"supplier_id"`
	Supplies   []SupplyInput `json:"supplies"`
}

type SupplyInput struct {
	SKU         string   `json:"sku"`
	Price       int64    `json:"price"`
	LeadTime    LeadTime `json:"lead_time"`
	MinOrderQty int      `json:"min_order_qty"`
}

func (in SupplyInput) ToDomain() domain.SupplyInput {
	return domain.SupplyInput{
		SKU:   in.SKU,
		Price: in.Price,
		LeadTime: domain.LeadTime{
			Days: in.LeadTime.Days,
		},
		MinOrderQty: in.MinOrderQty,
	}
}

type LeadTime struct {
	Days int `json:"days"`
}
