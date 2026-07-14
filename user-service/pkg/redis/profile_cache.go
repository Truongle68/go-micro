package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProfileCache struct {
	client *redis.Client
}

func NewProfileCache(client *redis.Client) *ProfileCache {
	return &ProfileCache{
		client: client,
	}
}

func (c *ProfileCache) InvalidateProfile(ctx context.Context, userID string) error {
	return c.client.Del(ctx, profileKey(userID)).Err()
}

func (c *ProfileCache) CacheProfile(ctx context.Context, userID string, data []byte, ttl time.Duration) error {
	return c.client.Set(ctx, profileKey(userID), data, ttl).Err()
}

func (c *ProfileCache) GetProfile(ctx context.Context, userID string) (string, error) {
	val, err := c.client.Get(ctx, profileKey(userID)).Result()
	return val, err
}

func profileKey(userID string) string {
	return fmt.Sprintf("user:profile:%s", userID)
}
