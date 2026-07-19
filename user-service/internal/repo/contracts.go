package repo

import (
	"context"
	"user-service/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User, cred *domain.UserCredential, profile *domain.Profile) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByIdentifier(ctx context.Context, identifier string) (*domain.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByIdentifier(ctx context.Context, identifier string) (bool, error)

	// Credentials
	FindCredentialByIdentifier(ctx context.Context, identifier string) (*domain.UserCredential, error)
	FindCredentialsByUserID(ctx context.Context, userID string) ([]*domain.UserCredential, error)
	UpdateCredential(ctx context.Context, cred *domain.UserCredential) error
	SaveCredential(ctx context.Context, cred *domain.UserCredential) error

	// Profiles
	FindProfileByUserID(ctx context.Context, userID string) (*domain.Profile, error)
	UpdateProfile(ctx context.Context, profile *domain.Profile) error

	// Updates
	UpdatePassword(ctx context.Context, userID string, passwordHash string) error
	Update(ctx context.Context, user *domain.User) error

	// Addresses
	SaveAddress(ctx context.Context, address *domain.Address) error
	FindAddressByID(ctx context.Context, id string) (*domain.Address, error)
	FindAddressesByUserID(ctx context.Context, userID string) ([]*domain.Address, error)
	UpdateAddress(ctx context.Context, address *domain.Address) error
	SetDefaultAddress(ctx context.Context, userID string, addressID string) error
	DeleteAddress(ctx context.Context, userID string, addressID string) error
}
