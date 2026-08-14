package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/reshap0318/serv-uam/internal/dtos"
	"github.com/reshap0318/serv-uam/internal/helpers"
	"github.com/reshap0318/serv-uam/internal/models"
)

// AuthVerify checks an email/password pair for the api-gateway's login flow
// and returns the resolved identity + roles/permissions to cache. There is
// no end-user identity yet at this point, so this is called directly by
// api-gateway on its own behalf, not on behalf of a forwarded request.
func (s *Services) AuthVerify(ctx context.Context, req dtos.AuthVerifyRequest) (*dtos.UserAccessDTO, error) {
	s.Logger.LogCtx(ctx, "AuthVerify", "Login attempt: %s", req.Email)

	user, err := s.repo.User.FindByEmail(nil, req.Email)
	if err != nil {
		s.Logger.LogCtx(ctx, "AuthVerify", "User not found: %s", req.Email)
		return nil, helpers.ErrInvalidCredential
	}

	if user.Status != "active" {
		s.Logger.LogCtx(ctx, "AuthVerify", "Account not active: %s", req.Email)
		return nil, &helpers.CustomError{Status: http.StatusUnauthorized, Message: "Account is not active"}
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		s.Logger.LogCtx(ctx, "AuthVerify", "Account locked: %s", req.Email)
		return nil, &helpers.CustomError{Status: http.StatusUnauthorized, Message: "Account temporarily locked, try again later"}
	}

	if err := helpers.VerifyString(req.Password, user.Password); err != nil {
		s.Logger.LogCtx(ctx, "AuthVerify", "Invalid password: %s", req.Email)
		s.registerFailedLogin(user)
		return nil, helpers.ErrInvalidCredential
	}

	if user.FailedLoginAttempts > 0 {
		if _, err := s.repo.User.UpdateMap(nil, &models.User{ID: user.ID}, map[string]interface{}{"failed_login_attempts": 0}); err != nil {
			s.Logger.LogWarn("AuthVerify", "Failed to reset failed_login_attempts for user %d: %v", user.ID, err)
		}
	}

	s.Logger.LogCtx(ctx, "AuthVerify", "Login successful: %s", req.Email)
	access := buildUserAccessDTO(user)
	access.TwoFARequired = user.TwoFA
	return access, nil
}

// UserGetAccess resolves a user's current roles/permissions by ID — used by
// the api-gateway to refill its Redis cache on a miss.
func (s *Services) UserGetAccess(ctx context.Context, userID uint) (*dtos.UserAccessDTO, error) {
	user, err := s.repo.User.FindByID(nil, userID, "Roles.Permissions")
	if err != nil {
		return nil, err
	}
	return buildUserAccessDTO(user), nil
}

// passwordResetKey is where a reset token's owning email is cached — TTL'd, so expiry is
// Redis's job, and deleted on use, so single-use is a plain DEL instead of a `used` column.
func passwordResetKey(hashedToken string) string {
	return "password_reset:" + hashedToken
}

const passwordResetTTL = 1 * time.Hour

// AuthRequestPasswordReset generates a reset token and caches it against the email, keyed by
// the token's hash. Returns nothing to the caller beyond success/failure — the token must only
// ever reach the account owner via email, never in an API response (that would let anyone
// reset anyone else's password without ever touching their inbox).
func (s *Services) AuthRequestPasswordReset(ctx context.Context, email string) error {
	s.Logger.LogCtx(ctx, "AuthRequestPasswordReset", "Reset requested for: %s", email)

	user, err := s.repo.User.FindByEmail(nil, email)
	if err != nil {
		s.Logger.LogCtx(ctx, "AuthRequestPasswordReset", "User not found: %s", email)
		return helpers.ErrNotFound
	}

	if s.RedisClient == nil || !s.RedisClient.IsCacheAvailable() {
		return &helpers.CustomError{Status: http.StatusServiceUnavailable, Message: "Password reset is temporarily unavailable"}
	}

	token, err := helpers.GenerateRandomString(32)
	if err != nil {
		return err
	}

	if err := s.RedisClient.Set(passwordResetKey(helpers.HashToken(token)), user.Email, passwordResetTTL); err != nil {
		s.Logger.LogCtx(ctx, "AuthRequestPasswordReset", "Failed: %v", err)
		return err
	}

	if err := s.sendResetPasswordEmail(ctx, user.Email, token); err != nil {
		s.Logger.LogCtx(ctx, "AuthRequestPasswordReset", "Failed to send reset email: %v", err)
		return err
	}

	s.Logger.LogCtx(ctx, "AuthRequestPasswordReset", "Reset token generated for: %s", email)
	return nil
}

