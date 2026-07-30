package usecase

import (
	"context"
	"encoding/json"
	"errors"
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
	var isEmailVerified bool
	for _, c := range creds {
		switch c.Type {
		case domain.CredentialTypePhone:
			phone = c.Identifier
		case domain.CredentialTypeEmail:
			email = c.Identifier
			isEmailVerified = c.IsVerified
		}
	}

	dto := &UserProfileDTO{
		ID:              user.ID,
		Username:        user.Username,
		Email:           email,
		IsEmailVerified: isEmailVerified,
		Phone:           phone,
		FullName:        profile.FullName,
		Gender:          profile.Gender,
		DOB:             profile.DOB,
		Role:            user.Role,
		Status:          user.Status,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}

	// cache in Redis
	if data, err := json.Marshal(dto); err == nil {
		uc.cache.CacheProfile(ctx, id, data, 15*time.Minute)
	}

	return dto, nil
}

func (uc *UserUC) UpdateProfile(ctx context.Context, in UpdatedProfileInput) (*UserProfileDTO, error) {
	err := uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {

		// check existing credentials
		creds, err := uc.repo.FindCredentialsByUserID(txCtx, in.UserID)
		if err != nil {
			return err
		}

		var emailCred *domain.UserCredential
		var phoneCred *domain.UserCredential
		for _, c := range creds {
			switch c.Type {
			case domain.CredentialTypeEmail:
				emailCred = &c
			case domain.CredentialTypePhone:
				phoneCred = &c
			}
		}

		// handle email update
		if in.Email != nil && *in.Email != "" {
			newEmail := *in.Email
			if emailCred != nil && emailCred.IsVerified {
				if emailCred.Identifier != newEmail {
					return domain.ErrVerifiedEmailCannotBeUpdatedDirectly
				}
			} else {
				if emailCred == nil || emailCred.Identifier != newEmail {
					exists, err := uc.repo.ExistsByIdentifier(txCtx, newEmail)
					if err != nil {
						return fmt.Errorf("checking email existence: %w", err)
					}
					if exists {
						return domain.ErrEmailAlreadyExists
					}

					if emailCred != nil {
						emailCred.Identifier = newEmail
						emailCred.IsVerified = false
						emailCred.UpdatedAt = time.Now()
						if err := uc.repo.UpdateCredential(txCtx, emailCred); err != nil {
							return fmt.Errorf("updating email credential: %w", err)
						}
					} else {
						var secretHash string
						if phoneCred != nil {
							secretHash = phoneCred.SecretHash
						}
						newCred := domain.NewCredentialWithHash(
							in.UserID,
							domain.CredentialTypeEmail,
							newEmail,
							secretHash,
							false,
							false,
						)
						if err := uc.repo.SaveCredential(txCtx, newCred); err != nil {
							return fmt.Errorf("saving email credential: %w", err)
						}
					}
				}
			}
		}

		// update profile fields
		profile, err := uc.repo.FindProfileByUserID(txCtx, in.UserID)
		if err != nil {
			return err
		}

		if err := profile.ApplyUpdate(domain.UpdateProfileParams{
			FullName: in.FullName,
			Gender:   in.Gender,
			Dob:      in.Dob,
		}); err != nil && !errors.Is(err, domain.ErrNoFieldsToUpdate) {
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

func (uc *UserUC) RequestEmailLink(ctx context.Context, in RequestEmailLinkInput) error {
	policy, ok := domain.GetEmailLinkPolicy(in.Purpose)
	if !ok {
		return domain.ErrInvalidOperation
	}

	if in.Email == "" {
		return domain.ErrEmailNotSet
	}

	targetEmail := in.Email
	switch in.Purpose {
	case domain.EmailLinkPurposeVerifyCurrent,
		domain.EmailLinkPurposeVerifyNew:
	default:
		return domain.ErrInvalidOperation
	}

	token, err := uc.tokens.GenerateEmailLinkToken(in.ActorUserID, targetEmail, in.Purpose, policy.TokenTTL)
	if err != nil {
		return fmt.Errorf("generating email link token: %w", err)
	}

	link := fmt.Sprintf("%s/users/verify-email/confirm?token=%s", uc.baseURL, token)

	uc.email.Dispatch(worker.EmailJob{
		Template: "verify_email",
		Subject:  "Verify Your Email Address",
		To:       targetEmail,
		Data: struct {
			Link string
			TTL  time.Duration
		}{
			Link: link,
			TTL:  policy.TokenTTL,
		},
	})

	return nil
}

func (uc *UserUC) ConfirmEmailLink(ctx context.Context, token string) (*ConfirmEmailLinkOutput, error) {
	claims, err := uc.tokens.VerifyEmailLinkToken(token)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	purpose := domain.EmailLinkPurpose(claims.Purpose)

	switch purpose {
	case domain.EmailLinkPurposeVerifyNew:
		if err := uc.attachVerifiedEmail(ctx, claims.UserID, claims.Email); err != nil {
			return nil, err
		}
		return &ConfirmEmailLinkOutput{
			Purpose: purpose,
		}, nil

	case domain.EmailLinkPurposeVerifyCurrent:
		changeToken, err := uc.tokens.GenerateChangeEmailToken(claims.UserID)
		if err != nil {
			return nil, fmt.Errorf("generating change email token: %w", err)
		}
		return &ConfirmEmailLinkOutput{
			Purpose:          purpose,
			ChangeEmailToken: changeToken,
		}, nil

	default:
		return nil, domain.ErrInvalidOperation
	}
}

func (uc *UserUC) sendAccountOTP(ctx context.Context, in RequestOTPInput) error {
	p, ok := domain.GetVerifiedOTPPolicy(in.Purpose)
	if !ok {
		return domain.ErrInvalidOperation
	}

	switch in.Purpose {
	case domain.VerifyPurposeChangeEmail:
		if in.ChangeToken == "" {
			return domain.ErrCurrentEmailNotVerified
		}
		claims, err := uc.tokens.VerifyChangeEmailToken(in.ChangeToken)
		if err != nil || claims.UserID != in.ActorUserID {
			return domain.ErrCurrentEmailNotVerified
		}

		if in.Identifier == "" {
			return domain.ErrEmailRequired
		}

		exists, err := uc.repo.ExistsByIdentifier(ctx, in.Identifier)
		if err != nil {
			return fmt.Errorf("checking email existence: %w", err)
		}
		if exists {
			return domain.ErrEmailAlreadyExists
		}

		code, err := redis.GenOTPCode()
		if err != nil {
			return err
		}

		if err := uc.cache.SetOTP(ctx, in.ActorUserID, in.Purpose, code, p.OTPTTL); err != nil {
			return fmt.Errorf("caching OTP: %w", err)
		}

		uc.email.Dispatch(worker.EmailJob{
			Template: "change_email_otp",
			Subject:  "OTP code for changed email verification",
			To:       in.Identifier,
			Data: struct {
				Code string
				TTL  time.Duration
			}{
				Code: code,
				TTL:  p.OTPTTL,
			},
		})
		return nil
	case domain.VerifyPurposeChangePhone:
		if in.ChangeToken == "" {
			return domain.ErrCurrentPhoneNotVerified
		}

		claims, err := uc.tokens.VerifyChangePhoneToken(in.ChangeToken)
		if err != nil || claims.UserID != in.ActorUserID {
			return domain.ErrCurrentPhoneNotVerified
		}

		if in.Identifier == "" {
			return domain.ErrEmptyPhone
		}

		exists, err := uc.repo.ExistsByIdentifier(ctx, in.Identifier)
		if err != nil {
			return fmt.Errorf("checking phone existence: %w", err)
		}
		if exists {
			return domain.ErrPhoneAlreadyExists
		}

		code, err := redis.GenOTPCode()
		if err != nil {
			return err
		}

		if err := uc.cache.SetOTP(ctx, in.ActorUserID, in.Purpose, code, p.OTPTTL); err != nil {
			return fmt.Errorf("caching OTP: %w", err)
		}

		uc.logger.Info(" [OTP DEBUG] Generated Change Phone OTP for User %s (%s): %s", in.ActorUserID, in.Identifier, code)
		return nil
	case domain.VerifyPurposeVerifyPhone:
		creds, err := uc.repo.FindCredentialByTypeAndUserID(ctx, domain.CredentialTypePhone, in.ActorUserID)
		if err != nil {
			return fmt.Errorf("finding credentials by type and userID: %w", err)
		}

		code, err := redis.GenOTPCode()
		if err != nil {
			return err
		}

		if err := uc.cache.SetOTP(ctx, in.ActorUserID, in.Purpose, code, p.OTPTTL); err != nil {
			return fmt.Errorf("caching OTP: %w", err)
		}

		uc.logger.Info(" [OTP DEBUG] Generated Verified Phone OTP for User %s (%s): %s", in.ActorUserID, creds.Identifier, code)
		return nil
	}

	return domain.ErrInvalidOperation
}

func (uc *UserUC) SendChangeEmailOTP(ctx context.Context, in RequestChangeEmailOTPInput) error {
	return uc.sendAccountOTP(ctx, RequestOTPInput{
		Identifier:  in.Identifier,
		Purpose:     domain.VerifyPurposeChangeEmail,
		ChangeToken: in.ChangeEmailToken,
		ActorUserID: in.ActorUserID,
	})
}

func (uc *UserUC) SendChangePhoneOTP(ctx context.Context, in RequestChangePhoneOTPInput) error {
	return uc.sendAccountOTP(ctx, RequestOTPInput{
		Identifier:  in.Identifier,
		Purpose:     domain.VerifyPurposeChangePhone,
		ChangeToken: in.ChangePhoneToken,
		ActorUserID: in.ActorUserID,
	})
}

func (uc *UserUC) SendPhoneVerificationOTP(ctx context.Context, userID string) error {
	return uc.sendAccountOTP(ctx, RequestOTPInput{
		Purpose:     domain.VerifyPurposeVerifyPhone,
		ActorUserID: userID,
	})
}

func (uc *UserUC) verifyAccountOTP(ctx context.Context, in VerifyOTPInput) (string, error) {
	cachedCode, err := uc.cache.GetOTP(ctx, in.ActorUserID, in.Purpose)
	if err != nil {
		return "", domain.ErrOTPExpired
	}
	if cachedCode != in.Code {
		return "", domain.ErrInvalidOTP
	}
	_ = uc.cache.DeleteOTP(ctx, in.ActorUserID, in.Purpose)

	switch in.Purpose {
	case domain.VerifyPurposeChangeEmail:
		if err := uc.attachVerifiedEmail(ctx, in.ActorUserID, in.Identifier); err != nil {
			return "", err
		}
		return "", nil
	case domain.VerifyPurposeChangePhone:
		if err := uc.attachVerifiedPhone(ctx, in.ActorUserID, in.Identifier); err != nil {
			return "", err
		}
		return "", nil
	case domain.VerifyPurposeVerifyPhone:
		token, err := uc.tokens.GenerateChangePhoneToken(in.ActorUserID)
		if err != nil {
			return "", fmt.Errorf("generating change phone token: %w", err)
		}
		return token, nil
	}

	return "", domain.ErrInvalidOperation
}

func (uc *UserUC) VerifyChangeEmailOTP(ctx context.Context, in VerifyChangeEmailOTPInput) error {
	_, err := uc.verifyAccountOTP(ctx, VerifyOTPInput{
		Identifier:  in.Identifier,
		Code:        in.Code,
		Purpose:     domain.VerifyPurposeChangeEmail,
		ActorUserID: in.ActorUserID,
	})
	return err
}

func (uc *UserUC) VerifyChangePhoneOTP(ctx context.Context, in VerifyChangePhoneOTPInput) error {
	_, err := uc.verifyAccountOTP(ctx, VerifyOTPInput{
		Identifier:  in.Identifier,
		Code:        in.Code,
		Purpose:     domain.VerifyPurposeChangePhone,
		ActorUserID: in.ActorUserID,
	})
	return err
}

func (uc *UserUC) VerifyPhoneVerificationOTP(ctx context.Context, in VerifyPhoneVerificationOTPInput) (string, error) {
	return uc.verifyAccountOTP(ctx, VerifyOTPInput{
		Code:        in.Code,
		Purpose:     domain.VerifyPurposeVerifyPhone,
		ActorUserID: in.ActorUserID,
	})
}

func (uc *UserUC) attachVerifiedEmail(ctx context.Context, userID, email string) error {
	creds, err := uc.repo.FindCredentialsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var emailCred *domain.UserCredential
	var phoneCred *domain.UserCredential
	for _, c := range creds {
		switch c.Type {
		case domain.CredentialTypeEmail:
			emailCred = &c
		case domain.CredentialTypePhone:
			phoneCred = &c
		}
	}

	return uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		if emailCred != nil {
			emailCred.Identifier = email
			emailCred.IsVerified = true
			emailCred.UpdatedAt = time.Now()
			if err := uc.repo.UpdateCredential(txCtx, emailCred); err != nil {
				return fmt.Errorf("updating email credential: %w", err)
			}
		} else {
			var secretHash string
			if phoneCred != nil {
				secretHash = phoneCred.SecretHash
			}
			newCred := domain.NewCredentialWithHash(userID, domain.CredentialTypeEmail, email, secretHash, true, false)
			if err := uc.repo.SaveCredential(txCtx, newCred); err != nil {
				return fmt.Errorf("saving email credential: %w", err)
			}
		}

		user, err := uc.repo.FindByID(txCtx, userID)
		if err == nil {
			user.UpdatedAt = time.Now()
			_ = uc.repo.Update(txCtx, user)
		}
		uc.cache.InvalidateProfile(ctx, userID)
		return nil
	})
}

func (uc *UserUC) attachVerifiedPhone(ctx context.Context, userID, phone string) error {
	creds, err := uc.repo.FindCredentialsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var phoneCred *domain.UserCredential
	var emailCred *domain.UserCredential
	for _, c := range creds {
		switch c.Type {
		case domain.CredentialTypePhone:
			phoneCred = &c
		case domain.CredentialTypeEmail:
			emailCred = &c
		}
	}

	return uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		if phoneCred != nil {
			phoneCred.Identifier = phone
			phoneCred.IsVerified = true
			phoneCred.UpdatedAt = time.Now()
			if err := uc.repo.UpdateCredential(txCtx, phoneCred); err != nil {
				return fmt.Errorf("updating phone credential: %w", err)
			}
		} else {
			var secretHash string
			if emailCred != nil {
				secretHash = emailCred.SecretHash
			}
			newCred := domain.NewCredentialWithHash(userID, domain.CredentialTypePhone, phone, secretHash, true, false)
			if err := uc.repo.SaveCredential(txCtx, newCred); err != nil {
				return fmt.Errorf("saving phone credential: %w", err)
			}
		}

		user, err := uc.repo.FindByID(txCtx, userID)
		if err == nil {
			user.UpdatedAt = time.Now()
			_ = uc.repo.Update(txCtx, user)
		}
		uc.cache.InvalidateProfile(ctx, userID)
		return nil
	})
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
