package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/reshap0318/serv-gateway/internal/dtos"
	"github.com/reshap0318/serv-gateway/internal/helpers"
)

// TwoFASend asks serv-uam to generate and email a 2FA code.
func (s *Services) TwoFASend(ctx context.Context, email string) error {
	body, err := json.Marshal(map[string]string{"email": email})
	if err != nil {
		return err
	}

	url := helpers.UamBaseURL() + "/auth/2fa/send"
	if _, err := helpers.HTTPCall(ctx, http.MethodPost, url, bytes.NewReader(body)); err != nil {
		return fmt.Errorf("uam: %w", err)
	}
	return nil
}

// TwoFAVerify completes a 2FA login: verifies the code via serv-uam then issues tokens.
func (s *Services) TwoFAVerify(ctx context.Context, email, code string) (*dtos.LoginResponse, error) {
	s.Logger.LogStart(ctx, "TwoFAVerify", "2FA verify attempt: %s", email)

	access, err := s.verifyCodeWithUAM(ctx, email, code)
	if err != nil {
		s.Logger.LogEndWithError(ctx, "TwoFAVerify", "2FA verify failed: %v", err)
		return nil, err
	}

	response, err := s.issueTokens(ctx, access)
	if err != nil {
		s.Logger.LogEndWithError(ctx, "TwoFAVerify", "2FA verify failed - token generation error")
		return nil, err
	}

	s.Logger.LogEnd(ctx, "TwoFAVerify", "2FA login successful for user: %s", email)
	return response, nil
}

// verifyCodeWithUAM checks an email/code pair against serv-uam.
func (s *Services) verifyCodeWithUAM(ctx context.Context, email, code string) (*dtos.UamAccessDTO, error) {
	body, err := json.Marshal(map[string]string{"email": email, "code": code})
	if err != nil {
		return nil, err
	}

	url := helpers.UamBaseURL() + "/auth/2fa/verify"
	data, err := helpers.HTTPCall(ctx, http.MethodPost, url, bytes.NewReader(body))

	var httpErr *helpers.HTTPError
	if errors.As(err, &httpErr) && (httpErr.Status == http.StatusUnauthorized || httpErr.Status == http.StatusTooManyRequests) {
		msg := httpErr.Message
		if msg == "" {
			msg = "Invalid or expired code"
		}
		return nil, &helpers.CustomError{Status: httpErr.Status, Message: msg}
	}
	if err != nil {
		return nil, fmt.Errorf("uam: %w", err)
	}

	var envelope struct {
		Message string            `json:"message"`
		Data    dtos.UamAccessDTO `json:"data"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("uam: decode response: %w", err)
	}
	return &envelope.Data, nil
}
