package redis

import (
	"context"
	"time"
)

type BlacklistCacher interface {
	Blacklist(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type ProfileCacheReader interface {
	GetProfile(ctx context.Context, userID string) (string, error)
}

type ProfileCacheWriter interface {
	CacheProfile(ctx context.Context, userID string, data []byte, ttl time.Duration) error
}

type ProfileInvalidator interface {
	InvalidateProfile(ctx context.Context, userID string) error
}

type ProfileCacher interface {
	ProfileCacheReader
	ProfileCacheWriter
	ProfileInvalidator
}

type OTPCacher interface {
	SetOTP(ctx context.Context, phone, purpose string, code string, ttl time.Duration) error
	GetOTP(ctx context.Context, phone, purpose string) (string, error)
	DeleteOTP(ctx context.Context, phone, purpose string) error
}
