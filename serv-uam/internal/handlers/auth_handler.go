package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/dtos"
	"github.com/reshap0318/serv-uam/internal/helpers"
)

// AuthVerify checks an email/password pair for the api-gateway's
// login flow — outside GatewayAuth since there's no end-user identity yet.
func (h *Handlers) AuthVerify(c *gin.Context) {
	var req dtos.AuthVerifyRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.AuthVerify(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to verify credentials") {
		return
	}

	helpers.OK(c, "Credentials verified successfully", result)
}

// UserGetAccess resolves a user's current roles/permissions by ID —
// called by the api-gateway to refill its Redis cache on a miss.
func (h *Handlers) UserGetAccess(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	result, err := h.svcs.UserGetAccess(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to resolve user access") {
		return
	}

	helpers.OK(c, "User access retrieved successfully", result)
}

// RequestPasswordReset generates a reset token for an email. The response never
// carries the token itself — it only ever reaches the account owner via email
// (see AuthRequestPasswordReset), never the API caller.
func (h *Handlers) RequestPasswordReset(c *gin.Context) {
	var req dtos.RequestPasswordResetRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	err := h.svcs.AuthRequestPasswordReset(c.Request.Context(), req.Email)
	if helpers.HandleError(c, err, "Failed to request password reset") {
		return
	}

	helpers.OK(c, "If the email is registered, a reset link has been sent", nil)
}

// ResetPassword validates a reset token and updates the user's password.
func (h *Handlers) ResetPassword(c *gin.Context) {
	var req dtos.ResetPasswordRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	email, err := h.svcs.AuthResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if helpers.HandleError(c, err, "Failed to reset password") {
		return
	}

	helpers.OK(c, "Password reset successfully", gin.H{"email": email})
}

// TwoFASend generates and emails a 2FA code, called by api-gateway after password verification.
func (h *Handlers) TwoFASend(c *gin.Context) {
	var req dtos.TwoFASendRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	if err := h.svcs.TwoFASend(c.Request.Context(), req.Email); helpers.HandleError(c, err, "Failed to send 2FA code") {
		return
	}

	helpers.OK(c, "2FA code sent", nil)
}

// TwoFAVerify checks a submitted 2FA code and resolves the user's access.
func (h *Handlers) TwoFAVerify(c *gin.Context) {
	var req dtos.TwoFAVerifyRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.TwoFAVerify(c.Request.Context(), req.Email, req.Code)
	if helpers.HandleError(c, err, "Failed to verify 2FA code") {
		return
	}

	helpers.OK(c, "2FA code verified", result)
}
