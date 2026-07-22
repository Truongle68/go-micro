package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/gin-gonic/gin"
)

func buildPanicMessage(c *gin.Context, err any) string {
	var result strings.Builder

	result.WriteString(c.ClientIP())
	result.WriteString(" - ")
	result.WriteString(c.Request.Method)
	result.WriteString(" ")
	result.WriteString(c.Request.URL.String())
	result.WriteString(" PANIC DETECTED: ")
	fmt.Fprintf(&result, "%v\n%s\n", err, debug.Stack())

	return result.String()
}

func Recovery(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l.Error(buildPanicMessage(c, err))

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
