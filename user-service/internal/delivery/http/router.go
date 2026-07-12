package http

import (
	"user-service/internal/delivery/http/middleware"
	v1 "user-service/internal/delivery/http/v1"

	"github.com/gin-gonic/gin"
)

func NewRouter(engine *gin.Engine, deps *v1.Dependencies) {
	engine.Use(middleware.Recovery(deps.Logger))
	engine.Use(middleware.Logger(deps.Logger))

	apiV1Group := engine.Group("/api/v1")
	v1.NewRoutes(apiV1Group, deps)
}
