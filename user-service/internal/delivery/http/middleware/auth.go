package middleware

import (
	"net/http"
	"user-service/pkg/jwt"
	"user-service/pkg/redis"

	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
	_ "github.com/redis/go-redis/v9"
)

const (
	authorizationHeader    = "Authorization"
	userCtx                = "userId"
	userRoleCtx            = "userRole"
	tokenCtx               = "token"
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

func Auth(jwtService jwt.TokenService, bl redis.BlacklistCacher) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieToken, err := c.Cookie(AccessTokenCookieName)
		if err != nil || cookieToken == "" {
			response.Unauthorized(c, "missing access token cookie")
			c.Abort()
			return
		}

		claims, err := jwtService.VerifyAccessToken(cookieToken)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, response.CodeInvalidToken, err.Error())
			c.Abort()
			return
		}
		// check if token is blacklisted in Redis
		inBlacklist, err := bl.IsBlacklisted(c, claims.JTI())
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		if inBlacklist {
			response.Error(c, http.StatusUnauthorized, response.CodeTokenLoggedOut, "token has been logged out")
			c.Abort()
			return
		}

		c.Set(userCtx, claims.UserID)
		c.Set(userRoleCtx, claims.Role)
		c.Set(tokenCtx, cookieToken)
		c.Next()
	}
}

func GetRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(userRoleCtx)
	if !exists {
		return "", false
	}
	roleStr, ok := role.(string)
	return roleStr, ok
}

func GetUserID(c *gin.Context) (string, bool) {
	id, exists := c.Get(userCtx)
	if !exists {
		return "", false
	}
	idStr, ok := id.(string)
	return idStr, ok
}

func GetToken(c *gin.Context) (string, bool) {
	tok, exists := c.Get(tokenCtx)
	if !exists {
		return "", false
	}
	tokStr, ok := tok.(string)
	return tokStr, ok
}
