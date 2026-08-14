package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-master/internal/handlers"
)

// RegisterSystemRoutes registers system-level public routes (health).
func RegisterSystemRoutes(r *gin.Engine, h *handlers.Handlers) {
	r.GET("/health", h.HealthCheck)
}
