package dtos

import (
	"time"

	"github.com/reshap0318/serv-uam/internal/models"
)

// PermissionRequest is the create/update payload for a permission.
type PermissionRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description *string `json:"description" validate:"omitempty,max=255"`
}

// PermissionDTO is the response shape for a permission.
type PermissionDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToPermissionDTO converts a Permission model to its response DTO.
func ToPermissionDTO(p *models.Permission) PermissionDTO {
	return PermissionDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

// ToPermissionDTOList converts a slice of Permission models to response DTOs.
func ToPermissionDTOList(permissions []models.Permission) []PermissionDTO {
	dtos := make([]PermissionDTO, len(permissions))
	for i, p := range permissions {
		dtos[i] = ToPermissionDTO(&p)
	}
	return dtos
}
