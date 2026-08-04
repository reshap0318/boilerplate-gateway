package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/reshap0318/serv-gateway/internal/dtos"
	"github.com/reshap0318/serv-gateway/internal/helpers"
)

// sessionKey is the Redis presence marker for one issued token (access or refresh) — its
// existence is what AuthValidateToken checks, replacing the old blacklist-on-logout model
// with an allow-list: no entry means invalid, regardless of why.
func sessionKey(userID uint, jti string) string {
	return fmt.Sprintf("session:%d:%s", userID, jti)
}

// AuthValidateToken validates a JWT token's signature/expiry and its session entry in Redis.
func (s *Services) AuthValidateToken(tokenString string) (*helpers.JWTClaims, error) {
	claims, err := helpers.ValidateToken(tokenString, s.JWKSManager.GetPublicKey())
	if err != nil {
		return nil, helpers.ErrInvalidToken
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, helpers.ErrExpiredToken
	}

	// Fail-closed: Redis is the only place session validity is tracked now, so if it's
	// unreachable we can't confirm the token wasn't logged out — reject rather than guess.
	if !s.RedisClient.IsCacheAvailable() {
		s.Logger.LogWarn("AuthValidateToken", "Redis unavailable, rejecting token for user %d", claims.UserID)
		return nil, helpers.ErrInvalidToken
	}

	exists, err := s.RedisClient.Exists(sessionKey(claims.UserID, claims.ID))
	if err != nil || !exists {
		return nil, helpers.ErrInvalidToken
	}

	return claims, nil
}

// AuthLogin authenticates a user via serv-uam and returns tokens, or triggers a 2FA code.
func (s *Services) AuthLogin(ctx context.Context, email, password string) (*dtos.LoginResponse, error) {
	s.Logger.LogStart("AuthLogin", "User login attempt: %s", email)

	access, err := s.verifyCredentialsWithUAM(ctx, email, password)
	if err != nil {
		s.Logger.LogEndWithError("AuthLogin", "Login failed: %v", err)
		return nil, err
	}
	s.Logger.LogStep("AuthLogin", "Credentials verified: %s", email)

	if access.TwoFARequired {
		if err := s.TwoFASend(ctx, email); err != nil {
			s.Logger.LogError("AuthLogin", "Failed to request 2FA code: %v", err)
			s.Logger.LogEndWithError("AuthLogin", "Login failed - 2FA code dispatch error")
			return nil, err
		}
		s.Logger.LogEnd("AuthLogin", "2FA code sent, awaiting verification: %s", email)
		return &dtos.LoginResponse{TwoFARequired: true}, nil
	}

	response, err := s.issueTokens(access)
	if err != nil {
		s.Logger.LogEndWithError("AuthLogin", "Login failed - token generation error")
		return nil, err
	}

	s.Logger.LogEnd("AuthLogin", "Login successful for user: %s", email)
	return response, nil
}

// issueTokens generates and stores tokens for an already-verified identity.
func (s *Services) issueTokens(access *dtos.UamAccessDTO) (*dtos.LoginResponse, error) {
	token, tokenClaims, err := s.generateTokenWithClaims(access.UserID, access.Email, access.Name, access.Roles, access.Permissions)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshClaims, err := s.generateRefreshTokenWithClaims(access.UserID, access.Email, access.Name, access.Roles, access.Permissions)
	if err != nil {
		return nil, err
	}

	s.storeSession(access.UserID, tokenClaims)
	s.storeSession(access.UserID, refreshClaims)
	s.Access.Seed(access.UserID, access.Roles, access.Permissions)

	return &dtos.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         buildUserDTO(access.UserID, access.Email, access.Name, access.Roles, access.Permissions),
	}, nil
}

// AuthRefreshToken refreshes the access token using a refresh token.
func (s *Services) AuthRefreshToken(ctx context.Context, refreshToken string) (*dtos.LoginResponse, error) {
	s.Logger.LogStart("AuthRefreshToken", "Token refresh attempt")

	claims, err := s.AuthValidateToken(refreshToken)
	if err != nil {
		s.Logger.LogEndWithError("AuthRefreshToken", "Token refresh failed - invalid token")
		return nil, err
	}

	token, tokenClaims, err := s.generateTokenWithClaims(claims.UserID, claims.Email, claims.Name, claims.Roles, claims.Permissions)
	if err != nil {
		s.Logger.LogError("AuthRefreshToken", "Failed to generate token: %v", err)
		return nil, err
	}

	newRefreshToken, newRefreshClaims, err := s.generateRefreshTokenWithClaims(claims.UserID, claims.Email, claims.Name, claims.Roles, claims.Permissions)
	if err != nil {
		s.Logger.LogError("AuthRefreshToken", "Failed to generate refresh token: %v", err)
		return nil, err
	}

	s.storeSession(claims.UserID, tokenClaims)
	s.storeSession(claims.UserID, newRefreshClaims)

	s.Logger.LogEnd("AuthRefreshToken", "Token refreshed successfully for user: %s", claims.Email)
	return &dtos.LoginResponse{
		Token:        token,
		RefreshToken: newRefreshToken,
		User:         buildUserDTO(claims.UserID, claims.Email, claims.Name, claims.Roles, claims.Permissions),
	}, nil
}

