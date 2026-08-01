package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/handlers"
)

// RegisterPermissionRoutes registers permission CRUD routes. Access control
// happens at the api-gateway (per-route required permissions) before the
// request ever reaches this service.
func RegisterPermissionRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	permissions := r.Group("/permissions")
	{
		permissions.POST("", h.PermissionCreate)
		permissions.GET("", h.PermissionGetAll)
		permissions.GET("/:id", h.PermissionGetByID)
		permissions.PUT("/:id", h.PermissionUpdate)
		permissions.DELETE("/:id", h.PermissionDelete)
	}
}
