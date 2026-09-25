package client

import "context"

type UserProfileDTO struct {
	UserID          string
	Email           string
	Phone           string
	FullName        string
	Status          string
	IsEmailVerified bool
}

type UserClient interface {
	GetProfile(ctx context.Context, userID string) (*UserProfileDTO, error)
}
