package redis

import (
	"context"
	"time"
)

type IdentityCacher interface {
	BlacklistCacher
}

type BlacklistCacher interface {
	Blacklist(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}
