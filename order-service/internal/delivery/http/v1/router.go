package v1

import (
	"order-service/internal/usecase"

	"github.com/GoProOrg/core-go-pkg/jwtmanager"
	redismanager "github.com/GoProOrg/core-go-pkg/redismanager/identity"
	"github.com/TruongLe68/go-micro/pkg/ginmw"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type V1 struct {
	o usecase.Order
	l logger.Interface
	v *validator.Validate
}

type Dependencies struct {
	Order    usecase.Order
	Verifier jwtmanager.JWTManager
	Cache    redismanager.BlacklistCacher
	Logger   logger.Interface
}

func NewDependencies(order usecase.Order, verifier jwtmanager.JWTManager, cache redismanager.BlacklistCacher, l logger.Interface) *Dependencies {
	return &Dependencies{
		Order:    order,
		Verifier: verifier,
		Cache:    cache,
		Logger:   l,
	}
}

func NewRoutes(apiV1Group *gin.RouterGroup, deps *Dependencies) {
	r := &V1{
		o: deps.Order,
		l: deps.Logger,
		v: validator.New(),
	}

	ordersGroup := apiV1Group.Group("/orders")
	ordersGroup.Use(ginmw.Auth(deps.Verifier, deps.Cache))

	ordersGroup.POST("/checkout", r.checkout)
	ordersGroup.GET("", r.listOrders)
	ordersGroup.GET("/:id", r.getOrder)
	ordersGroup.GET("/:id/tracking", r.getTracking)
	ordersGroup.POST("/:id/ship", r.shipOrder)
	ordersGroup.POST("/:id/deliver", r.deliverOrder)
	ordersGroup.POST("/:id/cancel", r.cancelOrder)
}
