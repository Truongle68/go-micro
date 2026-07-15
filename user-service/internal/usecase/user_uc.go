package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repo"
	"user-service/pkg/postgres"
	"user-service/pkg/redis"
)

type UserUC struct {
	repo       repo.UserRepository
	profile    redis.ProfileCacher
	transactor postgres.Transactor
}

func NewUserUC(repo repo.UserRepository, profile redis.ProfileCacher, transactor postgres.Transactor) *UserUC {
	return &UserUC{
		repo:       repo,
		profile:    profile,
		transactor: transactor,
	}
}

var _ User = (*UserUC)(nil)

func (uc *UserUC) GetProfile(ctx context.Context, id string) (*UserProfileDTO, error) {
	// get from cache
	val, err := uc.profile.GetProfile(ctx, id)
	if err == nil {
		var dto UserProfileDTO
		if err := json.Unmarshal([]byte(val), &dto); err == nil {
			return &dto, nil
		}
	}

	// fetch user from DB
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// fetch profile from DB
	profile, err := uc.repo.FindProfileByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	// fetch credentials to get email & phone
	creds, err := uc.repo.FindCredentialsByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	var email, phone string
	for _, c := range creds {
		switch c.Type {
		case domain.CredentialTypePhone:
			phone = c.Identifier
		case domain.CredentialTypeEmail:
			email = c.Identifier
		}
	}

	dto := &UserProfileDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     email,
		Phone:     phone,
		FullName:  profile.FullName,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// cache in Redis
	if data, err := json.Marshal(dto); err == nil {
		uc.profile.CacheProfile(ctx, id, data, 15*time.Minute)
	}

	return dto, nil
}

func (uc *UserUC) UpdateProfile(ctx context.Context, id string, fullName string, phone string) (*UserProfileDTO, error) {
	err := uc.transactor.WithTransaction(ctx, func(ctx context.Context) error {

		// update profile fullname
		profile, err := uc.repo.FindProfileByUserID(ctx, id)
		if err != nil {
			return err
		}

		profile.FullName = fullName
		profile.UpdatedAt = time.Now()
		err = uc.repo.UpdateProfile(ctx, profile)
		if err != nil {
			return fmt.Errorf("updating user profile: %w", err)
		}

		// update primary phone credential
		creds, err := uc.repo.FindCredentialsByUserID(ctx, id)
		if err != nil {
			return err
		}

		var phoneCred *domain.UserCredential
		for _, c := range creds {
			if c.Type == domain.CredentialTypePhone {
				phoneCred = c
				break
			}
		}

		if phoneCred != nil {
			phoneCred.Identifier = phone
			phoneCred.UpdatedAt = time.Now()
			err = uc.repo.UpdateCredential(ctx, phoneCred)
			if err != nil {
				return fmt.Errorf("updating user phone credential: %w", err)
			}
		}

		// update general user info (updated_at)
		user, err := uc.repo.FindByID(ctx, id)
		if err == nil {
			user.UpdatedAt = time.Now()
			_ = uc.repo.Update(ctx, user)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// invalidate cache
	uc.profile.InvalidateProfile(ctx, id)

	// return updated profile details
	return uc.GetProfile(ctx, id)
}
