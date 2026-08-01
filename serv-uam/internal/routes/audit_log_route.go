package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/handlers"
)

// RegisterAuditLogRoutes registers the read side of audit logs. Insert happens
// via RegisterDirectRoutes (api-gateway calls it directly, not through this
// protected group) — this is list-only.
func RegisterAuditLogRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	r.GET("/audit-logs", h.AuditLogGetAll)
}
