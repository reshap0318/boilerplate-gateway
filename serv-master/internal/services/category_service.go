package services

import (
	"context"

	"gorm.io/gorm"

	"github.com/reshap0318/serv-master/internal/dtos"
	"github.com/reshap0318/serv-master/internal/models"
	"github.com/reshap0318/serv-master/internal/repositories"
)

// CategoryCreate creates a new category.
func (s *Services) CategoryCreate(ctx context.Context, req dtos.CategoryCreateRequest) (*dtos.CategoryDTO, error) {
	s.Logger.LogCtx(ctx, "CategoryCreate", "Creating category: %s", req.Name)

	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		return s.repo.Category.Create(tx, category)
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "CategoryCreate", "Failed to create category: %v", err)
		return nil, err
	}

	dto := dtos.ToCategoryDTO(res.(*models.Category))
	s.Logger.LogCtx(ctx, "CategoryCreate", "Category created: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

// CategoryGetAll returns paginated categories.
func (s *Services) CategoryGetAll(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.CategoryDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "ASC"
	}

	result, err := s.repo.Category.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	return &repositories.PagedResult[dtos.CategoryDTO]{
		Data:       dtos.ToCategoryDTOList(result.Data),
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// CategoryGetByID returns a category by ID.
func (s *Services) CategoryGetByID(ctx context.Context, id uint) (*dtos.CategoryDTO, error) {
	category, err := s.repo.Category.FindByID(nil, id)
	if err != nil {
		return nil, err
	}

	dto := dtos.ToCategoryDTO(category)
	return &dto, nil
}

// CategoryUpdate updates an existing category.
func (s *Services) CategoryUpdate(ctx context.Context, id uint, req dtos.CategoryUpdateRequest) (*dtos.CategoryDTO, error) {
	s.Logger.LogCtx(ctx, "CategoryUpdate", "Updating category ID: %d", id)

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
	}

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		return s.repo.Category.UpdateMap(tx, &models.Category{ID: id}, updates)
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "CategoryUpdate", "Failed to update category: %v", err)
		return nil, err
	}

	dto := dtos.ToCategoryDTO(res.(*models.Category))
	s.Logger.LogCtx(ctx, "CategoryUpdate", "Category updated: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

// CategoryDelete soft deletes a category.
func (s *Services) CategoryDelete(ctx context.Context, id uint) error {
	s.Logger.LogCtx(ctx, "CategoryDelete", "Deleting category ID: %d", id)

	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		_, err := s.repo.Category.Delete(tx, id)
		return err
	}); err != nil {
		s.Logger.LogCtx(ctx, "CategoryDelete", "Failed to delete category: %v", err)
		return err
	}

	s.Logger.LogCtx(ctx, "CategoryDelete", "Category deleted: ID %d", id)
	return nil
}
