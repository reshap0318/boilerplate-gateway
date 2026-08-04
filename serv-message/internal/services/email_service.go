package services

import (
	"context"
	"fmt"

	"github.com/reshap0318/serv-message/internal/dtos"
	"github.com/reshap0318/serv-message/internal/helpers"
	"github.com/reshap0318/serv-message/internal/pkg/email"
	"github.com/reshap0318/serv-message/internal/pkg/email/templates"
)

// EmailSendTemplate renders a registered template with the given params and sends it.
func (s *Services) EmailSendTemplate(ctx context.Context, req dtos.EmailTemplateRequest) error {
	def, ok := templates.Registry[req.Template]
	if !ok {
		return &helpers.FieldError{Field: "template", Message: fmt.Sprintf("unknown template: %s", req.Template)}
	}

	for _, p := range def.RequiredParams {
		if req.Params[p] == "" {
			return &helpers.FieldError{Field: "params", Message: fmt.Sprintf("missing required param: %s", p)}
		}
	}

	appName := helpers.GetEnv("APP_NAME", "Rupiah Digital")
	subject, content := def.Render(req.Params, appName)

	s.Logger.LogCtx(ctx, "EmailSendTemplate", "Sending template email %s to %v", req.Template, req.To)

	emailReq := email.EmailRequest{
		To:      req.To,
		CC:      req.CC,
		BCC:     req.BCC,
		Subject: subject,
		Body:    templates.Layout(appName, content),
	}
	go func() {
		if err := s.Email.SendEmail(emailReq); err != nil {
			s.Logger.LogCtx(ctx, "EmailSendTemplate", "Failed to send template email: %v", err)
			return
		}
		s.Logger.LogCtx(ctx, "EmailSendTemplate", "Template email sent: %v", req.To)
	}()

	return nil
}

// EmailSendCustom sends a caller-supplied HTML body wrapped in the shared layout.
func (s *Services) EmailSendCustom(ctx context.Context, req dtos.EmailCustomRequest) error {
	appName := helpers.GetEnv("APP_NAME", "Rupiah Digital")

	s.Logger.LogCtx(ctx, "EmailSendCustom", "Sending custom email to %v", req.To)

	emailReq := email.EmailRequest{
		To:      req.To,
		CC:      req.CC,
		BCC:     req.BCC,
		Subject: req.Subject,
		Body:    templates.Layout(appName, req.Body),
	}
	go func() {
		if err := s.Email.SendEmail(emailReq); err != nil {
			s.Logger.LogCtx(ctx, "EmailSendCustom", "Failed to send custom email: %v", err)
			return
		}
		s.Logger.LogCtx(ctx, "EmailSendCustom", "Custom email sent: %v", req.To)
	}()

	return nil
}
