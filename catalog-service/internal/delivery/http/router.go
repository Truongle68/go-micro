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
	engine.Use(corsMiddleware())
	// define routers
	apiV1Group := engine.Group("/api/v1")
	v1.NewRoutes(apiV1Group, deps)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Header", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Method", "GET, OPTIONS, POST, PUT, DELETE")
		c.Next()
	}
}
