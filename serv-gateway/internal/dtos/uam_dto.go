package dtos

// UamAccessDTO is serv-uam's /auth/verify and /users/:id/access response shape — the resolved
// identity + roles/permissions used to sign a JWT and seed the Access cache.
type UamAccessDTO struct {
	UserID      uint     `json:"user_id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
