package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-message/internal/dtos"
	"github.com/reshap0318/serv-message/internal/helpers"
	"github.com/reshap0318/serv-message/internal/repositories"
)

// NotificationCreate handles POST /notifications (public — see routes/notification_route.go)
func (h *Handlers) NotificationCreate(c *gin.Context) {
	var req dtos.NotificationCreateRequest

	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	// GatewayPublic set this from X-User-Id if present, 0 ("system") otherwise.
	senderID := helpers.GetCallerID(c.Request.Context())

	dto, err := h.svcs.NotificationCreate(c.Request.Context(), req, senderID)
	if helpers.HandleError(c, err, "Failed to create notification") {
		return
	}

	helpers.Created(c, "Notification created successfully", dto)
}

// NotificationGetAll handles GET /notifications
func (h *Handlers) NotificationGetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.NotificationGetAll(c.Request.Context(), opts)
	if helpers.HandleError(c, err, "Failed to fetch notifications") {
		return
	}

	helpers.OKWithMetadata(c, "Notifications fetched successfully", result)
}

// NotificationCountUnread handles GET /notifications/unread-count
func (h *Handlers) NotificationCountUnread(c *gin.Context) {
	count, err := h.svcs.NotificationCountUnread(c.Request.Context())
	if helpers.HandleError(c, err, "Failed to count unread notifications") {
		return
	}

	helpers.OK(c, "Unread count fetched successfully", gin.H{"count": count})
}

// NotificationMarkAsRead handles PATCH /notifications/:id/read
func (h *Handlers) NotificationMarkAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid notification ID")
		return
	}

	if err := h.svcs.NotificationMarkAsRead(c.Request.Context(), uint(id)); helpers.HandleError(c, err, "Failed to mark notification as read") {
		return
	}

	helpers.OK(c, "Notification marked as read", nil)
}

// NotificationMarkAllAsRead handles PATCH /notifications/read-all
func (h *Handlers) NotificationMarkAllAsRead(c *gin.Context) {
	if err := h.svcs.NotificationMarkAllAsRead(c.Request.Context()); helpers.HandleError(c, err, "Failed to mark notifications as read") {
		return
	}

	helpers.OK(c, "All notifications marked as read", nil)
}

// NotificationDelete handles DELETE /notifications/:id
func (h *Handlers) NotificationDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid notification ID")
		return
	}

	if err := h.svcs.NotificationDelete(c.Request.Context(), uint(id)); helpers.HandleError(c, err, "Failed to delete notification") {
		return
	}

	helpers.OK(c, "Notification deleted successfully", nil)
}

// NotificationDeleteAll handles DELETE /notifications
func (h *Handlers) NotificationDeleteAll(c *gin.Context) {
	if err := h.svcs.NotificationDeleteAll(c.Request.Context()); helpers.HandleError(c, err, "Failed to delete notifications") {
		return
	}

	helpers.OK(c, "All notifications deleted successfully", nil)
}
