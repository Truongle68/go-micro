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

func (b *TokenBlacklist) Blacklist(ctx context.Context, token string, ttl time.Duration) error {
	return b.client.Set(ctx, blacklistKey(token), 1, ttl).Err()
}
func (b *TokenBlacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	n, err := b.client.Exists(ctx, blacklistKey(token)).Result()
	return n > 0, err
}

func blacklistKey(token string) string {
	return fmt.Sprintf("blacklist:%s", token)
}
