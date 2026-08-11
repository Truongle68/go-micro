package usecase

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repo"
	"user-service/pkg/jwt"
	"user-service/pkg/postgres"
	"user-service/worker"

	redismanager "github.com/GoProOrg/core-go-pkg/redismanager/identity"
	"github.com/TruongLe68/go-micro/pkg/logger"
)

type AuthUC struct {
	userRepo      repo.UserRepository
	tokens        jwt.TokenService
	cache         redismanager.IdentityCacher
	transactor    postgres.Transactor
	email         worker.EmailDispatcher
	otpDispatcher worker.OTPDispatcher
	logger        logger.Interface
	baseURL       string
}

func NewAuthUC(
	repo repo.UserRepository,
	tokens jwt.TokenService,
	cache redismanager.IdentityCacher,
	transactor postgres.Transactor,
	email worker.EmailDispatcher,
	otpDispatcher worker.OTPDispatcher,
	logger logger.Interface,
	baseURL string,
) *AuthUC {
	return &AuthUC{
		userRepo:      repo,
		tokens:        tokens,
		cache:         cache,
		transactor:    transactor,
		email:         email,
		otpDispatcher: otpDispatcher,
		logger:        logger,
		baseURL:       baseURL,
	}
}

var _ Auth = (*AuthUC)(nil)

func (uc *AuthUC) RequestOTP(ctx context.Context, in RequestOTPInput) error {
	p, ok := domain.GetVerifiedOTPPolicy(in.Purpose)
	if !ok {
		return domain.ErrInvalidOperation
	}

	identifier := in.Identifier
	if identifier == "" {
		return domain.ErrEmptyPhone
	}

	code, err := redismanager.GenOTPCode()
	if err != nil {
		return err
	}

	err = uc.cache.SetOTP(ctx, identifier, string(in.Purpose), code, p.OTPTTL)
	if err != nil {
		return fmt.Errorf("caching OTP: %w", err)
	}

	uc.logger.Info(" [OTP DEBUG] Generated OTP for %s (%s): %s\n", identifier, in.Purpose, code)
	uc.otpDispatcher.DispatchOTP(domain.OTPChannelSMS, identifier, code)
	uc.logger.Info("Dispatch OTP for %s (purpose: %s)", identifier, in.Purpose)

	return nil
}

func (uc *AuthUC) VerifyOTP(ctx context.Context, in VerifyOTPInput) (token string, exists bool, username string, err error) {
	identifier := in.Identifier
	if identifier == "" {
		return "", false, "", domain.ErrEmptyPhone
	}
	if in.Code == "" {
		return "", false, "", domain.ErrInvalidOTP
	}

	cachedCode, err := uc.cache.GetOTP(ctx, identifier, string(in.Purpose))
	if err != nil {
		return "", false, "", domain.ErrOTPExpired
	}

	if cachedCode != in.Code {
		return "", false, "", domain.ErrInvalidOTP
	}

	_ = uc.cache.DeleteOTP(ctx, identifier, string(in.Purpose))

	token, err = uc.tokens.GenerateVerificationToken(identifier, in.Purpose)
	if err != nil {
		return "", false, "", fmt.Errorf("generating verification token: %w", err)
	}

	// Check if user exists with this phone/identifier number
	if in.Purpose == domain.VerifyPurposeRegister {
		cred, err := uc.userRepo.FindCredentialByIdentifier(ctx, identifier)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return "", false, "", fmt.Errorf("finding credential by identifier: %w", err)
		}

		if cred != nil {
			exists = true
			user, err := uc.userRepo.FindByID(ctx, cred.UserID)
			if err != nil {
				return "", false, "", fmt.Errorf("finding user by id: %w", err)
			}
			username = user.Username
		}
	}

	return token, exists, username, nil
}

func (uc *AuthUC) CompleteRegister(ctx context.Context, in RegisterInput) (AuthOutput, error) {
	claims, err := uc.tokens.VerifyVerificationToken(in.Token)
	if err != nil {
		return AuthOutput{}, domain.ErrInvalidToken
	}

	isBlackList, err := uc.cache.IsBlacklisted(ctx, claims.JTI())
	if err != nil {
		return AuthOutput{}, fmt.Errorf("checking token blacklist: %w", err)
	}

	if isBlackList {
		return AuthOutput{}, domain.ErrInvalidToken
	}

	if claims.Purpose != string(domain.VerifyPurposeRegister) {
		return AuthOutput{}, domain.ErrInvalidToken
	}

	// check if username exists
	usernameExists, err := uc.userRepo.ExistsByUsername(ctx, in.Username)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("checking username: %w", err)
	}
	if usernameExists {
		return AuthOutput{}, domain.ErrUsernameExists
	}

	// check if phone exists
	phoneExists, err := uc.userRepo.ExistsByIdentifier(ctx, claims.Phone)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("checking phone: %w", err)
	}
	if phoneExists {
		return AuthOutput{}, domain.ErrPhoneAlreadyExists
	}

	user, err := domain.NewUser(in.Username)
	if err != nil {
		return AuthOutput{}, err
	}

	cred, err := domain.NewCredential(user.ID, domain.CredentialTypePhone, claims.Phone, in.Password, true, true)
	if err != nil {
		return AuthOutput{}, err
	}

	prof := domain.NewProfile(user.ID, in.FullName)

	err = uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		err := uc.userRepo.Save(txCtx, user, cred, prof)
		if err != nil {
			return fmt.Errorf("saving new user: %w", err)
		}
		return nil
	})

	if err != nil {
		return AuthOutput{}, fmt.Errorf("register user transaction: %w", err)
	}

	accessToken, err := uc.tokens.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, err := uc.tokens.GenerateRefreshToken(user.ID)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("generating refresh token: %w", err)
	}

	if err := uc.cache.Add(ctx, claims.JTI(), time.Until(claims.ExpiresAtTime())); err != nil {
		return AuthOutput{}, fmt.Errorf("adding to black list: %w", err)
	}

	return AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
	}, nil
}