// AuthLogout removes the token's session entry so it's rejected on its next use.
func (s *Services) AuthLogout(ctx context.Context, tokenString string) error {
	claims, err := helpers.ValidateToken(tokenString, s.JWKSManager.GetPublicKey())
	if err != nil {
		return helpers.ErrInvalidToken
	}

	if s.RedisClient.IsCacheAvailable() {
		_ = s.RedisClient.Delete(sessionKey(claims.UserID, claims.ID))
	}

	return nil
}

// AuthForgotPassword forwards a password-reset request to serv-uam. Errors are logged but
// swallowed into the same generic message serv-uam itself returns — this endpoint must never
// reveal whether an email is registered (see serv-uam AuthRequestPasswordReset).
func (s *Services) AuthForgotPassword(ctx context.Context, email string) (string, error) {
	body, err := json.Marshal(map[string]string{"email": email})
	if err != nil {
		return "", err
	}

	url := helpers.UamBaseURL() + "/auth/forgot-password"
	resp, err := helpers.HTTPCall(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("uam: request failed: %w", err)
	}
	defer resp.Body.Close()

	var envelope struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return "", fmt.Errorf("uam: decode response: %w", err)
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("uam: unexpected status %d", resp.StatusCode)
	}

	return envelope.Message, nil
}

// AuthResetPassword forwards a reset token + new password to serv-uam.
func (s *Services) AuthResetPassword(ctx context.Context, token, newPassword string) (string, error) {
	body, err := json.Marshal(map[string]string{"token": token, "new_password": newPassword})
	if err != nil {
		return "", err
	}

	url := helpers.UamBaseURL() + "/auth/reset-password"
	resp, err := helpers.HTTPCall(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("uam: request failed: %w", err)
	}
	defer resp.Body.Close()

	var envelope struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return "", fmt.Errorf("uam: decode response: %w", err)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		msg := envelope.Message
		if msg == "" {
			msg = "Invalid or expired reset token"
		}
		return "", &helpers.CustomError{Status: http.StatusUnauthorized, Message: msg}
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("uam: unexpected status %d", resp.StatusCode)
	}

	return envelope.Message, nil
}

// storeSession writes the session:{userID}:{jti} presence marker for a newly issued token,
// TTL'd to match its own expiry.
func (s *Services) storeSession(userID uint, claims *helpers.JWTClaims) {
	if !s.RedisClient.IsCacheAvailable() {
		s.Logger.LogWarn("storeSession", "Redis unavailable — token for user %d will fail validation", userID)
		return
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return
	}
	if err := s.RedisClient.Set(sessionKey(userID, claims.ID), "1", ttl); err != nil {
		s.Logger.LogWarn("storeSession", "Failed to store session for user %d: %v", userID, err)
	}
}

func (s *Services) generateTokenWithClaims(userID uint, email, name string, roles, permissions []string) (string, *helpers.JWTClaims, error) {
	return helpers.GenerateToken(userID, email, name, roles, permissions, s.JWKSManager.GetPrivateKey(), s.JWKSManager.GetKeyID(), helpers.GetEnvInt("JWT_EXPIRATION", 24))
}

func (s *Services) generateRefreshTokenWithClaims(userID uint, email, name string, roles, permissions []string) (string, *helpers.JWTClaims, error) {
	return helpers.GenerateRefreshToken(userID, email, name, roles, permissions, s.JWKSManager.GetPrivateKey(), s.JWKSManager.GetKeyID(), helpers.GetEnvInt("JWT_REFRESH_EXPIRATION", 168))
}

// verifyCredentialsWithUAM checks an email/password pair against serv-uam.
func (s *Services) verifyCredentialsWithUAM(ctx context.Context, email, password string) (*dtos.UamAccessDTO, error) {
	body, err := json.Marshal(map[string]string{"email": email, "password": password})
	if err != nil {
		return nil, err
	}

	url := helpers.UamBaseURL() + "/auth/verify"
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

	if resp.StatusCode == http.StatusUnauthorized {
		msg := envelope.Message
		if msg == "" {
			msg = "Invalid email or password"
		}
		return nil, &helpers.CustomError{Status: http.StatusUnauthorized, Message: msg}
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("uam: unexpected status %d", resp.StatusCode)
	}

	return &envelope.Data, nil
}

// buildUserDTO assembles the LoginResponse.User field from flat role/permission names —
// IDs/descriptions aren't available here (serv-uam's access endpoints return names only; the
// full objects live behind its own /roles, /permissions endpoints), so those are left zero.
func buildUserDTO(userID uint, email, name string, roles, permissions []string) *dtos.UserDTO {
	roleDTOs := make([]dtos.RoleMiniDTO, len(roles))
	for i, r := range roles {
		roleDTOs[i] = dtos.RoleMiniDTO{Name: r}
	}
	permDTOs := make([]dtos.PermissionDTO, len(permissions))
	for i, p := range permissions {
		permDTOs[i] = dtos.PermissionDTO{Name: p}
	}
	return &dtos.UserDTO{
		ID:          userID,
		Email:       email,
		Name:        name,
		Roles:       roleDTOs,
		Permissions: permDTOs,
	}
}
