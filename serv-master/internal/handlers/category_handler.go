package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-master/internal/dtos"
	"github.com/reshap0318/serv-master/internal/helpers"
	"github.com/reshap0318/serv-master/internal/repositories"
)

// CategoryCreate handles POST /categories
func (h *Handlers) CategoryCreate(c *gin.Context) {
	var req dtos.CategoryCreateRequest

	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.CategoryCreate(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to create category") {
		return
	}

	helpers.Created(c, "Category created successfully", dto)
}

// CategoryGetAll handles GET /categories
// page_size defaults to -1 (no pagination, return everything). Pass a
// positive page_size to paginate.
func (h *Handlers) CategoryGetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "-1"))

	if pageSize < 0 {
		page = 1
	}

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.CategoryGetAll(c.Request.Context(), opts)
	if helpers.HandleError(c, err, "Failed to fetch categories") {
		return
	}

	helpers.OKWithMetadata(c, "Categories fetched successfully", result)
}

// CategoryGetByID handles GET /categories/:id
func (h *Handlers) CategoryGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid category ID")
		return
	}

	dto, err := h.svcs.CategoryGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch category") {
		return
	}

	helpers.OK(c, "Category fetched successfully", dto)
}

// CategoryUpdate handles PUT /categories/:id
func (h *Handlers) CategoryUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid category ID")
		return
	}

	var req dtos.CategoryUpdateRequest

	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.CategoryUpdate(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update category") {
		return
	}

	helpers.OK(c, "Category updated successfully", dto)
}

// CategoryDelete handles DELETE /categories/:id
func (h *Handlers) CategoryDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid category ID")
		return
	}

	if err := h.svcs.CategoryDelete(c.Request.Context(), uint(id)); err != nil {
		if helpers.HandleError(c, err, "Failed to delete category") {
			return
		}
	}

	helpers.OK(c, "Category deleted successfully", nil)
}
