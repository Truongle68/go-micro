package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBlacklist struct {
	client *redis.Client
}

func NewTokenBlacklist(client *redis.Client) *TokenBlacklist {
	return &TokenBlacklist{
		client: client,
	}
}

var _ BlacklistCacher = (*TokenBlacklist)(nil)

func (b *TokenBlacklist) Add(ctx context.Context, jti string, ttl time.Duration) error {
	return b.client.Set(ctx, blacklistKey(jti), 1, ttl).Err()
}
func (b *TokenBlacklist) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	n, err := b.client.Exists(ctx, blacklistKey(jti)).Result()
	return n > 0, err
}

func blacklistKey(jti string) string {
	return fmt.Sprintf("blacklist:%s", jti)
}
