package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/gin-gonic/gin"
)

func buildRequestMessage(c *gin.Context, latency time.Duration) string {
	var result strings.Builder

	result.WriteString(c.ClientIP())
	result.WriteString(" - ")
	result.WriteString(c.Request.Method)
	result.WriteString(" ")
	result.WriteString(c.Request.URL.String())
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(c.Writer.Status()))
	result.WriteString(" ")
	result.WriteString(strconv.Itoa(c.Writer.Size()))
	result.WriteString(" ")
	result.WriteString(latency.String())

	return result.String()
}

func Logger(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		l.Info("%s", buildRequestMessage(c, latency))
	}
}
