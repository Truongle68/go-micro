package httpbind

import (
	"net/http"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindAndValidate[T any](c *gin.Context, v *validator.Validate, l logger.Interface, handlerName string) (T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		l.Warn("restapi - v1 - %s - ShouldBindJSON: %v", handlerName, err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return req, false
	}
	if err := v.Struct(req); err != nil {
		l.Warn("restapi - v1 - %s - validate: %v", handlerName, err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return req, false
	}
	return req, true
}
