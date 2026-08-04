package services

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/reshap0318/serv-uam/internal/dtos"
	"github.com/reshap0318/serv-uam/internal/helpers"
	"github.com/reshap0318/serv-uam/internal/models"
)

func twoFAOTPKey(email string) string      { return "2fa_otp:" + email }
func twoFAAttemptsKey(email string) string { return "2fa_attempts:" + email }

// TwoFASend generates and emails a numeric OTP for a user with 2FA enabled.
func (s *Services) TwoFASend(ctx context.Context, email string) error {
	s.Logger.LogCtx(ctx, "TwoFASend", "2FA code requested: %s", email)

	user, err := s.repo.User.FindByEmail(nil, email)
	if err != nil {
		return helpers.ErrNotFound
	}
	if !user.TwoFA {
		return &helpers.CustomError{Status: http.StatusBadRequest, Message: "Two-factor authentication is not enabled for this account"}
	}

	if s.RedisClient == nil || !s.RedisClient.IsCacheAvailable() {
		return &helpers.CustomError{Status: http.StatusServiceUnavailable, Message: "Two-factor authentication is temporarily unavailable"}
	}

	code, err := helpers.GenerateNumericOTP(6)
	if err != nil {
		return err
	}

	ttl := time.Duration(helpers.GetEnvInt("AUTH_2FA_OTP_TTL_MINUTES", 5)) * time.Minute
	if err := s.RedisClient.Set(twoFAOTPKey(email), helpers.HashToken(code), ttl); err != nil {
		s.Logger.LogCtx(ctx, "TwoFASend", "Failed to cache code: %v", err)
		return err
	}
	_ = s.RedisClient.Delete(twoFAAttemptsKey(email))

	if err := s.sendTwoFAEmail(ctx, email, code); err != nil {
		s.Logger.LogCtx(ctx, "TwoFASend", "Failed to send code email: %v", err)
		return err
	}

	s.Logger.LogCtx(ctx, "TwoFASend", "2FA code sent: %s", email)
	return nil
}

// TwoFAVerify checks a submitted code and resolves the user's access on success.
func (s *Services) TwoFAVerify(ctx context.Context, email, code string) (*dtos.UserAccessDTO, error) {
	s.Logger.LogCtx(ctx, "TwoFAVerify", "2FA verify attempt: %s", email)

	if s.RedisClient == nil || !s.RedisClient.IsCacheAvailable() {
		return nil, &helpers.CustomError{Status: http.StatusServiceUnavailable, Message: "Two-factor authentication is temporarily unavailable"}
	}

	otpKey := twoFAOTPKey(email)
	attemptsKey := twoFAAttemptsKey(email)

	storedHash, err := s.RedisClient.Get(otpKey)
	if err != nil {
		return nil, &helpers.CustomError{Status: http.StatusUnauthorized, Message: "Invalid or expired code"}
	}

	maxAttempts := helpers.GetEnvInt("AUTH_2FA_MAX_ATTEMPTS", 5)
	attempts, err := s.RedisClient.Incr(attemptsKey)
	if err == nil && attempts == 1 {
		ttl := time.Duration(helpers.GetEnvInt("AUTH_2FA_OTP_TTL_MINUTES", 5)) * time.Minute
		_ = s.RedisClient.Expire(attemptsKey, ttl)
	}
	if attempts > int64(maxAttempts) {
		_ = s.RedisClient.Delete(otpKey, attemptsKey)
		return nil, &helpers.CustomError{Status: http.StatusTooManyRequests, Message: "Too many attempts, request a new code"}
	}

	if subtle.ConstantTimeCompare([]byte(helpers.HashToken(code)), []byte(storedHash)) != 1 {
		return nil, &helpers.CustomError{Status: http.StatusUnauthorized, Message: "Invalid or expired code"}
	}

	user, err := s.repo.User.FindByEmail(nil, email)
	if err != nil {
		return nil, helpers.ErrInvalidCredential
	}

	_ = s.RedisClient.Delete(otpKey, attemptsKey)

	s.Logger.LogCtx(ctx, "TwoFAVerify", "2FA verified: %s", email)
	return buildUserAccessDTO(user), nil
}

// TwoFAToggle enables or disables a user's 2FA — admin-only, not self-service.
func (s *Services) TwoFAToggle(ctx context.Context, id uint, enabled bool) (*dtos.UserDTO, error) {
	s.Logger.LogCtx(ctx, "TwoFAToggle", "Setting user %d twofa=%t", id, enabled)

	if _, err := s.repo.User.FindByID(nil, id); err != nil {
		return nil, err
	}

	if _, err := s.repo.User.UpdateMap(nil, &models.User{ID: id}, map[string]interface{}{"twofa": enabled}); err != nil {
		s.Logger.LogCtx(ctx, "TwoFAToggle", "Failed: %v", err)
		return nil, err
	}

	result, err := s.repo.User.FindByID(nil, id, "Roles.Permissions")
	if err != nil {
		return nil, err
	}

	s.Logger.LogCtx(ctx, "TwoFAToggle", "User %d twofa set to %t", id, enabled)
	dto := dtos.ToUserDTO(result)
	return &dto, nil
}

// sendTwoFAEmail dispatches the code via serv-message's template endpoint.
func (s *Services) sendTwoFAEmail(ctx context.Context, email, code string) error {
	body, err := json.Marshal(map[string]interface{}{
		"to":       []string{email},
		"template": "twofa_code",
		"params":   map[string]string{"code": code},
	})
	if err != nil {
		return err
	}

	url := helpers.MessageBaseURL() + "/emails/template"
	resp, err := helpers.HTTPCall(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("message: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("message: unexpected status %d", resp.StatusCode)
	}
	return nil
}
