package handlers

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/dtos"
	"github.com/reshap0318/serv-uam/internal/helpers"
	"github.com/reshap0318/serv-uam/internal/repositories"
)

// AuditLogCreate records a config-change reported by api-gateway —
// outside GatewayAuth since the gateway calls it on its own behalf.
func (h *Handlers) AuditLogCreate(c *gin.Context) {
	var req dtos.AuditLogRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	ctx := c.Request.Context()
	if uid, err := strconv.ParseUint(c.GetHeader(helpers.HeaderUserID), 10, 64); err == nil {
		ctx = context.WithValue(ctx, helpers.KeyUserID, uint(uid))
	}

	h.svcs.AuditLogCreate(ctx, req)

	helpers.Created(c, "Audit log recorded successfully", nil)
}

// AuditLogGetAll returns a paginated list of audit log entries.
func (h *Handlers) AuditLogGetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "-1"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
		SortBy:   c.DefaultQuery("sort", "id"),
		Order:    c.DefaultQuery("order", "DESC"),
	}

	if entityType := c.Query("entity_type"); entityType != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic:      "AND",
			Conditions: []repositories.QueryCondition{{Column: "entity_type", Operator: "=", Value: entityType}},
		})
	}

	result, err := h.svcs.AuditLogGetAll(c.Request.Context(), opts)
	if helpers.HandleError(c, err, "Failed to fetch audit logs") {
		return
	}

	helpers.OKWithMetadata(c, "Audit logs retrieved successfully", result)
}
