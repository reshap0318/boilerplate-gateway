package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-gateway/internal/handlers"
)

// RegisterAuthRoutes registers public authentication routes. Forgot/reset-password aren't
// here — they're public GatewayRoute entries proxied straight to serv-uam (see
// cmd/migration/seeders/gateway_seeder.go), not part of the Management API.
func RegisterAuthRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.AuthLogin)
		auth.POST("/refresh", h.AuthRefreshToken)
	}
}

// RegisterAuthProtectedRoutes registers protected authentication routes.
func RegisterAuthProtectedRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	auth := r.Group("/auth")
	{
		auth.POST("/logout", h.AuthLogout)
	}
}
