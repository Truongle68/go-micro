package redis

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
	"user-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

type OTPPurpose string

type OTPCache struct {
	client *redis.Client
}

func NewOTPCache(client *redis.Client) *OTPCache {
	return &OTPCache{
		client: client,
	}
}

var _ OTPCacher = (*OTPCache)(nil)

func (c *OTPCache) SetOTP(ctx context.Context, phone string, purpose domain.VerifyPurpose, code string, ttl time.Duration) error {
	return c.client.Set(ctx, otpKey(phone, purpose), code, ttl).Err()
}

func (c *OTPCache) GetOTP(ctx context.Context, phone string, purpose domain.VerifyPurpose) (string, error) {
	return c.client.Get(ctx, otpKey(phone, purpose)).Result()
}

func (c *OTPCache) DeleteOTP(ctx context.Context, phone string, purpose domain.VerifyPurpose) error {
	return c.client.Del(ctx, otpKey(phone, purpose)).Err()
}

func otpKey(phone string, purpose domain.VerifyPurpose) string {
	return fmt.Sprintf("otp:%s:%s", string(purpose), phone)
}

func GenOTPCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("generating random OTP: %w", err)
	}
	code := fmt.Sprintf("%06d", n.Int64())
	return code, nil
}
