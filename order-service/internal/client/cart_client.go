package client

import (
	"context"
)

type CartItemDTO struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type CartDTO struct {
	UserID string        `json:"user_id"`
	Items  []CartItemDTO `json:"items"`
}

type CartClient interface {
	GetCart(ctx context.Context, userID string, token string) (*CartDTO, error)
	ClearCart(ctx context.Context, userID string, token string) error
}
