package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/dtos"
	"github.com/reshap0318/serv-uam/internal/helpers"
	"github.com/reshap0318/serv-uam/internal/repositories"
)

// RoleCreate creates a new role.
func (h *Handlers) RoleCreate(c *gin.Context) {
	var req dtos.RoleRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.RoleCreate(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to create role") {
		return
	}

	helpers.Created(c, "Role created successfully", result)
}

// RoleGetAll returns a paginated list of roles.
func (h *Handlers) RoleGetAll(c *gin.Context) {
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

	result, err := h.svcs.RoleGetAll(c.Request.Context(), opts)
	if helpers.HandleError(c, err, "Failed to fetch roles") {
		return
	}

	helpers.OKWithMetadata(c, "Roles retrieved successfully", result)
}

// RoleGetByID returns a single role by ID.
func (h *Handlers) RoleGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid role ID")
		return
	}

	result, err := h.svcs.RoleGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch role") {
		return
	}

	helpers.OK(c, "Role retrieved successfully", result)
}

// RoleUpdate updates an existing role.
func (h *Handlers) RoleUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid role ID")
		return
	}

	var req dtos.RoleRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.RoleUpdate(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update role") {
		return
	}

	helpers.OK(c, "Role updated successfully", result)
}

// RoleGetPermissions returns all permissions for a role.
func (h *Handlers) RoleGetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid role ID")
		return
	}

	result, err := h.svcs.RoleGetPermissions(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch role permissions") {
		return
	}

	helpers.OK(c, "Role permissions retrieved successfully", result)
}

// RoleDelete deletes a role.
func (h *Handlers) RoleDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid role ID")
		return
	}

	err = h.svcs.RoleDelete(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to delete role") {
		return
	}

	helpers.OK(c, "Role deleted successfully", nil)
}
