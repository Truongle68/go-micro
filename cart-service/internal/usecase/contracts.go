package usecase

import (
	"cart-service/internal/domain"
	"cart-service/internal/repo/redis"
	"context"
)

type CartRepository interface {
	GetByUserID(ctx context.Context, userID string) (*domain.Cart, error)
	Save(ctx context.Context, cart *domain.Cart) error
	Delete(ctx context.Context, userID string) error
}

var _ CartRepository = (*redis.CartRepo)(nil)
