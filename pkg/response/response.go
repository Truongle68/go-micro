package response

import "github.com/gin-gonic/gin"

type envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(code, envelope{Success: true, Message: msg, Data: data})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, envelope{Success: false, Message: msg})
}
