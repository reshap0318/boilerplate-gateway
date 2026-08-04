package services

import (
	"bytes"
	"context"
	"encoding/json"
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
	resp, err := helpers.HTTPCall(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("uam: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("uam: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// TwoFAVerify completes a 2FA login: verifies the code via serv-uam then issues tokens.
func (s *Services) TwoFAVerify(ctx context.Context, email, code string) (*dtos.LoginResponse, error) {
	s.Logger.LogStart("TwoFAVerify", "2FA verify attempt: %s", email)

	access, err := s.verifyCodeWithUAM(ctx, email, code)
	if err != nil {
		s.Logger.LogEndWithError("TwoFAVerify", "2FA verify failed: %v", err)
		return nil, err
	}

	response, err := s.issueTokens(access)
	if err != nil {
		s.Logger.LogEndWithError("TwoFAVerify", "2FA verify failed - token generation error")
		return nil, err
	}

	s.Logger.LogEnd("TwoFAVerify", "2FA login successful for user: %s", email)
	return response, nil
}

// verifyCodeWithUAM checks an email/code pair against serv-uam.
func (s *Services) verifyCodeWithUAM(ctx context.Context, email, code string) (*dtos.UamAccessDTO, error) {
	body, err := json.Marshal(map[string]string{"email": email, "code": code})
	if err != nil {
		return nil, err
	}

	url := helpers.UamBaseURL() + "/auth/2fa/verify"
	resp, err := helpers.HTTPCall(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("uam: request failed: %w", err)
	}
	defer resp.Body.Close()

	var envelope struct {
		Message string            `json:"message"`
		Data    dtos.UamAccessDTO `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("uam: decode response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusTooManyRequests {
		msg := envelope.Message
		if msg == "" {
			msg = "Invalid or expired code"
		}
		return nil, &helpers.CustomError{Status: resp.StatusCode, Message: msg}
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("uam: unexpected status %d", resp.StatusCode)
	}

	return &envelope.Data, nil
}
