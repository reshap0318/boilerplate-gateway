package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-message/internal/handlers"
	"github.com/reshap0318/serv-message/internal/middleware"
)

// RegisterAll registers all routes into the app.
// Add new feature routes here — main.go never needs to change.
func RegisterAll(r *gin.Engine, h *handlers.Handlers) {
	protected := r.Group("")
	protected.Use(middleware.GatewayAuth())

	public := r.Group("")
	public.Use(middleware.GatewayPublic())

	// Public routes
	RegisterSystemRoutes(r, h)
	RegisterEmailRoutes(r, h)

	// Protected routes
	RegisterNotificationRoutes(public, protected, h)
}
