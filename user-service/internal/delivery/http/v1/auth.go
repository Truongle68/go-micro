package v1

import (
	"net/http"
	"strings"
	"user-service/internal/delivery/http/v1/req"
	"user-service/internal/delivery/http/v1/res"
	"user-service/internal/domain"
	"user-service/internal/usecase"
	"user-service/pkg/jwt"

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
		response.Error(c, http.StatusBadRequest, string(domain.CodeNotMatchPassword), domain.ErrNotMatchPassword.Error())
		return
	}

	out, err := r.a.CompleteRegister(c.Request.Context(), request.ToInput())
	if err != nil {
		r.l.Error(err, "restapi - v1 - completeRegister - usecase failed")
		r.handleError(c, err)
		return
	}

	setAuthCookies(c, out.AccessToken, out.RefreshToken)
	response.Success(c, http.StatusOK, "register success", res.ToAuthResponse(out))
}

func (r *V1) checkAvailableUsername(c *gin.Context) {
	username := c.Query("username")
	trimU := strings.TrimSpace(username)
	if len(trimU) < 3 {
		response.Error(c, http.StatusBadRequest, "INVALID_USERNAME_LENGTH", "username must have at least 3 characters")
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

	setAuthCookies(c, out.AccessToken, out.RefreshToken)
	response.Success(c, http.StatusOK, "login success", res.ToAuthResponse(out))
}

func (r *V1) logout(c *gin.Context) {
	cookieAccessToken, err := c.Cookie("access_token")
	if err != nil || cookieAccessToken == "" {
		response.Unauthorized(c, "missing access token cookie")
		return
	}
	cookieRefreshToken, err := c.Cookie("refresh_token")
	if err != nil || cookieRefreshToken == "" {
		response.Unauthorized(c, "missing refresh token cookie")
		return
	}

	err = r.a.Logout(c.Request.Context(), usecase.LogoutInput{
		AccessToken:  cookieAccessToken,
		RefreshToken: cookieRefreshToken,
	})
	if err != nil {
		r.handleError(c, err)
		return
	}

	clearAuthCookies(c)
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

func (r *V1) verifyPassword(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Unauthorized(c)
		return
	}
	request, ok := httpbind.BindAndValidate[req.VerifyPassword](c, r.v, r.l, "verifyPassword")
	if !ok {
		return
	}

	token, err := r.a.VerifyPassword(c.Request.Context(), request.ToInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "verify password success", res.ToVerifyPasswordResponse(token))
}

func (r *V1) changePassword(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Unauthorized(c)
		return
	}
	request, ok := httpbind.BindAndValidate[req.ChangePassword](c, r.v, r.l, "changePassword")
	if !ok {
		return
	}

	err := r.a.ChangePassword(c.Request.Context(), request.ToInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "change password success", nil)
}

func (r *V1) resetPassword(c *gin.Context) {
	request, ok := httpbind.BindAndValidate[req.ResetPassword](c, r.v, r.l, "resetPassword")
	if !ok {
		return
	}

	if !domain.IsConfirmMatch(request.NewPassword, request.ConfirmedPassword) {
		response.Error(c, http.StatusBadRequest, string(domain.CodeNotMatchPassword), domain.ErrNotMatchPassword.Error())
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
	cookieRefreshToken, err := c.Cookie("refresh_token")
	if err != nil || cookieRefreshToken == "" {
		response.Unauthorized(c, "missing refresh token cookie")
		return
	}

	authOut, err := r.a.RefreshToken(c.Request.Context(), cookieRefreshToken)
	if err != nil {
		r.l.Error(err, "restapi - v1 - refreshToken - usecase failed")
		r.handleError(c, err)
		return
	}

	setAuthCookies(c, authOut.AccessToken, authOut.RefreshToken)
	response.Success(c, http.StatusOK, "token refreshed successfully", res.ToAuthResponse(authOut))
}

func setAuthCookies(c *gin.Context, accessToken, refreshToken jwt.GeneratedTokenOutput) {
	c.SetCookie("access_token", accessToken.Token, int(accessToken.Exp.Seconds()), "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken.Token, int(refreshToken.Exp.Seconds()), "/", "", false, true)
}

func clearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
}
