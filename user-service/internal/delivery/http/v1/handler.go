package v1

import (
	"errors"
	"net/http"
	"user-service/internal/delivery/http/middleware"
	"user-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *V1) handleError(c *gin.Context, err error) {
	appErr := domain.ToAppError(err)
	if appErr != nil {
		var status int
		switch {
		case errors.Is(err, domain.ErrEmailAlreadyExists),
			errors.Is(err, domain.ErrPhoneAlreadyExists),
			errors.Is(err, domain.ErrUsernameExists):
			status = http.StatusConflict
		case errors.Is(err, domain.ErrInvalidCredentials):
			status = http.StatusUnauthorized
		case errors.Is(err, domain.ErrUserBanned),
			errors.Is(err, domain.ErrUnauthorizedUser),
			errors.Is(err, domain.ErrUserInactive):
			status = http.StatusForbidden
		case errors.Is(err, domain.ErrUserNotFound),
			errors.Is(err, domain.ErrAddressNotFound):
			status = http.StatusNotFound
		default:
			status = http.StatusBadRequest
		}
		response.Error(c, status, string(appErr.Code), appErr.Message)
		return
	}

	r.l.Error("auth handler unexpected error: %v", err)
	response.InternalServerError(c)
}

func (r *V1) getUserId(c *gin.Context) (string, bool) {
	return middleware.GetUserID(c)
}

func (r *V1) getToken(c *gin.Context) (string, bool) {
	return middleware.GetToken(c)
}

func (r *V1) getRole(c *gin.Context) (string, bool) {
	return middleware.GetRole(c)
}
