package v1

import (
	"github.com/GoProOrg/core-go-pkg/jwtmanager"
	redismanager "github.com/GoProOrg/core-go-pkg/redismanager/identity"
	"github.com/TruongLe68/go-micro/pkg/ginmw"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type V1 struct {
	c CartUC
	l logger.Interface
	v *validator.Validate
}

type Dependencies struct {
	Cart     CartUC
	Verifier jwtmanager.JWTManager
	Cache    redismanager.BlacklistCacher
	Logger   logger.Interface
}

func NewDependencies(cart CartUC, verifier jwtmanager.JWTManager, cache redismanager.BlacklistCacher, l logger.Interface) *Dependencies {
	return &Dependencies{
		Cart:     cart,
		Verifier: verifier,
		Cache:    cache,
		Logger:   l,
	}
}

func NewRoutes(apiV1Group *gin.RouterGroup, deps *Dependencies) {
	r := &V1{
		c: deps.Cart,
		l: deps.Logger,
		v: validator.New(),
	}

	cartGroup := apiV1Group.Group("/cart")
	cartGroup.Use(ginmw.Auth(deps.Verifier, deps.Cache))

	cartGroup.GET("", r.getCart)
	cartGroup.POST("/items", r.addItem)
	cartGroup.PUT("/items/:sku", r.updateItemQuantity)
	cartGroup.DELETE("/items/:sku", r.removeItem)
	cartGroup.DELETE("/items", r.removeItems)
	cartGroup.DELETE("", r.clearCart)
}
