package dtos

import "time"

// LoginRequest represents the login request payload.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents the login response payload.
type LoginResponse struct {
	TwoFARequired bool     `json:"twofa_required"`
	Token         string   `json:"token,omitempty"`
	RefreshToken  string   `json:"refresh_token,omitempty"`
	User          *UserDTO `json:"user,omitempty"`
}

// TwoFAVerifyRequest is the email + code payload for completing a 2FA login.
type TwoFAVerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6,numeric"`
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

// ForgotPasswordRequest represents the email to generate a reset token for.
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents the token + new password payload.
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
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
