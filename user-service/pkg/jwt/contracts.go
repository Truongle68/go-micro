package jwt

type TokenService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	GenerateResetToken(userID string) (string, error)
	GenerateVerificationToken(phone, purpose string) (string, error)
	VerifyAccessToken(tokenStr string) (*Claims, error)
	VerifyRefreshToken(tokenStr string) (*Claims, error)
	VerifyResetToken(tokenStr string) (*Claims, error)
	VerifyVerificationToken(tokenStr string) (*Claims, error)
}
