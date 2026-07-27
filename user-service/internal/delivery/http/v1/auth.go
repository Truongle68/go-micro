package v1

import (
	"net/http"
	"strings"
	"user-service/internal/delivery/http/v1/req"
	"user-service/internal/delivery/http/v1/res"
	"user-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/httpbind"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *V1) requestOTP(c *gin.Context, purpose domain.VerifyPurpose) {
	request, ok := httpbind.BindAndValidate[req.RequestOTP](c, r.v, r.l, "requestOTP")
	if !ok {
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
	request, ok := httpbind.BindAndValidate[req.VerifyOTP](c, r.v, r.l, "verifyOTP")
	if !ok {
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
	request, ok := httpbind.BindAndValidate[req.CompleteRegister](c, r.v, r.l, "completeRegister")
	if !ok {
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
	request, ok := httpbind.BindAndValidate[req.Login](c, r.v, r.l, contextTag)
	if !ok {
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
	request, ok := httpbind.BindAndValidate[req.Logout](c, r.v, r.l, "logout")
	if !ok {
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
	request, ok := httpbind.BindAndValidate[req.ForgotPassword](c, r.v, r.l, "forgotPassword")
	if !ok {
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
	request, ok := httpbind.BindAndValidate[req.ResetPassword](c, r.v, r.l, "resetPassword")
	if !ok {
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

func (r *V1) refreshToken(c *gin.Context) {
	request, ok := httpbind.BindAndValidate[req.RefreshToken](c, r.v, r.l, "refreshToken")
	if !ok {
		return
	}

	authOut, err := r.a.RefreshToken(c.Request.Context(), request.RefreshToken)
	if err != nil {
		r.l.Error(err, "restapi - v1 - refreshToken - usecase failed")
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "token refreshed successfully", res.ToAuthResponse(authOut))
}
