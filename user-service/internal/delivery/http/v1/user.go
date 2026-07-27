package v1

import (
	"net/http"
	"user-service/internal/delivery/http/v1/req"
	"user-service/internal/delivery/http/v1/res"
	"user-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/httpbind"
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

	request, ok := httpbind.BindAndValidate[req.UpdateProfile](c, r.v, r.l, "updateProfile")
	if !ok {
		return
	}

	profile, err := r.u.UpdateProfile(c.Request.Context(), request.ToInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "update profile success", res.ToProfileResponse(profile))
}

func (r *V1) requestEmailLink(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	request, ok := httpbind.BindAndValidate[req.RequestEmailLink](c, r.v, r.l, "requestEmailLink")
	if !ok {
		return
	}

	err := r.u.RequestEmailLink(c.Request.Context(), request.ToInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "email link sent successfully", nil)
}

func (r *V1) confirmEmailLink(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.Error(c, http.StatusBadRequest, "token is required")
		return
	}

	out, err := r.u.ConfirmEmailLink(c.Request.Context(), token)
	if err != nil {
		r.handleError(c, err)
		return
	}

	if out.Purpose == domain.EmailLinkPurposeVerifyCurrent {
		response.Success(c, http.StatusOK, "current email verified successfully", gin.H{
			"change_email_token": out.ChangeEmailToken,
		})
		return
	}

	response.Success(c, http.StatusOK, "email verified successfully", nil)
}

func (r *V1) changeEmail(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	request, ok := httpbind.BindAndValidate[req.ChangeEmail](c, r.v, r.l, "changeEmail")
	if !ok {
		return
	}

	err := r.u.SendChangeEmailOTP(c.Request.Context(), request.ToInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "otp code sent successfully", nil)
}

func (r *V1) changeEmailVerify(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	request, ok := httpbind.BindAndValidate[req.ChangeEmailConfirm](c, r.v, r.l, "changeEmailVerify")
	if !ok {
		return
	}
	err := r.u.VerifyChangeEmailOTP(c.Request.Context(), request.ToInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "email changed successfully", nil)
}

func (r *V1) requestPhoneVerification(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := r.u.SendPhoneVerificationOTP(c.Request.Context(), userID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "otp code sent successfully", nil)
}

func (r *V1) verifyPhone(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	request, ok := httpbind.BindAndValidate[req.VerifyPhone](c, r.v, r.l, "verifyPhone")
	if !ok {
		return
	}

	token, err := r.u.VerifyPhoneVerificationOTP(c.Request.Context(), request.ToInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "otp verified successfully", res.VerifyPhoneResponse{
		ChangePhoneToken: token,
	})
}

func (r *V1) changePhone(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	request, ok := httpbind.BindAndValidate[req.ChangePhone](c, r.v, r.l, "changePhone")
	if !ok {
		return
	}

	err := r.u.SendChangePhoneOTP(c.Request.Context(), request.ToInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "otp code sent successfully", nil)
}

func (r *V1) changePhoneVerify(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	request, ok := httpbind.BindAndValidate[req.ChangePhoneVerify](c, r.v, r.l, "changePhoneVerify")
	if !ok {
		return
	}

	err := r.u.VerifyChangePhoneOTP(c.Request.Context(), request.ToInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "phone changed successfully", nil)
}

func (r *V1) getAddressList(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	addresses, err := r.u.GetAddressList(c.Request.Context(), userID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "get address list success", res.ToAddressListResponse(addresses))
}

func (r *V1) addAddress(c *gin.Context) {
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	request, ok := httpbind.BindAndValidate[req.CreateAddress](c, r.v, r.l, "addAddress")
	if !ok {
		return
	}

	address, err := r.u.CreateAddress(c.Request.Context(), request.ToInput(userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "create user address success", res.ToAddressResponse(address))
}

func (r *V1) setDefaultAddress(c *gin.Context) {
	id := c.Param("id")
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := r.u.SetDefaultAddress(c.Request.Context(), userID, id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "set default address success", nil)
}

func (r *V1) updateAddress(c *gin.Context) {
	id := c.Param("id")
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	request, ok := httpbind.BindAndValidate[req.UpdateAddress](c, r.v, r.l, "updateAddress")
	if !ok {
		return
	}

	address, err := r.u.UpdateAddress(c.Request.Context(), request.ToInput(id, userID))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "update user address success", res.ToAddressResponse(address))
}

func (r *V1) deleteAddress(c *gin.Context) {
	id := c.Param("id")
	userID, ok := r.getUserId(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := r.u.DeleteAddress(c.Request.Context(), userID, id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "delete user address success", nil)
}
