package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-gateway/internal/dtos"
	"github.com/reshap0318/serv-gateway/internal/helpers"
)

// AuthLogin handles user login.
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.LoginRequest true "Login credentials"
// @Success 200 {object} dtos.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func (h *Handlers) AuthLogin(c *gin.Context) {
	var req dtos.LoginRequest

	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	response, err := h.svcs.AuthLogin(c.Request.Context(), req.Email, req.Password)
	if helpers.HandleError(c, err, "Login failed") {
		return
	}

	helpers.OK(c, "Login successful", response)
}

// AuthRefreshToken handles token refresh.
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dtos.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/refresh [post]
func (h *Handlers) AuthRefreshToken(c *gin.Context) {
	var req dtos.RefreshTokenRequest

	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	response, err := h.svcs.AuthRefreshToken(c.Request.Context(), req.RefreshToken)
	if helpers.HandleError(c, err, "Token refresh failed") {
		return
	}

	helpers.OK(c, "Token refreshed successfully", response)
}

// AuthForgotPassword handles a password-reset request, forwarded to serv-uam.
// @Summary Request password reset
// @Description Generates a password reset token and emails it (never reveals if the email exists)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.ForgotPasswordRequest true "Email"
// @Success 200 {object} map[string]string
// @Router /api/auth/forgot-password [post]
func (h *Handlers) AuthForgotPassword(c *gin.Context) {
	var req dtos.ForgotPasswordRequest

	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	message, err := h.svcs.AuthForgotPassword(c.Request.Context(), req.Email)
	if helpers.HandleError(c, err, "Failed to request password reset") {
		return
	}

	helpers.OK(c, message, nil)
}

// AuthResetPassword handles a password reset via token, forwarded to serv-uam.
// @Summary Reset password
// @Description Validates a reset token and updates the user's password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.ResetPasswordRequest true "Token + new password"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/reset-password [post]
func (h *Handlers) AuthResetPassword(c *gin.Context) {
	var req dtos.ResetPasswordRequest

	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	message, err := h.svcs.AuthResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if helpers.HandleError(c, err, "Failed to reset password") {
		return
	}

	helpers.OK(c, message, nil)
}

// AuthLogout handles user logout.
// @Summary User logout
// @Description Blacklists the current token so it cannot be reused
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /api/auth/logout [post]
func (h *Handlers) AuthLogout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	tokenString := parts[1]

	if err := h.svcs.AuthLogout(c.Request.Context(), tokenString); err != nil {
		if helpers.HandleError(c, err, "Logout failed") {
			return
		}
	}

	helpers.OK(c, "Logout successful", nil)
}
