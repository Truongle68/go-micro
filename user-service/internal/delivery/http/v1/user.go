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

	response.Success(c, http.StatusOK, "get profile success", res.ToProfileResponse(profile))
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

	profile, err := r.u.UpdateProfile(c.Request.Context(), request.ToUpdatedProfileInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "update profile success", res.ToProfileResponse(profile))
}

func (r *V1) requestChangeEmail(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := r.u.RequestChangeEmailOTP(c.Request.Context(), userID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "verification code sent to your email", nil)
}

func (r *V1) verifyChangeEmail(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var request req.VerifyChangeEmail
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(request); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := r.u.VerifyChangeEmailOTP(c.Request.Context(), request.ToVerifyChangeEmailInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "otp verified successfully", res.VerifyChangePhoneResponse{
		Token: token,
	})
}

func (r *V1) completeChangeEmail(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var request req.CompleteChangeEmail
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(request); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := r.u.CompleteChangeEmail(c.Request.Context(), request.ToCompleteChangeEmailInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "email updated successfully", nil)
}

func (r *V1) requestChangePhone(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := r.u.RequestChangePhoneLink(c.Request.Context(), userID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "phone reset link sent to your email", nil)
}

func (r *V1) completeChangePhone(c *gin.Context) {
	var request req.CompleteChangePhone
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(request); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := r.u.CompleteChangePhone(c.Request.Context(), request.ToCompleteChangePhoneInput())
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "phone updated successfully", nil)
}
