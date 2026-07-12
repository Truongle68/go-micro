package v1

import (
	"user-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type V1 struct {
	a usecase.Auth
	l logger.Interface
	v *validator.Validate
}

type Dependencies struct {
	Auth   usecase.Auth
	Logger logger.Interface
}

func NewDependencies(auth usecase.Auth, l logger.Interface) *Dependencies {
	return &Dependencies{
		Auth:   auth,
		Logger: l,
	}
}

func NewRoutes(apiV1Group *gin.RouterGroup, deps *Dependencies) {
	r := &V1{a: deps.Auth, l: deps.Logger, v: validator.New()}
	// public routes
	auth := apiV1Group.Group("/auth")
	auth.POST("/register", r.register)
	auth.POST("/login", r.login)
	// private routes
}
