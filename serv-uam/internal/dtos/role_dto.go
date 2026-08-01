package dtos

import (
	"time"

	"github.com/reshap0318/serv-uam/internal/models"
)

// RoleRequest is the create/update payload for a role.
type RoleRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description *string `json:"description" validate:"omitempty,max=255"`
	Permissions []uint  `json:"permissions" validate:"omitempty,dive,required"`
}

// RoleDTO is the response shape for a role.
type RoleDTO struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Permissions []PermissionDTO `json:"permissions"`
	CreatedAt   time.Time       `json:"created_at"`
}

// ToRoleDTO converts a Role model to its response DTO.
func ToRoleDTO(r *models.Role) RoleDTO {
	return RoleDTO{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: ToPermissionDTOList(r.Permissions),
		CreatedAt:   r.CreatedAt,
	}
}

// ToRoleDTOList converts a slice of Role models to response DTOs.
func ToRoleDTOList(roles []models.Role) []RoleDTO {
	dtos := make([]RoleDTO, len(roles))
	for i, r := range roles {
		dtos[i] = ToRoleDTO(&r)
	}
	return dtos
}
