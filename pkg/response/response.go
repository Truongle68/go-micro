package response

import (
	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(code, envelope{Success: true, Message: msg, Data: data})
}

func SuccessPaginated[T any](c *gin.Context, code int, msg string, result pagination.Result[T]) {
	c.JSON(code, envelope{
		Success: true,
		Message: msg,
		Data:    result.Items,
		Meta:    result.Meta,
	})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, envelope{Success: false, Message: msg})
}
