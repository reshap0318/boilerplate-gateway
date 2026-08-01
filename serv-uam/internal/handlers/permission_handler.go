package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/dtos"
	"github.com/reshap0318/serv-uam/internal/helpers"
	"github.com/reshap0318/serv-uam/internal/repositories"
)

// PermissionCreate creates a new permission.
func (h *Handlers) PermissionCreate(c *gin.Context) {
	var req dtos.PermissionRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.PermissionCreate(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to create permission") {
		return
	}

	helpers.Created(c, "Permission created successfully", result)
}

// PermissionGetAll returns a paginated list of permissions.
func (h *Handlers) PermissionGetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "-1"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
		SortBy:   c.DefaultQuery("sort", "id"),
		Order:    c.DefaultQuery("order", "ASC"),
	}

	if search := c.Query("search"); search != "" {
		opts.ConditionGroups = []repositories.ConditionGroup{
			{
				Logic: "OR",
				Conditions: []repositories.QueryCondition{
					{Column: "name", Operator: "LIKE", Value: "%" + search + "%"},
				},
			},
		}
	}

	result, err := h.svcs.PermissionGetAll(c.Request.Context(), opts)
	if helpers.HandleError(c, err, "Failed to fetch permissions") {
		return
	}

	helpers.OKWithMetadata(c, "Permissions retrieved successfully", result)
}

// PermissionGetByID returns a single permission by ID.
func (h *Handlers) PermissionGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	result, err := h.svcs.PermissionGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch permission") {
		return
	}

	helpers.OK(c, "Permission retrieved successfully", result)
}

// PermissionUpdate updates an existing permission.
func (h *Handlers) PermissionUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	var req dtos.PermissionRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.PermissionUpdate(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update permission") {
		return
	}

	helpers.OK(c, "Permission updated successfully", result)
}

// PermissionDelete deletes a permission.
func (h *Handlers) PermissionDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	err = h.svcs.PermissionDelete(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to delete permission") {
		return
	}

	helpers.OK(c, "Permission deleted successfully", nil)
}
