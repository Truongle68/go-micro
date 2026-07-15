package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repo"
	"user-service/pkg/jwt"
	"user-service/pkg/mailer"
	"user-service/pkg/postgres"
	"user-service/pkg/redis"

	"github.com/TruongLe68/go-micro/pkg/logger"
)

type AuthUC struct {
	userRepo   repo.UserRepository
	tokens     jwt.TokenService
	cache      redis.IdentityCacher
	transactor postgres.Transactor
	mailer     mailer.Mailer
	logger     logger.Interface
}

func NewAuthUC(repo repo.UserRepository, tokens jwt.TokenService, cache redis.IdentityCacher, transactor postgres.Transactor, mailer mailer.Mailer, logger logger.Interface) *AuthUC {
	return &AuthUC{
		userRepo:   repo,
		tokens:     tokens,
		cache:      cache,
		transactor: transactor,
		mailer:     mailer,
		logger:     logger,
	}
}

var _ Auth = (*AuthUC)(nil)

func (uc *AuthUC) RequestOTP(ctx context.Context, in RequestOTPInput) error {
	if in.Phone == "" {
		return domain.ErrEmptyPhone
	}

	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return fmt.Errorf("generating random OTP: %w", err)
	}
	code := fmt.Sprintf("%06d", n.Int64())

	fmt.Printf(" [OTP DEBUG] Generated OTP for %s (%s): %s\n", in.Phone, in.Purpose, code)

	err = uc.cache.SetOTP(ctx, in.Phone, in.Purpose, code, 1*time.Minute)
	if err != nil {
		return fmt.Errorf("caching OTP: %w", err)
	}

	return nil
}

func (uc *AuthUC) VerifyOTP(ctx context.Context, in VerifyOTPInput) (token string, exists bool, username string, err error) {
	if in.Phone == "" {
		return "", false, "", domain.ErrEmptyPhone
	}
	if in.Code == "" {
		return "", false, "", domain.ErrInvalidOTP
	}

	cachedCode, err := uc.cache.GetOTP(ctx, in.Phone, in.Purpose)
	if err != nil {
		return "", false, "", domain.ErrOTPExpired
	}

	if cachedCode != in.Code {
		return "", false, "", domain.ErrInvalidOTP
	}

	_ = uc.cache.DeleteOTP(ctx, in.Phone, in.Purpose)

	token, err = uc.tokens.GenerateVerificationToken(in.Phone, in.Purpose)
	if err != nil {
		return "", false, "", fmt.Errorf("generating verification token: %w", err)
	}

	// Check if user exists with this phone number
	cred, err := uc.userRepo.FindCredentialByIdentifier(ctx, in.Phone)
	if err == nil && cred != nil {
		exists = true
		user, err := uc.userRepo.FindByID(ctx, cred.UserID)
		if err == nil && user != nil {
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

	cred, err := domain.NewUserCredential(user.ID, domain.CredentialTypePhone, claims.Phone, in.Password, true, true)
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

	accessToken, err := uc.tokens.GenerateAccessToken(user.ID, string(user.Role))
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
				for _, c := range creds {
					if c.IsPrimary {
						cred = c
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

	switch user.Status {
	case domain.UserStatusBanned:
		return AuthOutput{}, domain.ErrUserBanned
	case domain.UserStatusUnverified, domain.UserStatusVerified:
		// OK, only making order requires verified user
	default:
		return AuthOutput{}, fmt.Errorf("unknown user status: %s", user.Status)
	}

	accessToken, err := uc.tokens.GenerateAccessToken(user.ID, string(user.Role))
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

func (uc *AuthUC) ForgotPassword(ctx context.Context, email string) (string, error) {
	cred, err := uc.userRepo.FindCredentialByIdentifier(ctx, email)
	if err != nil {
		return "", domain.ErrUserNotFound
	}

	resetToken, err := uc.tokens.GenerateResetToken(cred.UserID)
	if err != nil {
		return "", fmt.Errorf("generating reset token: %w", err)
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// protect against panics in the goroutine
		defer func() {
			if r := recover(); r != nil {
				uc.logger.Error("panic occurred in ForgotPassword email sending goroutine: %v", r)
			}
		}()

		// define email template later
		subject := "Reset Your Password"
		body := fmt.Sprintf(
			"<p>You requested a password reset. Please use the following token to reset your password:</p>"+
				"<p><strong>%s</strong></p>"+
				"<p>Or click this link: <a href=\"http://localhost:3000/reset-password?token=%s\">Reset Password</a></p>",
			resetToken, resetToken,
		)

		err := uc.mailer.Send(bgCtx, email, subject, body)
		if err != nil {
			uc.logger.Error("failed to send reset password email to %s: %v", email, err)
		} else {
			uc.logger.Info("sent reset password email to %s", email)
		}
	}()

	return resetToken, nil
}

func (uc *AuthUC) ResetPassword(ctx context.Context, in ResetPasswordInput) error {
	claims, err := uc.tokens.VerifyResetToken(in.Token)
	if err != nil {
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

	uc.cache.InvalidateProfile(ctx, claims.UserID)

	return nil
}

func (uc *AuthUC) Logout(ctx context.Context, in LogoutInput) error {
	// blacklist access token
	accClaims, err := uc.tokens.VerifyAccessToken(in.AccessToken)
	if err == nil {
		accTtl := time.Until(accClaims.ExpiresAt.Time)
		if accTtl > 0 {
			err = uc.cache.Blacklist(ctx, in.AccessToken, accTtl)
			if err != nil {
				return fmt.Errorf("blacklisting access token: %w", err)
			}
		}
	}

	// blacklist refresh token
	refClaims, err := uc.tokens.VerifyRefreshToken(in.RefreshToken)
	if err == nil {
		refTtl := time.Until(refClaims.ExpiresAt.Time)
		if refTtl > 0 {
			err = uc.cache.Blacklist(ctx, in.RefreshToken, refTtl)
			if err != nil {
				return fmt.Errorf("blacklisting refresh token: %w", err)
			}
		}
	}

	return nil
}
