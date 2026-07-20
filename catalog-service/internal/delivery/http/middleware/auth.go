package middleware

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

type Claims struct {
	UserID string `json:"user_id,omitempty"`
	Role   string `json:"role,omitempty"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func AuthMiddleware(accessSecret string) gin.HandlerFunc {
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

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			return []byte(accessSecret), nil
		})

		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		if claims.Type != "access" {
			response.Error(c, http.StatusUnauthorized, "invalid token type")
			c.Abort()
			return
		}

		c.Set(userCtx, claims.UserID)
		c.Set(userRoleCtx, claims.Role)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
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
