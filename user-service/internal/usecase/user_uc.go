package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repo"
	"user-service/pkg/jwt"
	"user-service/pkg/postgres"
	"user-service/pkg/redis"
	"user-service/worker"

	"github.com/TruongLe68/go-micro/pkg/logger"
)

type UserUC struct {
	repo       repo.UserRepository
	cache      redis.IdentityCacher
	transactor postgres.Transactor
	tokens     jwt.TokenService
	email      worker.EmailDispatcher
	logger     logger.Interface
	baseURL    string
}

func NewUserUC(
	repo repo.UserRepository,
	tokens jwt.TokenService,
	cache redis.IdentityCacher,
	transactor postgres.Transactor,
	email worker.EmailDispatcher,
	logger logger.Interface,
	baseURL string,
) *UserUC {
	return &UserUC{
		repo:       repo,
		cache:      cache,
		transactor: transactor,
		tokens:     tokens,
		email:      email,
		logger:     logger,
		baseURL:    baseURL,
	}
}

var _ User = (*UserUC)(nil)

func (uc *UserUC) GetProfile(ctx context.Context, id string) (*UserProfileDTO, error) {
	// get from cache
	val, err := uc.cache.GetProfile(ctx, id)
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
		Gender:    profile.Gender,
		DOB:       profile.DOB,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// cache in Redis
	if data, err := json.Marshal(dto); err == nil {
		uc.cache.CacheProfile(ctx, id, data, 15*time.Minute)
	}

	return dto, nil
}

func (uc *UserUC) UpdateProfile(ctx context.Context, in UpdatedProfileInput) (*UserProfileDTO, error) {
	err := uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {

		// update profile
		profile, err := uc.repo.FindProfileByUserID(txCtx, in.UserID)
		if err != nil {
			return err
		}

		if err := profile.ApplyUpdate(domain.UpdateProfileParams{
			FullName: in.FullName,
			Gender:   in.Gender,
			Dob:      in.Dob,
		}); err != nil {
			return err
		}

		err = uc.repo.UpdateProfile(txCtx, profile)
		if err != nil {
			return fmt.Errorf("updating user profile: %w", err)
		}

		// update general user info (updated_at)
		user, err := uc.repo.FindByID(txCtx, in.UserID)
		if err != nil {
			return fmt.Errorf("finding user: %w", err)
		}
		user.UpdatedAt = time.Now()
		return uc.repo.Update(txCtx, user)
	})

	if err != nil {
		return nil, err
	}

	// invalidate cache
	uc.cache.InvalidateProfile(ctx, in.UserID)

	// return updated profile details
	return uc.GetProfile(ctx, in.UserID)
}

func (uc *UserUC) RequestChangeEmailOTP(ctx context.Context, userID string) error {
	creds, err := uc.repo.FindCredentialsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var currentEmail string
	for _, c := range creds {
		if c.Type == domain.CredentialTypeEmail {
			currentEmail = c.Identifier
			break
		}
	}

	if currentEmail == "" {
		return domain.ErrEmailNotSet
	}

	code, err := redis.GenOTPCode()
	if err != nil {
		return err
	}

	uc.logger.Info(" [OTP DEBUG] Generated Change Email OTP for User %s: %s", userID, code)

	err = uc.cache.SetOTP(ctx, userID, domain.VerifyPurposeChangeEmail, code, 5*time.Minute)
	if err != nil {
		return fmt.Errorf("caching OTP: %w", err)
	}

	uc.email.Dispatch(worker.EmailJob{
		Template: "change_email",
		Subject:  "Verify your request to change email address",
		To:       currentEmail,
		Data: struct{ Code string }{
			Code: code,
		},
	})

	return nil
}

func (uc *UserUC) VerifyChangeEmailOTP(ctx context.Context, in VerifyChangeEmailInput) (string, error) {
	cachedCode, err := uc.cache.GetOTP(ctx, in.UserID, domain.VerifyPurposeChangeEmail)
	if err != nil {
		return "", domain.ErrOTPExpired
	}

	if cachedCode != in.Code {
		return "", domain.ErrInvalidOTP
	}

	_ = uc.cache.DeleteOTP(ctx, in.UserID, domain.VerifyPurposeChangeEmail)

	token, err := uc.tokens.GenerateChangeEmailToken(in.UserID)
	if err != nil {
		return "", fmt.Errorf("generating change email token: %w", err)
	}

	return token, nil
}

