package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-master/internal/handlers"
	"github.com/reshap0318/serv-master/internal/middleware"
)

// RegisterAll registers all routes into the app.
// Add new feature routes here — main.go never needs to change.
func RegisterAll(
	r *gin.Engine,
	h *handlers.Handlers,
) {
	// Public routes — no identity required
	RegisterSystemRoutes(r, h)

	// Public routes with optional identity — add here when needed:
	// publicOptional := r.Group("")
	// publicOptional.Use(middleware.GatewayPublic())

	// Protected routes — identity required
	protected := r.Group("")
	protected.Use(middleware.GatewayAuth())
	RegisterCategoryRoutes(protected, h)
}
