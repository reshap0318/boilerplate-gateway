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

// RoleCreate creates a new role and attaches the given permissions.
func (s *Services) RoleCreate(ctx context.Context, req dtos.RoleRequest) (*dtos.RoleDTO, error) {
	s.Logger.LogCtx(ctx, "RoleCreate", "Creating role: %s", req.Name)

	var result *models.Role
	err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		role := &models.Role{
			Name:        req.Name,
			Description: req.Description,
		}
		var err error
		result, err = s.repo.Role.Create(tx, role)
		if err != nil {
			return err
		}

		return s.repo.Role.SyncPermissions(tx, result, req.Permissions)
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "RoleCreate", "Failed: %v", err)
		return nil, err
	}

	result, err = s.repo.Role.FindByID(nil, result.ID, "Permissions")
	if err != nil {
		return nil, err
	}

	dto := dtos.ToRoleDTO(result)
	s.AuditLogCreate(ctx, dtos.AuditLogRequest{
		Action:      "create",
		EntityType:  "role",
		EntityID:    result.ID,
		Description: fmt.Sprintf("Created role %s", result.Name),
		Payloads:    map[string]interface{}{"after": dto},
	})

	s.Logger.LogCtx(ctx, "RoleCreate", "Role created: %s (ID: %d)", result.Name, result.ID)
	return &dto, nil
}

// RoleGetAll returns a paginated list of roles with their permissions.
func (s *Services) RoleGetAll(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.RoleDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	opts.Preloads = append(opts.Preloads, "Permissions")

	// Role ID 1 (Super Admin) is hidden from listing unless the caller holds
	// the elevated role.index-su permission.
	if !s.Access.HasPermission(ctx, "role.index-su") {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic: "AND",
			Conditions: []repositories.QueryCondition{
				{Column: "id", Operator: "!=", Value: 1},
			},
		})
	}

	result, err := s.repo.Role.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	return &repositories.PagedResult[dtos.RoleDTO]{
		Data:       dtos.ToRoleDTOList(result.Data),
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// RoleGetByID returns a single role with its permissions.
func (s *Services) RoleGetByID(ctx context.Context, id uint) (*dtos.RoleDTO, error) {
	if id == 1 && !s.Access.HasPermission(ctx, "role.index-su") {
		return nil, helpers.ErrForbidden
	}

	role, err := s.repo.Role.FindByID(nil, id, "Permissions")
	if err != nil {
		return nil, err
	}
	dto := dtos.ToRoleDTO(role)
	return &dto, nil
}

// RoleUpdate updates a role and replaces its permission set.
func (s *Services) RoleUpdate(ctx context.Context, id uint, req dtos.RoleRequest) (*dtos.RoleDTO, error) {
	s.Logger.LogCtx(ctx, "RoleUpdate", "Updating role %d", id)

	existing, err := s.repo.Role.FindByID(nil, id, "Permissions")
	if err != nil {
		return nil, err
	}

	var result *models.Role
	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.repo.Role.Update(tx, &models.Role{ID: id}, &models.Role{
			Name:        req.Name,
			Description: req.Description,
		})
		if err != nil {
			return err
		}

		return s.repo.Role.SyncPermissions(tx, result, req.Permissions)
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "RoleUpdate", "Failed: %v", err)
		return nil, err
	}

	result, err = s.repo.Role.FindByID(nil, id, "Permissions")
	if err != nil {
		return nil, err
	}

	s.invalidateAllAccess()

	dto := dtos.ToRoleDTO(result)
	beforeDTO := dtos.ToRoleDTO(existing)
	s.AuditLogCreate(ctx, dtos.AuditLogRequest{
		Action:      "update",
		EntityType:  "role",
		EntityID:    id,
		Description: fmt.Sprintf("Updated role %s", result.Name),
		Payloads:    map[string]interface{}{"before": beforeDTO, "after": dto},
	})

	s.Logger.LogCtx(ctx, "RoleUpdate", "Role %d updated", id)
	return &dto, nil
}

// RoleGetPermissions returns all permissions for a role.
func (s *Services) RoleGetPermissions(ctx context.Context, roleID uint) ([]dtos.PermissionDTO, error) {
	role, err := s.repo.Role.FindByID(nil, roleID, "Permissions")
	if err != nil {
		return nil, err
	}
	return dtos.ToPermissionDTOList(role.Permissions), nil
}

// RoleDelete deletes a role.
func (s *Services) RoleDelete(ctx context.Context, id uint) error {
	s.Logger.LogCtx(ctx, "RoleDelete", "Deleting role %d", id)

	existing, err := s.repo.Role.FindByID(nil, id, "Permissions")
	if err != nil {
		return err
	}

	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		_, err := s.repo.Role.Delete(tx, id)
		return err
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "RoleDelete", "Failed: %v", err)
		return err
	}

	s.invalidateAllAccess()

	s.AuditLogCreate(ctx, dtos.AuditLogRequest{
		Action:      "delete",
		EntityType:  "role",
		EntityID:    id,
		Description: fmt.Sprintf("Deleted role %s", existing.Name),
		Payloads:    map[string]interface{}{"before": dtos.ToRoleDTO(existing)},
	})

	s.Logger.LogCtx(ctx, "RoleDelete", "Role %d deleted", id)
	return nil
}
