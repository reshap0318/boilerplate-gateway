package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/dtos"
	"github.com/reshap0318/serv-uam/internal/helpers"
	"github.com/reshap0318/serv-uam/internal/repositories"
)

// UserCreate creates a new user.
func (h *Handlers) UserCreate(c *gin.Context) {
	var req dtos.UserCreateRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.UserCreate(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to create user") {
		return
	}

	helpers.Created(c, "User created successfully", result)
}

// UserGetAll returns a paginated list of users.
func (h *Handlers) UserGetAll(c *gin.Context) {
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
					{Column: "email", Operator: "LIKE", Value: "%" + search + "%"},
				},
			},
		}
	}

	result, err := h.svcs.UserGetAll(c.Request.Context(), opts)
	if helpers.HandleError(c, err, "Failed to fetch users") {
		return
	}

	helpers.OKWithMetadata(c, "Users retrieved successfully", result)
}

// UserGetByID returns a single user by ID.
func (h *Handlers) UserGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	result, err := h.svcs.UserGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch user") {
		return
	}

	helpers.OK(c, "User retrieved successfully", result)
}

// UserUpdate updates an existing user.
func (h *Handlers) UserUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	var req dtos.UserUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.UserUpdate(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update user") {
		return
	}

	helpers.OK(c, "User updated successfully", result)
}

// UserUnlock manually clears a user's account lock.
func (h *Handlers) UserUnlock(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	result, err := h.svcs.UserUnlock(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to unlock user") {
		return
	}

	helpers.OK(c, "User unlocked successfully", result)
}

// UserDelete deletes a user.
func (h *Handlers) UserDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	err = h.svcs.UserDelete(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to delete user") {
		return
	}

	helpers.OK(c, "User deleted successfully", nil)
}
