package v1

import (
	"errors"
	"net/http"
	"user-service/internal/delivery/http/v1/req"
	"user-service/internal/delivery/http/v1/res"
	"user-service/internal/domain"
	"user-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *V1) register(c *gin.Context) {
	var req req.Register
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Warn("restapi - v1 - register - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.v.Struct(req); err != nil {
		r.l.Warn("restapi - v1 - register - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	out, err := r.a.Register(c.Request.Context(), usecase.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		r.l.Error(err, "restapi - v1 - register - usecase failed")
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "register success", res.AuthResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		UserID:       out.UserID,
	})
}

func (r *V1) login(c *gin.Context) {
	var req req.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Warn("restapi - v1 - login - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(req); err != nil {
		r.l.Warn("restapi - v1 - login - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	out, err := r.a.Login(c.Request.Context(), usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "login success", res.AuthResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		UserID:       out.UserID,
	})
}

func (r *V1) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		response.Error(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrUserBanned), errors.Is(err, domain.ErrUserInactive):
		response.Error(c, http.StatusForbidden, err.Error())
	case errors.Is(err, domain.ErrWeakPassword), errors.Is(err, domain.ErrEmailRequired):
		response.Error(c, http.StatusBadRequest, err.Error())
	default:
		r.l.Error(err, "auth handler unexpected error")
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
