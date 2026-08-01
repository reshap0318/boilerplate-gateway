package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/handlers"
	"github.com/reshap0318/serv-uam/internal/middleware"
)

// RegisterAll registers all routes into the app.
// Add new feature routes here — main.go never needs to change.
func RegisterAll(r *gin.Engine, h *handlers.Handlers) {
	protected := r.Group("")
	protected.Use(middleware.GatewayAuth())

	// Public routes
	RegisterSystemRoutes(r, h)

	// Routes the api-gateway calls directly on its own behalf, outside GatewayAuth.
	RegisterDirectRoutes(r, h)

	// Protected routes (api-gateway identity headers via middleware.GatewayAuth)
	RegisterUserRoutes(protected, h)
	RegisterRoleRoutes(protected, h)
	RegisterPermissionRoutes(protected, h)
	RegisterAuditLogRoutes(protected, h)
}
