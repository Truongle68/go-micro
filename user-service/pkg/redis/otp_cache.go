package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPCache struct {
	client *redis.Client
}

func NewOTPCache(client *redis.Client) *OTPCache {
	return &OTPCache{
		client: client,
	}
}

var _ OTPCacher = (*OTPCache)(nil)

func (c *OTPCache) SetOTP(ctx context.Context, phone, purpose string, code string, ttl time.Duration) error {
	return c.client.Set(ctx, otpKey(phone, purpose), code, ttl).Err()
}

func (c *OTPCache) GetOTP(ctx context.Context, phone, purpose string) (string, error) {
	return c.client.Get(ctx, otpKey(phone, purpose)).Result()
}

func (c *OTPCache) DeleteOTP(ctx context.Context, phone, purpose string) error {
	return c.client.Del(ctx, otpKey(phone, purpose)).Err()
}

func otpKey(phone, purpose string) string {
	return fmt.Sprintf("otp:%s:%s", purpose, phone)
}
