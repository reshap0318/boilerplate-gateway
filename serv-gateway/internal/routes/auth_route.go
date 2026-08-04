package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-gateway/internal/handlers"
)

// RegisterAuthRoutes registers public authentication routes. Forgot/reset-password are
// handled here directly (forwarded to serv-uam server-to-server, see
// services/auth_service.go) rather than through the Dynamic Proxy Engine, like login/refresh.
func RegisterAuthRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.AuthLogin)
		auth.POST("/2fa/verify", h.TwoFAVerify)
		auth.POST("/refresh", h.AuthRefreshToken)
		auth.POST("/forgot-password", h.AuthForgotPassword)
		auth.POST("/reset-password", h.AuthResetPassword)
	}
}

// RegisterAuthProtectedRoutes registers protected authentication routes.
func RegisterAuthProtectedRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	auth := r.Group("/auth")
	{
		auth.POST("/logout", h.AuthLogout)
	}
}
