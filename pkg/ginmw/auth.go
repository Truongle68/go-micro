package ginmw

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/GoProOrg/core-go-pkg/jwtmanager"
	redismanager "github.com/GoProOrg/core-go-pkg/redismanager/identity"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	CtxUserID = "user_id"
	CtxRole   = "user_role"

	DefaultAccessCookie = "access_token"
)

const (
	AdminRole = "admin"
	UserRole  = "user"
)

type middlewareConfig struct {
	cookieName string
	useBearer  bool
}

type MiddlewareOption func(*middlewareConfig)

func CookieName(name string) MiddlewareOption {
	return func(c *middlewareConfig) {
		c.cookieName = name
	}
}

func Bearer() MiddlewareOption {
	return func(c *middlewareConfig) {
		c.useBearer = true
	}
}

func Auth(v jwtmanager.JWTManager, bl redismanager.BlacklistCacher, opts ...MiddlewareOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := &middlewareConfig{
			cookieName: DefaultAccessCookie,
		}
		for _, opt := range opts {
			opt(cfg)
		}

		token, err := extractToken(c, cfg)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			c.Abort()
			return
		}

		if token == "" {
			response.Unauthorized(c, "missing access token")
			c.Abort()
			return
		}

		claims, err := v.VerifyAccessToken(token)
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

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

func extractToken(c *gin.Context, cfg *middlewareConfig) (string, error) {
	if cfg.useBearer {
		h := c.GetHeader("Authorization")
		if h == "" {
			return "", errors.New("missing authorization header")
		}
		const prefix = "Bearer "
		if !strings.HasPrefix(h, prefix) {
			return "", errors.New("invalid authorization header")
		}
		return strings.TrimSpace(h[len(prefix):]), nil
	}
	return c.Cookie(cfg.cookieName)
}

func Role(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(CtxRole)
		if !exists {
			response.Forbidden(c, "forbidden: missing role")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			response.Forbidden(c, "forbidden: invalid role format")
			c.Abort()
			return
		}

		if slices.Contains(allowedRoles, roleStr) {
			c.Next()
			return
		}

		response.Forbidden(c, "forbidden: insufficient permissions")
		c.Abort()
	}
}

func GetRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(CtxRole)
	if !exists {
		return "", false
	}
	roleStr, ok := role.(string)
	return roleStr, ok
}

func GetUserID(c *gin.Context) (string, bool) {
	id, exists := c.Get(CtxUserID)
	if !exists {
		return "", false
	}
	idStr, ok := id.(string)
	return idStr, ok
}