func (uc *UserUC) CompleteChangeEmail(ctx context.Context, in CompleteChangeEmailInput) error {
	claims, err := uc.tokens.VerifyChangeEmailToken(in.Token)
	if err != nil {
		return domain.ErrInvalidToken
	}

	if claims.UserID != in.UserID {
		return domain.ErrInvalidToken
	}

	// Find credentials to check if new email is different and unique
	creds, err := uc.repo.FindCredentialsByUserID(ctx, in.UserID)
	if err != nil {
		return err
	}

	var emailCred *domain.UserCredential
	var phoneCred *domain.UserCredential
	for _, c := range creds {
		switch c.Type {
		case domain.CredentialTypeEmail:
			emailCred = c
		case domain.CredentialTypePhone:
			phoneCred = c
		}
	}

	if emailCred != nil && emailCred.Identifier == in.NewEmail {
		return domain.ErrSameEmail
	}

	// Check email uniqueness
	exists, err := uc.repo.ExistsByIdentifier(ctx, in.NewEmail)
	if err != nil {
		return fmt.Errorf("checking email existence: %w", err)
	}
	if exists {
		return domain.ErrEmailAlreadyExists
	}

	err = uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		if emailCred != nil {
			emailCred.Identifier = in.NewEmail
			emailCred.UpdatedAt = time.Now()
			err = uc.repo.UpdateCredential(txCtx, emailCred)
			if err != nil {
				return fmt.Errorf("updating user email credential: %w", err)
			}
		} else {
			var secretHash string
			if phoneCred != nil {
				secretHash = phoneCred.SecretHash
			}
			newEmailCred := domain.NewCredentialWithHash(
				in.UserID,
				domain.CredentialTypeEmail,
				in.NewEmail,
				secretHash,
				true,
				false,
			)
			err = uc.repo.SaveCredential(txCtx, newEmailCred)
			if err != nil {
				return fmt.Errorf("saving user email credential: %w", err)
			}
		}

		// update general user info (updated_at)
		user, err := uc.repo.FindByID(txCtx, in.UserID)
		if err == nil {
			user.UpdatedAt = time.Now()
			_ = uc.repo.Update(txCtx, user)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// invalidate cache
	uc.cache.InvalidateProfile(ctx, in.UserID)
	return nil
}

func (uc *UserUC) RequestChangePhoneLink(ctx context.Context, userID string) error {
	creds, err := uc.repo.FindCredentialsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var currentEmail string
	for _, c := range creds {
		if c.Type == domain.CredentialTypeEmail {
			currentEmail = c.Identifier
			break
		}
	}

	if currentEmail == "" {
		return domain.ErrEmailNotSet
	}

	token, err := uc.tokens.GenerateChangePhoneToken(userID)
	if err != nil {
		return fmt.Errorf("generating change phone token: %w", err)
	}

	resetLink := fmt.Sprintf("%s/profile/change-phone?token=%s", uc.baseURL, token)

	uc.email.Dispatch(worker.EmailJob{
		Template: "change_phone",
		Subject:  "Verify your request to change phone number",
		To:       currentEmail,
		Data: struct{ ResetLink string }{
			ResetLink: resetLink,
		},
	})

	return nil
}

func (uc *UserUC) CompleteChangePhone(ctx context.Context, in CompleteChangePhoneInput) error {
	claims, err := uc.tokens.VerifyChangePhoneToken(in.Token)
	if err != nil {
		return domain.ErrInvalidToken
	}

	userID := claims.UserID

	creds, err := uc.repo.FindCredentialsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var phoneCred *domain.UserCredential
	var hash string
	for _, c := range creds {
		if hash == "" {
			hash = c.SecretHash
		}
		if c.Type == domain.CredentialTypePhone {
			phoneCred = c
			break
		}
	}

	if phoneCred != nil && phoneCred.Identifier == in.NewPhone {
		return domain.ErrSamePhone
	}

	// Check uniqueness
	exists, err := uc.repo.ExistsByIdentifier(ctx, in.NewPhone)
	if err != nil {
		return fmt.Errorf("checking phone existence: %w", err)
	}
	if exists {
		return domain.ErrPhoneAlreadyExists
	}

	err = uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		if phoneCred != nil {
			phoneCred.Identifier = in.NewPhone
			phoneCred.UpdatedAt = time.Now()
			err = uc.repo.UpdateCredential(txCtx, phoneCred)
			if err != nil {
				return fmt.Errorf("updating user phone credential: %w", err)
			}
		} else {
			newPhoneCred := domain.NewCredentialWithHash(
				userID,
				domain.CredentialTypeEmail,
				in.NewPhone,
				hash,
				true,
				false,
			)

			err = uc.repo.SaveCredential(txCtx, newPhoneCred)
			if err != nil {
				return fmt.Errorf("saving user phone credential: %w", err)
			}
		}

		// update general user info (updated_at)
		user, err := uc.repo.FindByID(txCtx, userID)
		if err == nil {
			user.UpdatedAt = time.Now()
			_ = uc.repo.Update(txCtx, user)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// invalidate cache
	uc.cache.InvalidateProfile(ctx, userID)
	return nil
}

func (uc *UserUC) GetAddressList(ctx context.Context, userID string) ([]*AddressDTO, error) {
	_, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	addresses, err := uc.repo.FindAddressesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*AddressDTO, len(addresses))
	for i, a := range addresses {
		dtos[i] = &AddressDTO{
			ID:          a.ID,
			UserID:      a.UserID,
			Label:       a.Label,
			AddressLine: a.AddressLine,
			Ward:        a.Ward,
			District:    a.District,
			City:        a.City,
			Lat:         a.Lat,
			Lng:         a.Lng,
			IsDefault:   a.IsDefault,
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		}
	}
	return dtos, nil
}

func (uc *UserUC) CreateAddress(ctx context.Context, in CreateAddressInput) (*AddressDTO, error) {
	_, err := uc.repo.FindByID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	existing, err := uc.repo.FindAddressesByUserID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}
	isFirst := len(existing) == 0

	opts := []domain.AddressOption{
		domain.Coordinates(in.Lat, in.Lng),
	}
	if in.Ward != "" {
		opts = append(opts, func(a *domain.Address) { a.Ward = in.Ward })
	}
	if in.District != "" {
		opts = append(opts, func(a *domain.Address) { a.District = in.District })
	}

	address, err := domain.NewAddress(in.UserID, in.AddressLine, in.City, in.Label, opts...)
	if err != nil {
		return nil, err
	}

	if isFirst {
		address.IsDefault = true
	}

	err = uc.repo.SaveAddress(ctx, address)
	if err != nil {
		return nil, err
	}

	return &AddressDTO{
		ID:          address.ID,
		UserID:      address.UserID,
		Label:       address.Label,
		AddressLine: address.AddressLine,
		Ward:        address.Ward,
		District:    address.District,
		City:        address.City,
		Lat:         address.Lat,
		Lng:         address.Lng,
		IsDefault:   address.IsDefault,
		CreatedAt:   address.CreatedAt,
		UpdatedAt:   address.UpdatedAt,
	}, nil
}

