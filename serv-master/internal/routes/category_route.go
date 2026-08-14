package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-master/internal/handlers"
)

// RegisterCategoryRoutes registers protected category routes — the reference CRUD example for this boilerplate.
// No permission middleware here — access checks (when needed) happen inside
// the service layer via Services.Access (see internal/services/00_services.go).
func RegisterCategoryRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	categories := r.Group("/categories")
	{
		categories.POST("", h.CategoryCreate)
		categories.GET("", h.CategoryGetAll)
		categories.GET("/:id", h.CategoryGetByID)
		categories.PUT("/:id", h.CategoryUpdate)
		categories.DELETE("/:id", h.CategoryDelete)
	}
}
