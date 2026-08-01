package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/handlers"
)

// RegisterProfileRoutes registers "my own profile" routes.
func RegisterProfileRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	r.GET("/me", h.ProfileGetMe)
	r.PUT("/me", h.ProfileUpdateMe)
}
