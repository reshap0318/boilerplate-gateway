package dtos

// TwoFASendRequest is the email to send a 2FA code to.
type TwoFASendRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// TwoFAVerifyRequest is the email + code payload for completing a 2FA login.
type TwoFAVerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6,numeric"`
}

// UserUpdateTwoFARequest is the payload for enabling/disabling a user's 2FA (admin action).
type UserUpdateTwoFARequest struct {
	TwoFA bool `json:"twofa"`
}
