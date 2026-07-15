package redis

import "github.com/redis/go-redis/v9"

type IdentityCache struct {
	*TokenBlacklist
	*ProfileCache
	*OTPCache
}

func NewIdentityCache(client *redis.Client) *IdentityCache {
	return &IdentityCache{
		TokenBlacklist: NewTokenBlacklist(client),
		ProfileCache:   NewProfileCache(client),
		OTPCache:       NewOTPCache(client),
	}
}
