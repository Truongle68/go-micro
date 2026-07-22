package redis

import "github.com/redis/go-redis/v9"

type IdentityCache struct {
	*TokenBlacklist
}

func NewIdentityCache(client *redis.Client) *IdentityCache {
	return &IdentityCache{
		TokenBlacklist: NewTokenBlacklist(client),
	}
}
