package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-message/internal/handlers"
)

// RegisterNotificationRoutes registers notification routes.
// POST is on the public group (GatewayPublic applied in 00_route.go) —
// created by other services, with an optional caller identity (0/"system"
// when absent). The rest act on "the caller's own notifications" and stay
// behind GatewayAuth via the protected group.
func RegisterNotificationRoutes(public, protected *gin.RouterGroup, h *handlers.Handlers) {
	public.POST("/notifications", h.NotificationCreate)

	notifications := protected.Group("/notifications")
	{
		notifications.GET("", h.NotificationGetAll)
		notifications.GET("/unread-count", h.NotificationCountUnread)
		notifications.PATCH("/read-all", h.NotificationMarkAllAsRead)
		notifications.PATCH("/:id/read", h.NotificationMarkAsRead)
		notifications.DELETE("", h.NotificationDeleteAll)
		notifications.DELETE("/:id", h.NotificationDelete)
	}
}
