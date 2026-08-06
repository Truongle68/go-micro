package response

import (
	"net/http"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/gin-gonic/gin"
)

// Common Error Codes
const (
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeForbidden           = "FORBIDDEN"
	CodeValidationError     = "VALIDATION_ERROR"
	CodeInvalidRequestBody  = "INVALID_REQUEST_BODY"
	CodeInvalidQueryParams  = "INVALID_QUERY_PARAMS"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
	CodeInvalidAuthHeader   = "INVALID_AUTH_HEADER"
	CodeTokenLoggedOut      = "TOKEN_LOGGED_OUT"
	CodeInvalidToken        = "INVALID_TOKEN"
	CodeMissingToken        = "MISSING_TOKEN"
)

type envelope struct {
	Success bool        `json:"success"`
	Code    string      `json:"code,omitempty"`
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

func Error(c *gin.Context, httpStatus int, errCode string, msg string) {
	c.JSON(httpStatus, envelope{
		Success: false,
		Code:    errCode,
		Message: msg,
	})
}

// Common error response helpers

func Unauthorized(c *gin.Context, msg ...string) {
	m := "unauthorized"
	if len(msg) > 0 && msg[0] != "" {
		m = msg[0]
	}
	Error(c, http.StatusUnauthorized, CodeUnauthorized, m)
}

func Forbidden(c *gin.Context, msg ...string) {
	m := "forbidden"
	if len(msg) > 0 && msg[0] != "" {
		m = msg[0]
	}
	Error(c, http.StatusForbidden, CodeForbidden, m)
}

func ValidationError(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, CodeValidationError, msg)
}

func InvalidRequestBody(c *gin.Context, msg ...string) {
	m := "invalid request body"
	if len(msg) > 0 && msg[0] != "" {
		m = msg[0]
	}
	Error(c, http.StatusBadRequest, CodeInvalidRequestBody, m)
}

func InvalidQueryParams(c *gin.Context, msg ...string) {
	m := "invalid query parameters"
	if len(msg) > 0 && msg[0] != "" {
		m = msg[0]
	}
	Error(c, http.StatusBadRequest, CodeInvalidQueryParams, m)
}

func InternalServerError(c *gin.Context, msg ...string) {
	m := "internal server error"
	if len(msg) > 0 && msg[0] != "" {
		m = msg[0]
	}
	Error(c, http.StatusInternalServerError, CodeInternalServerError, m)
}
