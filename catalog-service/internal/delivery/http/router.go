package http

import (
	"catalog-service/internal/delivery/http/middleware"
	v1 "catalog-service/internal/delivery/http/v1"

	"github.com/gin-gonic/gin"
)

func NewRouter(engine *gin.Engine, deps v1.Dependencies) {
	// use middleware
	engine.Use(middleware.Recovery(deps.Logger))
	engine.Use(middleware.Logger(deps.Logger))
	// define routers
	apiV1Group := engine.Group("/v1")
	v1.NewRoutes(apiV1Group, deps)
}
