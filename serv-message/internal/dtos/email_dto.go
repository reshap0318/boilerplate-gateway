package dtos

// EmailTemplateRequest sends an email rendered from a registered template
// (see internal/pkg/email/templates.Registry). Params are validated against
// the template's required params in the service layer.
type EmailTemplateRequest struct {
	To       []string          `json:"to" validate:"required,min=1,dive,email"`
	CC       []string          `json:"cc" validate:"omitempty,dive,email"`
	BCC      []string          `json:"bcc" validate:"omitempty,dive,email"`
	Template string            `json:"template" validate:"required"`
	Params   map[string]string `json:"params"`
}

// EmailCustomRequest sends a caller-supplied HTML body wrapped in the
// shared email layout (header/footer).
type EmailCustomRequest struct {
	To      []string `json:"to" validate:"required,min=1,dive,email"`
	CC      []string `json:"cc" validate:"omitempty,dive,email"`
	BCC     []string `json:"bcc" validate:"omitempty,dive,email"`
	Subject string   `json:"subject" validate:"required"`
	Body    string   `json:"body" validate:"required"`
}