func (uc *AuthUC) CheckUsernameAvailable(ctx context.Context, username string) (bool, error) {
	exists, err := uc.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return false, fmt.Errorf("checking username: %w", err)
	}
	return !exists, nil
}

func (uc *AuthUC) Login(ctx context.Context, in LoginInput) (AuthOutput, error) {
	var cred *domain.UserCredential
	var err error

	// finding by identifier (e.g. phone or email)
	cred, err = uc.userRepo.FindCredentialByIdentifier(ctx, in.Identifier)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return AuthOutput{}, fmt.Errorf("finding credential: %w", err)
		}
		// if not found, check if it matches a username
		user, uErr := uc.userRepo.FindByUsername(ctx, in.Identifier)
		if uErr == nil {
			// find primary credential for this user
			creds, cErr := uc.userRepo.FindCredentialsByUserID(ctx, user.ID)
			if cErr == nil {
				for i := range creds {
					if creds[i].IsPrimary {
						cred = &creds[i]
						break
					}
				}
			}
		}
	}

	if cred == nil || !cred.CheckPassword(in.Password) {
		return AuthOutput{}, domain.ErrInvalidCredentials
	}

	user, err := uc.userRepo.FindByID(ctx, cred.UserID)
	if err != nil {
		return AuthOutput{}, domain.ErrInvalidCredentials
	}

	if len(in.RequiredRoles) > 0 {
		if isAuthorized := slices.Contains(in.RequiredRoles, user.Role); !isAuthorized {
			return AuthOutput{}, domain.ErrUnauthorizedUser
		}
	}

	switch user.Status {
	case domain.UserStatusBanned:
		return AuthOutput{}, domain.ErrUserBanned
	case domain.UserStatusUnverified, domain.UserStatusVerified:
		// OK, only making order requires verified user
	default:
		return AuthOutput{}, fmt.Errorf("unknown user status: %s", user.Status)
	}

	accessToken, err := uc.tokens.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, err := uc.tokens.GenerateRefreshToken(user.ID)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("generating refresh token: %w", err)
	}

	return AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
	}, nil
}

func (uc *AuthUC) VerifyPassword(ctx context.Context, in VerifyPasswordInput) (string, error) {
	if err := uc.verifyCurrentPassword(ctx, in.UserID, in.Password); err != nil {
		return "", err
	}

	changePasswordToken, err := uc.tokens.GenerateChangePasswordToken(in.UserID)
	if err != nil {
		return "", fmt.Errorf("generating change password token: %w", err)
	}

	return changePasswordToken, nil
}

