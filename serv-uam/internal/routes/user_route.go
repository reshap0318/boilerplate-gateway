package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/handlers"
)

// RegisterUserRoutes registers user CRUD routes. Access control happens at
// the api-gateway (per-route required permissions) before the request ever
// reaches this service.
func RegisterUserRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	users := r.Group("/users")
	{
		users.POST("", h.UserCreate)
		users.GET("", h.UserGetAll)
		users.GET("/:id", h.UserGetByID)
		users.PUT("/:id", h.UserUpdate)
		users.POST("/:id/unlock", h.UserUnlock)
		users.DELETE("/:id", h.UserDelete)
	}
}
