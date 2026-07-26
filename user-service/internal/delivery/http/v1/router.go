package v1

import (
	"user-service/internal/delivery/http/middleware"
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
	Cache  redis.BlacklistCacher
	Logger logger.Interface
}

func NewDependencies(auth usecase.Auth, user usecase.User, jwt jwt.TokenService, cache redis.BlacklistCacher, l logger.Interface) *Dependencies {
	return &Dependencies{
		Auth:   auth,
		User:   user,
		JWT:    jwt,
		Cache:  cache,
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

	// public auth routes
	auth := apiV1Group.Group("/auth")
	auth.POST("/register", r.requestRegisterOTP)
	auth.POST("/register/verify", r.verifyRegisterOTP)
	auth.POST("/register/complete", r.completeRegister)
	auth.GET("/check-username", r.checkAvailableUsername)
	auth.POST("/login", r.login)
	auth.POST("/portal/login", r.portalLogin)
	auth.POST("/forgot-password", r.forgotPassword)
	auth.POST("/reset-password", r.resetPassword)
	auth.POST("/refresh", r.refreshToken)

	// public email link confirm
	apiV1Group.GET("/users/verify-email/confirm", r.confirmEmailLink)

	// protected routes
	protected := apiV1Group.Group("")
	protected.Use(middleware.Auth(deps.JWT, deps.Cache))

	protected.POST("/auth/logout", r.logout)

	users := protected.Group("/users")
	// profile
	users.GET("/profile", r.getProfile)
	users.PUT("/profile", r.updateProfile)

	// verification & credential change
	users.POST("/verify-email", r.requestEmailLink)
	users.POST("/change-email", r.changeEmail)
	users.POST("/change-email/confirm", r.changeEmailVerify)
	users.POST("/verify-phone", r.requestPhoneVerification)
	users.POST("/verify-phone/confirm", r.verifyPhone)
	users.POST("/change-phone", r.changePhone)
	users.POST("/change-phone/confirm", r.changePhoneVerify)

	// address
	address := users.Group("/addresses")
	address.GET("", r.getAddressList)
	address.POST("", r.addAddress)
	address.POST("/:id/set-default", r.setDefaultAddress)
	address.PUT("/:id/update", r.updateAddress)
	address.DELETE("/:id/delete", r.deleteAddress)
}
