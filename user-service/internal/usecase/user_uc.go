package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repo"

	"github.com/redis/go-redis/v9"
)

type UserUC struct {
	repo        repo.UserRepository
	redisClient *redis.Client
}

func NewUserUC(repo repo.UserRepository, redisClient *redis.Client) *UserUC {
	return &UserUC{
		repo:        repo,
		redisClient: redisClient,
	}
}

var _ UserUsecase = (*UserUC)(nil)

func (uc *UserUC) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	cacheKey := fmt.Sprintf("user:profile:%s", id)

	// getting from cache
	val, err := uc.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var user domain.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			return &user, nil
		}
	}

	// fetch from DB
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// cache in Redis
	if data, err := json.Marshal(user); err == nil {
		uc.redisClient.Set(ctx, cacheKey, data, 15*time.Minute)
	}

	return user, nil
}

func (uc *UserUC) UpdateProfile(ctx context.Context, id string, fullName string, phone string) (*domain.User, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = fullName
	user.Phone = phone
	user.UpdatedAt = time.Now()

	err = uc.repo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	// invalidate cache
	cacheKey := fmt.Sprintf("user:profile:%s", id)
	uc.redisClient.Del(ctx, cacheKey)

	return user, nil
}