// sendResetPasswordEmail dispatches the reset link via serv-message's template endpoint.
func (s *Services) sendResetPasswordEmail(ctx context.Context, email, token string) error {
	resetURL := helpers.FEBaseURL() + "/reset-password?token=" + token

	body, err := json.Marshal(map[string]interface{}{
		"to":       []string{email},
		"template": "reset_password",
		"params":   map[string]string{"token": token, "reset_url": resetURL},
	})
	if err != nil {
		return err
	}

	url := helpers.MessageBaseURL() + "/emails/template"
	if _, err := helpers.HTTPCall(ctx, http.MethodPost, url, bytes.NewReader(body)); err != nil {
		return fmt.Errorf("message: %w", err)
	}
	return nil
}

// AuthResetPassword validates a reset token and updates the user's password. Returns the
// affected email so the caller (api-gateway) can use it, e.g. for its own notification.
func (s *Services) AuthResetPassword(ctx context.Context, token, newPassword string) (string, error) {
	s.Logger.LogCtx(ctx, "AuthResetPassword", "Password reset attempt")

	if s.RedisClient == nil || !s.RedisClient.IsCacheAvailable() {
		return "", &helpers.CustomError{Status: http.StatusServiceUnavailable, Message: "Password reset is temporarily unavailable"}
	}

	key := passwordResetKey(helpers.HashToken(token))
	email, err := s.RedisClient.Get(key)
	if err != nil {
		return "", &helpers.CustomError{Status: http.StatusUnauthorized, Message: "Invalid or expired reset token"}
	}

	hashedPassword, err := helpers.HashString(newPassword)
	if err != nil {
		return "", err
	}

	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		_, err := s.repo.User.Update(tx, &models.User{Email: email}, &models.User{Password: hashedPassword})
		return err
	}); err != nil {
		s.Logger.LogCtx(ctx, "AuthResetPassword", "Failed: %v", err)
		return "", err
	}

	_ = s.RedisClient.Delete(key)

	s.Logger.LogCtx(ctx, "AuthResetPassword", "Password reset successful")
	return email, nil
}

// registerFailedLogin increments the failed-login counter and locks the account once the
// configured threshold is reached. Errors are logged only — a failure here must never surface
// as a different error to the caller than the plain "invalid credential".
func (s *Services) registerFailedLogin(user *models.User) {
	maxAttempts := helpers.GetEnvInt("AUTH_MAX_FAILED_LOGIN_ATTEMPTS", 5)
	lockDurationMin := helpers.GetEnvInt("AUTH_LOCK_DURATION_MINUTES", 15)

	attempts := user.FailedLoginAttempts + 1
	update := map[string]interface{}{"failed_login_attempts": attempts}

	if attempts >= maxAttempts {
		lockedUntil := time.Now().Add(time.Duration(lockDurationMin) * time.Minute)
		update["locked_until"] = lockedUntil
		s.Logger.LogWarn("registerFailedLogin", "User %d locked until %s after %d failed attempts", user.ID, lockedUntil, attempts)
	}

	if _, err := s.repo.User.UpdateMap(nil, &models.User{ID: user.ID}, update); err != nil {
		s.Logger.LogWarn("registerFailedLogin", "Failed to update failed_login_attempts for user %d: %v", user.ID, err)
	}
}

// buildUserAccessDTO collapses a user's roles into deduplicated role/permission name lists.
func buildUserAccessDTO(user *models.User) *dtos.UserAccessDTO {
	roleSet := make(map[string]bool)
	permSet := make(map[string]bool)

	for _, r := range user.Roles {
		roleSet[r.Name] = true
		for _, p := range r.Permissions {
			permSet[p.Name] = true
		}
	}

	roles := make([]string, 0, len(roleSet))
	for name := range roleSet {
		roles = append(roles, name)
	}
	permissions := make([]string, 0, len(permSet))
	for name := range permSet {
		permissions = append(permissions, name)
	}

	return &dtos.UserAccessDTO{
		UserID:      user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Roles:       roles,
		Permissions: permissions,
	}
}
