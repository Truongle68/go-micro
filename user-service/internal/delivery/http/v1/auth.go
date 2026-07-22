package v1

import (
	"net/http"
	"strings"
	"user-service/internal/delivery/http/v1/req"
	"user-service/internal/delivery/http/v1/res"
	"user-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *V1) requestOTP(c *gin.Context, purpose domain.VerifyPurpose) {
	var request req.RequestOTP
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Warn("restapi - v1 - requestOTP - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.v.Struct(request); err != nil {
		r.l.Warn("restapi - v1 - requestOTP - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := r.a.RequestOTP(c.Request.Context(), request.ToInput(purpose))
	if err != nil {
		r.l.Error(err, "restapi - v1 - requestOTP - usecase failed")
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "otp code generated", nil)
}

func (r *V1) requestRegisterOTP(c *gin.Context) {
	r.requestOTP(c, domain.VerifyPurposeRegister)
}

func (r *V1) verifyOTP(c *gin.Context, purpose domain.VerifyPurpose) {
	var request req.VerifyOTP
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Warn("restapi - v1 - verifyOTP - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.v.Struct(request); err != nil {
		r.l.Warn("restapi - v1 - verifyOTP - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, exists, username, err := r.a.VerifyOTP(c.Request.Context(), request.ToInput(purpose))
	if err != nil {
		r.l.Error(err, "restapi - v1 - verifyOTP - usecase failed")
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "otp code verified", res.ToVerifyOTPResponse(token, exists, username))
}

func (r *V1) verifyRegisterOTP(c *gin.Context) {
	r.verifyOTP(c, domain.VerifyPurposeRegister)
}

func (r *V1) completeRegister(c *gin.Context) {
	var request req.CompleteRegister
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Warn("restapi - v1 - completeRegister - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.v.Struct(request); err != nil {
		r.l.Warn("restapi - v1 - completeRegister - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if !domain.IsConfirmMatch(request.Password, request.ConfirmedPassword) {
		response.Error(c, http.StatusBadRequest, domain.ErrNotMatchPassword.Error())
		return
	}

	out, err := r.a.CompleteRegister(c.Request.Context(), request.ToInput())
	if err != nil {
		r.l.Error(err, "restapi - v1 - completeRegister - usecase failed")
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "register success", res.ToAuthResponse(out))
}

func (r *V1) checkAvailableUsername(c *gin.Context) {
	username := c.Query("username")
	trimU := strings.TrimSpace(username)
	if len(trimU) < 3 {
		response.Error(c, http.StatusBadRequest, "username must have at least 3 characters")
		return
	}

	available, err := r.a.CheckUsernameAvailable(c.Request.Context(), username)
	if err != nil {
		r.l.Error(err, "restapi - v1 - checkAvailableUsername - usecase failed")
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "check username success", res.ToCheckUsernameResponse(available))
}

func (r *V1) login(c *gin.Context) {
	r.executeLoginFlow(c, "customer", []domain.UserRole{domain.UserRoleUser})
}

func (r *V1) portalLogin(c *gin.Context) {
	r.executeLoginFlow(c, "portal", domain.PortalRole)
}

func (r *V1) executeLoginFlow(c *gin.Context, contextTag string, requiredRoles []domain.UserRole) {
	var request req.Login
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Warn("restapi - v1 - %s - ShouldBindJSON: %v", contextTag, err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(request); err != nil {
		r.l.Warn("restapi - v1 - %s - validate: %v", contextTag, err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	out, err := r.a.Login(c.Request.Context(), request.ToInput(requiredRoles))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "login success", res.ToAuthResponse(out))
}

func (r *V1) logout(c *gin.Context) {
	var request req.Logout
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Warn("restapi - v1 - logout - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(request); err != nil {
		r.l.Warn("restapi - v1 - logout - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, ok := r.getToken(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := r.a.Logout(c.Request.Context(), request.ToInput(accessToken))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "logout success", nil)
}

func (r *V1) forgotPassword(c *gin.Context) {
	var request req.ForgotPassword
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Warn("restapi - v1 - forgotPassword - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(request); err != nil {
		r.l.Warn("restapi - v1 - forgotPassword - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := r.a.ForgotPassword(c.Request.Context(), request.Email)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "forgot password success", res.ToForgotPasswordResponse(token))
}

func (r *V1) resetPassword(c *gin.Context) {
	var request req.ResetPassword
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Warn("restapi - v1 - resetPassword - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(request); err != nil {
		r.l.Warn("restapi - v1 - resetPassword - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if !domain.IsConfirmMatch(request.NewPassword, request.ConfirmedPassword) {
		response.Error(c, http.StatusBadRequest, domain.ErrNotMatchPassword.Error())
		return
	}

	err := r.a.ResetPassword(c.Request.Context(), request.ToInput())
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "reset password success", nil)
}
