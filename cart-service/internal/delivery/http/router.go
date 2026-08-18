package http

import (
	v1 "cart-service/internal/delivery/http/v1"

	"github.com/TruongLe68/go-micro/pkg/ginmw"
	"github.com/gin-gonic/gin"
)

func NewRouter(engine *gin.Engine, deps *v1.Dependencies) {
	engine.Use(ginmw.Recovery(deps.Logger))
	engine.Use(ginmw.Logger(deps.Logger))
	engine.Use(corsMiddleware())

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
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
