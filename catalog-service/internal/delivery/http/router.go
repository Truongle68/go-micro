package http

import (
	v1 "catalog-service/internal/delivery/http/v1"

	"github.com/gin-gonic/gin"
)

func NewRouter(engine *gin.Engine, deps v1.Dependencies) {
	// use middleware
	// define routers
	apiV1Group := engine.Group("/v1")
	v1.NewRoutes(apiV1Group, deps)
}
