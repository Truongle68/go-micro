package domain

import "time"

type CredentialType string

const (
	CredentialTypePhone  CredentialType = "phone"
	CredentialTypeEmail  CredentialType = "email"
	CredentialTypeGoogle CredentialType = "google"
)

type UserCredential struct {
	ID         string         `json:"id"`
	UserID     string         `json:"user_id"`
	Type       CredentialType `json:"type"`
	Identifier string         `json:"identifier"`
	SecretHash string         `json:"-"`
	IsVerified bool           `json:"is_verified"`
	IsPrimary  bool           `json:"is_primary"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (c *UserCredential) CheckPassword(pass string) bool {
	return checkPassword(pass, c.SecretHash)
}

func NewCredential(userID string, credType CredentialType, identifier, plainPassword string, isVerified, isPrimary bool) (*UserCredential, error) {
	var hash string
	if plainPassword != "" {
		var err error
		hash, err = ValidateAndHashPassword(plainPassword)
		if err != nil {
			return nil, err
		}
	}
	c := baseUserCredential(
		userID,
		credType,
		identifier,
		hash,
		isVerified,
		isPrimary,
	)
	return c, nil
}

func NewCredentialWithHash(userID string, credType CredentialType, identifier, hash string, isVerified, isPrimary bool) *UserCredential {
	return baseUserCredential(userID, credType, identifier, hash, isVerified, isPrimary)
}

func baseUserCredential(userID string, credType CredentialType, identifier, hash string, isVerified, isPrimary bool) *UserCredential {

	now := time.Now()
	return &UserCredential{
		UserID:     userID,
		Type:       credType,
		Identifier: identifier,
		SecretHash: hash,
		IsVerified: isVerified,
		IsPrimary:  isPrimary,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
