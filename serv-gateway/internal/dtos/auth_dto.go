package dtos

import "time"

// LoginRequest represents the login request payload.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents the login response payload.
type LoginResponse struct {
	Token        string  `json:"token"`
	RefreshToken string  `json:"refresh_token"`
	User         UserDTO `json:"user"`
}

// RefreshTokenRequest represents the refresh token request payload.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RegisterRequest represents the register request payload.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=3"`
}

// RoleMiniDTO is a lightweight role reference, embedded in UserDTO.
type RoleMiniDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// PermissionDTO is a lightweight permission reference, embedded in UserDTO.
type PermissionDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// UserDTO represents user data transfer object.
type UserDTO struct {
	ID          uint            `json:"id"`
	Email       string          `json:"email"`
	Name        string          `json:"name"`
	Avatar      string          `json:"avatar"`
	Status      string          `json:"status"`
	LockedUntil *time.Time      `json:"locked_until"`
	CreatedAt   time.Time       `json:"created_at"`
	Roles       []RoleMiniDTO   `json:"roles"`
	Permissions []PermissionDTO `json:"permissions"`
}

// MessageResponse represents a simple message response.
type MessageResponse struct {
	Message string `json:"message"`
}

// UploadFileResponse represents the file upload response payload.
type UploadFileResponse struct {
	UUID string `json:"uuid"`
	URL  string `json:"url"`
}
