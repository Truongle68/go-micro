package jwt

type TokenService interface {
	VerifyAccessToken(tokenStr string) (*Claims, error)
	VerifyRefreshToken(tokenStr string) (*Claims, error)
}
