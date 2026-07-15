package v1

import (
	"net/http"
	"strings"
	"user-service/pkg/jwt"
	"user-service/pkg/redis"

	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
	_ "github.com/redis/go-redis/v9"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	userRoleCtx         = "userRole"
	tokenCtx            = "token"
)

func authMiddleware(jwtService jwt.TokenService, bl redis.BlacklistCacher) gin.HandlerFunc {
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

		claims, err := jwtService.VerifyAccessToken(tokenStr)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set(userCtx, claims.UserID)
		c.Set(userRoleCtx, claims.Role)
		c.Set(tokenCtx, tokenStr)
		c.Next()
	}
}

func getRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(userRoleCtx)
	if !exists {
		return "", false
	}
	roleStr, ok := role.(string)
	return roleStr, ok
}

func (r *V1) getUserId(c *gin.Context) (string, bool) {
	id, exists := c.Get(userCtx)
	if !exists {
		return "", false
	}
	idStr, ok := id.(string)
	return idStr, ok
}

func (r *V1) getToken(c *gin.Context) (string, bool) {
	tok, exists := c.Get(tokenCtx)
	if !exists {
		return "", false
	}
	tokStr, ok := tok.(string)
	return tokStr, ok
}
