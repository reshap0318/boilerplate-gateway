package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-message/internal/dtos"
	"github.com/reshap0318/serv-message/internal/helpers"
)

// EmailSendTemplate handles POST /emails/template
func (h *Handlers) EmailSendTemplate(c *gin.Context) {
	var req dtos.EmailTemplateRequest

	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	if err := h.svcs.EmailSendTemplate(c.Request.Context(), req); helpers.HandleError(c, err, "Failed to send email") {
		return
	}

	helpers.OK(c, "Email sent successfully", nil)
}

// EmailSendCustom handles POST /emails/custom
func (h *Handlers) EmailSendCustom(c *gin.Context) {
	var req dtos.EmailCustomRequest

	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	if err := h.svcs.EmailSendCustom(c.Request.Context(), req); helpers.HandleError(c, err, "Failed to send email") {
		return
	}

	helpers.OK(c, "Email sent successfully", nil)
}
