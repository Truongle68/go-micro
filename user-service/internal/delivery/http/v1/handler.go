package v1

import (
	"errors"
	"net/http"
	"user-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *V1) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyExists),
		errors.Is(err, domain.ErrPhoneAlreadyExists),
		errors.Is(err, domain.ErrUsernameExists):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		response.Error(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrUserBanned),
		errors.Is(err, domain.ErrUnauthorizedUser),
		errors.Is(err, domain.ErrUserInactive):
		response.Error(c, http.StatusForbidden, err.Error())
	case errors.Is(err, domain.ErrWeakPassword),
		errors.Is(err, domain.ErrEmailRequired),
		errors.Is(err, domain.ErrInvalidToken),
		errors.Is(err, domain.ErrInvalidOTP),
		errors.Is(err, domain.ErrOTPExpired),
		errors.Is(err, domain.ErrNotMatchPassword),
		errors.Is(err, domain.ErrEmptyUsername),
		errors.Is(err, domain.ErrEmptyEmail),
		errors.Is(err, domain.ErrEmptyPhone),
		errors.Is(err, domain.ErrEmptyAddressLine),
		errors.Is(err, domain.ErrEmptyUserID),
		errors.Is(err, domain.ErrEmptyCity),
		errors.Is(err, domain.ErrSameEmail),
		errors.Is(err, domain.ErrSamePhone),
		errors.Is(err, domain.ErrEmailNotSet),
		errors.Is(err, domain.ErrInvalidFullName),
		errors.Is(err, domain.ErrInvalidGender),
		errors.Is(err, domain.ErrInvalidDob),
		errors.Is(err, domain.ErrNoFieldsToUpdate):
		response.Error(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrUserNotFound),
		errors.Is(err, domain.ErrAddressNotFound):
		response.Error(c, http.StatusNotFound, err.Error())
	default:
		r.l.Error(err, "auth handler unexpected error")
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
