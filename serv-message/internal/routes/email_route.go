package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-message/internal/handlers"
)

// RegisterEmailRoutes registers public email-send routes — hit by other
// services to dispatch a transactional email. No auth middleware: there is
// no service-to-service auth mechanism in this project, so these are
// protected by network placement only (reachable via gateway/internal
// network), same as system routes.
func RegisterEmailRoutes(r *gin.Engine, h *handlers.Handlers) {
	emails := r.Group("/emails")
	{
		emails.POST("/template", h.EmailSendTemplate)
		emails.POST("/custom", h.EmailSendCustom)
	}
}