func (uc *AuthUC) verifyCurrentPassword(ctx context.Context, userID, password string) error {
	if userID == "" {
		return domain.ErrEmptyUserID
	}
	if password == "" {
		return domain.ErrInvalidCredentials
	}

	creds, err := uc.userRepo.FindCredentialsByUserID(ctx, userID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	var cred *domain.UserCredential
	for _, c := range creds {
		if c.IsPrimary {
			cred = &c
			break
		}
	}
	if cred == nil {
		return domain.ErrUserNotFound
	}

	if !cred.CheckPassword(password) {
		return domain.ErrIncorrectPassword
	}
	return nil
}

func (uc *AuthUC) ChangePassword(ctx context.Context, in ChangePasswordInput) error {
	if !domain.IsConfirmMatch(in.NewPassword, in.ConfirmedPassword) {
		return domain.ErrNotMatchPassword
	}

	if err := domain.ValidatePassword(in.NewPassword); err != nil {
		return err
	}

	hash, err := domain.HashPassword(in.NewPassword)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	if err := uc.verifyCurrentPassword(ctx, in.UserID, in.CurrentPassword); err != nil {
		return err
	}

	if in.NewPassword == in.CurrentPassword {
		return domain.ErrMatchCurrentPassword
	}

	err = uc.userRepo.UpdatePassword(ctx, in.UserID, hash)
	if err != nil {
		return fmt.Errorf("updating password: %w", err)
	}

	uc.cache.InvalidateProfile(ctx, in.UserID)
	return nil
}

func (uc *AuthUC) ForgotPassword(ctx context.Context, email string) (string, error) {
	cred, err := uc.userRepo.FindCredentialByIdentifier(ctx, email)
	if err != nil {
		return "", domain.ErrUserNotFound
	}

	resetToken, err := uc.tokens.GenerateResetToken(cred.UserID)
	if err != nil {
		return "", fmt.Errorf("generating reset token: %w", err)
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", uc.baseURL, resetToken)

	uc.email.Dispatch(worker.EmailJob{
		Template: "forgot_password",
		Subject:  "Reset Your Password",
		To:       email,
		Data: struct{ ResetLink string }{
			ResetLink: resetLink,
		},
	})

	return resetToken, nil
}

func (uc *AuthUC) ResetPassword(ctx context.Context, in ResetPasswordInput) error {
	claims, err := uc.tokens.VerifyResetToken(in.Token)
	if err != nil {
		return domain.ErrInvalidToken
	}

	isBlacklist, err := uc.cache.IsBlacklisted(ctx, claims.JTI())
	if err != nil {
		return fmt.Errorf("checking token blacklist: %w", err)
	}

	if isBlacklist {
		return domain.ErrInvalidToken
	}

	if err := domain.ValidatePassword(in.NewPassword); err != nil {
		return err
	}

	hash, err := domain.HashPassword(in.NewPassword)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	err = uc.userRepo.UpdatePassword(ctx, claims.UserID, hash)
	if err != nil {
		return fmt.Errorf("updating password: %w", err)
	}

	if err := uc.cache.Add(ctx, claims.JTI(), time.Until(claims.ExpiresAtTime())); err != nil {
		return fmt.Errorf("adding to blacklist: %w", err)
	}
	uc.cache.InvalidateProfile(ctx, claims.UserID)

	return nil
}

func (uc *AuthUC) Logout(ctx context.Context, in LogoutInput) error {
	// blacklist access token
	accClaims, err := uc.tokens.VerifyAccessToken(in.AccessToken)
	if err == nil {
		accTtl := time.Until(accClaims.ExpiresAtTime())
		if accTtl > 0 {
			err = uc.cache.Add(ctx, accClaims.JTI(), accTtl)
			if err != nil {
				return fmt.Errorf("blacklisting access token: %w", err)
			}
		}
	}

	// blacklist refresh token
	refClaims, err := uc.tokens.VerifyRefreshToken(in.RefreshToken)
	if err == nil {
		refTtl := time.Until(refClaims.ExpiresAtTime())
		if refTtl > 0 {
			err = uc.cache.Add(ctx, refClaims.JTI(), refTtl)
			if err != nil {
				return fmt.Errorf("blacklisting refresh token: %w", err)
			}
		}
	}

	return nil
}

func (uc *AuthUC) RefreshToken(ctx context.Context, refreshTokenStr string) (AuthOutput, error) {
	claims, err := uc.tokens.VerifyRefreshToken(refreshTokenStr)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("%w: %v", domain.ErrInvalidToken, err)
	}

	isBlacklisted, err := uc.cache.IsBlacklisted(ctx, refreshTokenStr)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("checking token blacklist: %w", err)
	}
	if isBlacklisted {
		return AuthOutput{}, domain.ErrInvalidToken
	}

	u, err := uc.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return AuthOutput{}, domain.ErrUserNotFound
		}
		return AuthOutput{}, fmt.Errorf("fetching user: %w", err)
	}

	switch u.Status {
	case domain.UserStatusBanned:
		return AuthOutput{}, domain.ErrUserBanned
	case domain.UserStatusUnverified, domain.UserStatusVerified:
		// OK
	default:
		return AuthOutput{}, fmt.Errorf("unknown user status: %s", u.Status)
	}

	newAccessToken, err := uc.tokens.GenerateAccessToken(u.ID, u.Role)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("generating access token: %w", err)
	}

	newRefreshToken, err := uc.tokens.GenerateRefreshToken(u.ID)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("generating refresh token: %w", err)
	}

	refTtl := time.Until(claims.ExpiresAtTime())
	if refTtl > 0 {
		_ = uc.cache.Add(ctx, claims.JTI(), refTtl)
	}

	return AuthOutput{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		UserID:       u.ID,
	}, nil
}

func (uc *AuthUC) GetRoleByIdentifier(ctx context.Context, identifier string) (domain.UserRole, error) {
	u, err := uc.getUserByIdentifier(ctx, identifier)
	if err != nil {
		return "", err
	}
	return u.Role, nil
}

func (uc *AuthUC) getUserByIdentifier(ctx context.Context, identifier string) (*domain.User, error) {
	u, err := uc.userRepo.FindByIdentifier(ctx, identifier)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return nil, fmt.Errorf("finding by identifier: %w", err)
		}
		u, err = uc.userRepo.FindByUsername(ctx, identifier)
		if err != nil {
			return nil, fmt.Errorf("finding by username: %w", err)
		}
	}
	return u, nil
}
