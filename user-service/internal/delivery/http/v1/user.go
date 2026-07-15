package v1

import (
	"net/http"
	"user-service/internal/delivery/http/v1/req"
	"user-service/internal/delivery/http/v1/res"

	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *V1) getProfile(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := r.u.GetProfile(c.Request.Context(), userID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "get profile success", res.ProfileResponse{
		ID:        profile.ID,
		Username:  profile.Username,
		Email:     profile.Email,
		Phone:     profile.Phone,
		FullName:  profile.FullName,
		Role:      profile.Role,
		Status:    string(profile.Status),
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	})
}

func (r *V1) updateProfile(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var request req.UpdateProfile
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Warn("restapi - v1 - updateProfile - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(request); err != nil {
		r.l.Warn("restapi - v1 - updateProfile - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	profile, err := r.u.UpdateProfile(c.Request.Context(), userID, request.FullName, request.Phone, request.Email)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "update profile success", res.ToProfileResponse(profile))
}
