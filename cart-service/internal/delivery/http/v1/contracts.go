package v1

import (
	"cart-service/internal/domain"
	"cart-service/internal/usecase"
	"context"
)

type CartUC interface {
	GetCart(ctx context.Context, userID string) (*domain.Cart, error)
	AddItem(ctx context.Context, userID string, sku string, quantity int) (*domain.Cart, error)
	UpdateItemQuantity(ctx context.Context, userID string, sku string, quantity int) (*domain.Cart, error)
	RemoveItem(ctx context.Context, userID string, sku string) (*domain.Cart, error)
	RemoveItems(ctx context.Context, userID string, skus []string) (*domain.Cart, error)
	ClearCart(ctx context.Context, userID string) error
}

var _ CartUC = (*usecase.CartUC)(nil) 
