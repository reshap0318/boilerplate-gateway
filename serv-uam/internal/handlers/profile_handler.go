package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/dtos"
	"github.com/reshap0318/serv-uam/internal/helpers"
)

// ProfileGetMe returns the caller's own profile.
func (h *Handlers) ProfileGetMe(c *gin.Context) {
	result, err := h.svcs.ProfileGetMe(c.Request.Context())
	if helpers.HandleError(c, err, "Failed to fetch profile") {
		return
	}

	helpers.OK(c, "Profile retrieved successfully", result)
}

// ProfileUpdateMe updates the caller's own profile.
func (h *Handlers) ProfileUpdateMe(c *gin.Context) {
	var req dtos.ProfileUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.ProfileUpdateMe(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to update profile") {
		return
	}

	helpers.OK(c, "Profile updated successfully", result)
}
