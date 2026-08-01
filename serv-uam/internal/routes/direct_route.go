package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/handlers"
)

// RegisterDirectRoutes registers the routes the api-gateway calls directly on its
// own behalf (login verify, access-cache refill, audit log insert) — outside the
// protected (GatewayAuth) group since no end-user identity exists at this point.
func RegisterDirectRoutes(r *gin.Engine, h *handlers.Handlers) {
	r.POST("/auth/verify", h.AuthVerify)
	r.GET("/users/:id/access", h.UserGetAccess)
	r.POST("/auth/password-reset/request", h.RequestPasswordReset)
	r.POST("/auth/password-reset/confirm", h.ResetPassword)
	r.POST("/audit-logs", h.AuditLogCreate)
}
