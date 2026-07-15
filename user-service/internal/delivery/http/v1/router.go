package v1

import (
	"user-service/internal/usecase"
	"user-service/pkg/jwt"
	"user-service/pkg/redis"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type V1 struct {
	a usecase.Auth
	u usecase.User
	l logger.Interface
	v *validator.Validate
}

type Dependencies struct {
	Auth   usecase.Auth
	User   usecase.User
	JWT    jwt.TokenService
	cache  redis.BlacklistCacher
	Logger logger.Interface
}

func NewDependencies(auth usecase.Auth, user usecase.User, jwt jwt.TokenService, cache redis.BlacklistCacher, l logger.Interface) *Dependencies {
	return &Dependencies{
		Auth:   auth,
		User:   user,
		JWT:    jwt,
		cache:  cache,
		Logger: l,
	}
}

func NewRoutes(apiV1Group *gin.RouterGroup, deps *Dependencies) {
	r := &V1{
		a: deps.Auth,
		u: deps.User,
		l: deps.Logger,
		v: validator.New(),
	}

	// public routes
	auth := apiV1Group.Group("/auth")
	auth.POST("/register/request-otp", r.requestOTP)
	auth.POST("/register/verify-otp", r.verifyOTP)
	auth.POST("/register/complete", r.completeRegister)
	auth.GET("/check-username", r.checkAvailableUsername)
	auth.POST("/login", r.login)
	auth.POST("/forgot-password", r.forgotPassword)
	auth.POST("/reset-password", r.resetPassword)

	// private/protected routes
	protected := apiV1Group.Group("")
	protected.Use(authMiddleware(deps.JWT, deps.cache))

	protected.POST("/auth/logout", r.logout)

	users := protected.Group("/users")
	users.GET("/profile", r.getProfile)
	users.PUT("/profile", r.updateProfile)
}