func (uc *UserUC) UpdateAddress(ctx context.Context, in UpdateAddressInput) (*AddressDTO, error) {
	address, err := uc.repo.FindAddressByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if address.UserID != in.UserID {
		return nil, domain.ErrUnauthorizedUser
	}

	params := domain.UpdateAddressParams{
		Label:       in.Label,
		AddressLine: in.AddressLine,
		Ward:        in.Ward,
		District:    in.District,
		City:        in.City,
		Lat:         in.Lat,
		Lng:         in.Lng,
	}

	if err := address.ApplyUpdate(params); err != nil {
		return nil, err
	}

	err = uc.repo.UpdateAddress(ctx, address)
	if err != nil {
		return nil, err
	}

	return &AddressDTO{
		ID:          address.ID,
		UserID:      address.UserID,
		Label:       address.Label,
		AddressLine: address.AddressLine,
		Ward:        address.Ward,
		District:    address.District,
		City:        address.City,
		Lat:         address.Lat,
		Lng:         address.Lng,
		IsDefault:   address.IsDefault,
		CreatedAt:   address.CreatedAt,
		UpdatedAt:   address.UpdatedAt,
	}, nil
}

func (uc *UserUC) SetDefaultAddress(ctx context.Context, userID string, addressID string) error {
	address, err := uc.repo.FindAddressByID(ctx, addressID)
	if err != nil {
		return err
	}
	if address.UserID != userID {
		return domain.ErrUnauthorizedUser
	}

	return uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		return uc.repo.SetDefaultAddress(txCtx, userID, addressID)
	})
}

func (uc *UserUC) DeleteAddress(ctx context.Context, userID string, addressID string) error {
	address, err := uc.repo.FindAddressByID(ctx, addressID)
	if err != nil {
		return err
	}
	if address.UserID != userID {
		return domain.ErrUnauthorizedUser
	}

	return uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		err := uc.repo.DeleteAddress(txCtx, userID, addressID)
		if err != nil {
			return err
		}

		if address.IsDefault {
			remaining, err := uc.repo.FindAddressesByUserID(txCtx, userID)
			if err != nil {
				return err
			}
			if len(remaining) > 0 {
				err = uc.repo.SetDefaultAddress(txCtx, userID, remaining[0].ID)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
