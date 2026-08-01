package dtos

import (
	"time"

	"github.com/reshap0318/serv-uam/internal/models"
)

// ProfileUpdateRequest is the payload for a user updating their own profile.
// Password is optional — omitted means "leave unchanged".
type ProfileUpdateRequest struct {
	Name                 string `json:"name" validate:"required,max=255"`
	Email                string `json:"email" validate:"required,email,max=255"`
	Password             string `json:"password" validate:"omitempty,min=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"omitempty,eqfield=Password"`
}

// ProfileDTO is the response shape for "my own profile" — unlike UserDTO, Permissions is a
// flat, deduplicated list across all of the caller's roles (the FE profile page shows "what
// can I do", not "which role grants what"). No Avatar field — not backed by storage yet.
type ProfileDTO struct {
	ID          uint            `json:"id"`
	Email       string          `json:"email"`
	Name        string          `json:"name"`
	Avatar      *string         `json:"avatar"`
	CreatedAt   time.Time       `json:"created_at"`
	Roles       []RoleDTO       `json:"roles"`
	Permissions []PermissionDTO `json:"permissions"`
}

// ToProfileDTO converts a User model (with Roles.Permissions preloaded) to a ProfileDTO.
func ToProfileDTO(u *models.User) ProfileDTO {
	seen := make(map[uint]bool)
	permissions := make([]PermissionDTO, 0)
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if seen[perm.ID] {
				continue
			}
			seen[perm.ID] = true
			permissions = append(permissions, ToPermissionDTO(&perm))
		}
	}

	return ProfileDTO{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		Avatar:      nil,
		CreatedAt:   u.CreatedAt,
		Roles:       ToRoleDTOList(u.Roles),
		Permissions: permissions,
	}
}
