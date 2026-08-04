package dtos

// AuthVerifyRequest is the login payload forwarded by api-gateway.
type AuthVerifyRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserAccessDTO is the resolved identity/access api-gateway caches to answer permission checks.
type UserAccessDTO struct {
	UserID        uint     `json:"user_id"`
	Email         string   `json:"email"`
	Name          string   `json:"name"`
	Roles         []string `json:"roles"`
	Permissions   []string `json:"permissions"`
	TwoFARequired bool     `json:"twofa_required"`
}

// RequestPasswordResetRequest is the email to generate a reset token for.
type RequestPasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest is the token + new password payload.
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
