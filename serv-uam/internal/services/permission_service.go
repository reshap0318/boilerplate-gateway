package services

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/reshap0318/serv-uam/internal/dtos"
	"github.com/reshap0318/serv-uam/internal/helpers"
	"github.com/reshap0318/serv-uam/internal/models"
	"github.com/reshap0318/serv-uam/internal/repositories"
)

// PermissionCreate creates a new permission.
func (s *Services) PermissionCreate(ctx context.Context, req dtos.PermissionRequest) (*dtos.PermissionDTO, error) {
	s.Logger.LogCtx(ctx, "PermissionCreate", "Creating permission: %s", req.Name)

	var result *models.Permission
	err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		permission := &models.Permission{
			Name:        req.Name,
			Description: req.Description,
		}
		var err error
		result, err = s.repo.Permission.Create(tx, permission)
		return err
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "PermissionCreate", "Failed: %v", err)
		return nil, err
	}

	dto := dtos.ToPermissionDTO(result)
	s.AuditLogCreate(ctx, dtos.AuditLogRequest{
		Action:      "create",
		EntityType:  "permission",
		EntityID:    result.ID,
		Description: fmt.Sprintf("Created permission %s", result.Name),
		Payloads:    map[string]interface{}{"after": dto},
	})

	s.Logger.LogCtx(ctx, "PermissionCreate", "Permission created: %s (ID: %d)", result.Name, result.ID)
	return &dto, nil
}

// PermissionGetAll returns a paginated list of permissions.
func (s *Services) PermissionGetAll(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.PermissionDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}

	// IDs 1-23 are the built-in RBAC permissions — hidden from listing unless
	// the caller holds the elevated role.index-su permission.
	if !s.Access.HasPermission(ctx, "role.index-su") {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic: "AND",
			Conditions: []repositories.QueryCondition{
				{Column: "id", Operator: ">", Value: 23},
			},
		})
	}

	result, err := s.repo.Permission.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	return &repositories.PagedResult[dtos.PermissionDTO]{
		Data:       dtos.ToPermissionDTOList(result.Data),
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// PermissionGetByID returns a single permission by ID.
func (s *Services) PermissionGetByID(ctx context.Context, id uint) (*dtos.PermissionDTO, error) {
	if id <= 23 && !s.Access.HasPermission(ctx, "role.index-su") {
		return nil, helpers.ErrForbidden
	}

	permission, err := s.repo.Permission.FindByID(nil, id)
	if err != nil {
		return nil, err
	}
	dto := dtos.ToPermissionDTO(permission)
	return &dto, nil
}

// PermissionUpdate updates an existing permission.
func (s *Services) PermissionUpdate(ctx context.Context, id uint, req dtos.PermissionRequest) (*dtos.PermissionDTO, error) {
	s.Logger.LogCtx(ctx, "PermissionUpdate", "Updating permission %d", id)

	existing, err := s.repo.Permission.FindByID(nil, id)
	if err != nil {
		return nil, err
	}

	var result *models.Permission
	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.repo.Permission.Update(tx, &models.Permission{ID: id}, &models.Permission{
			Name:        req.Name,
			Description: req.Description,
		})
		return err
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "PermissionUpdate", "Failed: %v", err)
		return nil, err
	}

	s.invalidateAllAccess()

	dto := dtos.ToPermissionDTO(result)
	beforeDTO := dtos.ToPermissionDTO(existing)
	s.AuditLogCreate(ctx, dtos.AuditLogRequest{
		Action:      "update",
		EntityType:  "permission",
		EntityID:    id,
		Description: fmt.Sprintf("Updated permission %s", result.Name),
		Payloads:    map[string]interface{}{"before": beforeDTO, "after": dto},
	})

	s.Logger.LogCtx(ctx, "PermissionUpdate", "Permission %d updated", id)
	return &dto, nil
}

// PermissionDelete deletes a permission.
func (s *Services) PermissionDelete(ctx context.Context, id uint) error {
	s.Logger.LogCtx(ctx, "PermissionDelete", "Deleting permission %d", id)

	existing, err := s.repo.Permission.FindByID(nil, id)
	if err != nil {
		return err
	}

	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		_, err := s.repo.Permission.Delete(tx, id)
		return err
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "PermissionDelete", "Failed: %v", err)
		return err
	}

	s.invalidateAllAccess()

	s.AuditLogCreate(ctx, dtos.AuditLogRequest{
		Action:      "delete",
		EntityType:  "permission",
		EntityID:    id,
		Description: fmt.Sprintf("Deleted permission %s", existing.Name),
		Payloads:    map[string]interface{}{"before": dtos.ToPermissionDTO(existing)},
	})

	s.Logger.LogCtx(ctx, "PermissionDelete", "Permission %d deleted", id)
	return nil
}
