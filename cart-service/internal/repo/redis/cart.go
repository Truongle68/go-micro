package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cart-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

type CartRepo struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCartRepo(client *redis.Client, ttl time.Duration) *CartRepo {
	return &CartRepo{
		client: client,
		ttl:    ttl,
	}
}

func (r *CartRepo) key(userID string) string {
	return "cart:" + userID
}

func (r *CartRepo) GetByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	data, err := r.client.Get(ctx, r.key(userID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, domain.ErrCartNotFound
		}
		return nil, fmt.Errorf("redis Get failed for user %s: %w", userID, err)
	}

	var cart domain.Cart
	if err := json.Unmarshal(data, &cart); err != nil {
		return nil, fmt.Errorf("json unmarshal cart failed for user %s: %w", userID, err)
	}

	return &cart, nil
}

func (r *CartRepo) Save(ctx context.Context, cart *domain.Cart) error {
	data, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("json marshal cart failed: %w", err)
	}

	if err := r.client.Set(ctx, r.key(cart.UserID), data, r.ttl).Err(); err != nil {
		return fmt.Errorf("redis Set failed for cart user %s: %w", cart.UserID, err)
	}

	return nil
}

func (r *CartRepo) Delete(ctx context.Context, userID string) error {
	if err := r.client.Del(ctx, r.key(userID)).Err(); err != nil {
		return fmt.Errorf("redis Del failed for cart user %s: %w", userID, err)
	}

	return nil
}
