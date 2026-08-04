package dtos

import (
	"time"

	"github.com/reshap0318/serv-uam/internal/models"
)

// UserCreateRequest is the create payload for a user. Password is required.
type UserCreateRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,max=255"`
	Status   string `json:"status" validate:"omitempty,oneof=active suspended"`
	Roles    []uint `json:"roles" validate:"omitempty,dive,required"`
}

// UserUpdateRequest is the update payload for a user. Password is optional —
// omitted means "leave unchanged".
type UserUpdateRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"omitempty,min=8"`
	Name     string `json:"name" validate:"required,max=255"`
	Status   string `json:"status" validate:"omitempty,oneof=active suspended"`
	Roles    []uint `json:"roles" validate:"omitempty,dive,required"`
}

// UserUpdateStatusRequest is the payload for suspending/activating a user.
type UserUpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active suspended"`
}

// UserDTO is the response shape for a user.
type UserDTO struct {
	ID          uint       `json:"id"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	LockedUntil *time.Time `json:"locked_until"`
	TwoFA       bool       `json:"twofa"`
	Roles       []RoleDTO  `json:"roles"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ToUserDTO converts a User model to its response DTO.
func ToUserDTO(u *models.User) UserDTO {
	return UserDTO{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		Status:      u.Status,
		LockedUntil: u.LockedUntil,
		TwoFA:       u.TwoFA,
		Roles:       ToRoleDTOList(u.Roles),
		CreatedAt:   u.CreatedAt,
	}
}

// ToUserDTOList converts a slice of User models to response DTOs.
func ToUserDTOList(users []models.User) []UserDTO {
	dtos := make([]UserDTO, len(users))
	for i, u := range users {
		dtos[i] = ToUserDTO(&u)
	}
	return dtos
}
