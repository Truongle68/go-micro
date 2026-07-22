package middleware

import (
	"catalog-service/pkg/redis"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/GoProOrg/core-go-pkg/jwtmanager"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	userRoleCtx         = "userRole"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)

func Auth(v jwtmanager.JWTManager, bl redis.BlacklistCacher) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "empty auth header")
			c.Abort()
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "invalid auth header format")
			c.Abort()
			return
		}

		tokenStr := headerParts[1]

		// check if token is blacklisted in Redis
		inBlacklist, err := bl.IsBlacklisted(c, tokenStr)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		if inBlacklist {
			response.Error(c, http.StatusUnauthorized, "token has been logged out")
			c.Abort()
			return
		}

		claims, err := v.VerifyAccessToken(tokenStr)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set(userCtx, claims.UserID)
		c.Set(userRoleCtx, claims.Role)
		c.Next()
	}
}

func Role(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(userRoleCtx)
		if !exists {
			response.Error(c, http.StatusForbidden, "forbidden: missing role")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			response.Error(c, http.StatusForbidden, "forbidden: invalid role format")
			c.Abort()
			return
		}

		if slices.Contains(allowedRoles, roleStr) {
			c.Next()
			return
		}

		response.Error(c, http.StatusForbidden, "forbidden: insufficient permissions")
		c.Abort()
	}
}

func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(userCtx)
	if !exists {
		return "", false
	}
	userIDStr, ok := userID.(string)
	return userIDStr, ok
}
