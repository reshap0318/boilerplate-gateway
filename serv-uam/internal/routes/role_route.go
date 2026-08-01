package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/handlers"
)

// RegisterRoleRoutes registers role CRUD routes. Access control happens at
// the api-gateway (per-route required permissions) before the request ever
// reaches this service.
func RegisterRoleRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	roles := r.Group("/roles")
	{
		roles.POST("", h.RoleCreate)
		roles.GET("", h.RoleGetAll)
		roles.GET("/:id", h.RoleGetByID)
		roles.PUT("/:id", h.RoleUpdate)
		roles.GET("/:id/permissions", h.RoleGetPermissions)
		roles.DELETE("/:id", h.RoleDelete)
	}
}
