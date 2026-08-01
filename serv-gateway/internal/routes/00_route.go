package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-gateway/internal/handlers"
	"github.com/reshap0318/serv-gateway/internal/helpers"
)

// RegisterAll registers all routes into the app.
// Add new feature routes here — main.go never needs to change.
func RegisterAll(
	r *gin.Engine,
	api *gin.RouterGroup,
	protected *gin.RouterGroup,
	h *handlers.Handlers,
	acc *helpers.Access,
) {
	// Public routes
	RegisterSystemRoutes(r, h)
	RegisterAuthRoutes(api, h)

	// Protected routes (JWT required)
	RegisterAuthProtectedRoutes(protected, h)
	RegisterSystemProtectedRoutes(protected, h)
	RegisterGatewayServiceRoutes(protected, h, acc)
	RegisterGatewayRouteRoutes(protected, h, acc)
	RegisterGatewayCacheRoutes(protected, h, acc)
}
