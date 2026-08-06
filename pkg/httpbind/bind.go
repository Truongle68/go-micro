package httpbind

import (
	"fmt"
	"strings"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindAndValidate[T any](c *gin.Context, v *validator.Validate, l logger.Interface, handlerName string) (T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		l.Warn("restapi - v1 - %s - ShouldBindJSON: %v", handlerName, err)
		response.InvalidRequestBody(c)
		return req, false
	}
	if err := v.Struct(req); err != nil {
		l.Warn("restapi - v1 - %s - validate: %v", handlerName, err)
		response.ValidationError(c, formatValidationError(err))
		return req, false
	}
	return req, true
}

func formatValidationError(err error) string {
	var errMsgs []string

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			errMsgs = append(errMsgs,
				fmt.Sprintf("%s is %s", fe.Field(), fe.Tag()))
		}
		return strings.Join(errMsgs, ", ")
	}

	return err.Error()
}
