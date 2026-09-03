package client

import "context"

// VariantDTO represents product variant data from catalog-service.
type VariantDTO struct {
	ID          string            `json:"id"`
	ProductID   string            `json:"product_id"`
	ProductName string            `json:"product_name"`
	SKU         string            `json:"sku"`
	Attributes  map[string]string `json:"attributes"`
	Price       Price             `json:"price"`
	Image       string            `json:"image"`
	IsActive    bool              `json:"is_active"`
}

// Price represents price amount and currency.
type Price struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

// CatalogClient defines interactions with the catalog-service.
type CatalogClient interface {
	GetVariantsBySKUs(ctx context.Context, skus []string) ([]VariantDTO, error)
}
