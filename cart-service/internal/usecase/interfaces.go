package usecase

import (
	"context"

	"cart-service/internal/domain"
)

type Cart interface {
	GetCart(ctx context.Context, userID string) (*domain.Cart, error)
	AddItem(ctx context.Context, userID string, sku string, quantity int) (*domain.Cart, error)
	UpdateItemQuantity(ctx context.Context, userID string, sku string, quantity int) (*domain.Cart, error)
	RemoveItem(ctx context.Context, userID string, sku string) (*domain.Cart, error)
	ClearCart(ctx context.Context, userID string) error
}
