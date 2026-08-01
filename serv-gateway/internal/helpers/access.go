package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/reshap0318/serv-gateway/internal/database"
)

// userAccessData is the shape cached at access:{userID}.
type userAccessData struct {
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

// accessCacheTTL is a safety net if Invalidate is ever missed — not the primary invalidation path.
const accessCacheTTL = 15 * time.Minute

func accessKey(userID uint) string {
	return fmt.Sprintf("access:%d", userID)
}

// UamBaseURL is where serv-uam is reachable.
func UamBaseURL() string {
	return GetEnv("UAM_SERVICE_URL", "http://localhost:8081")
}

// Access checks permissions/roles: Redis cache, serv-uam as source of truth on a miss.
type Access struct {
	redis *database.RedisCache
}

// NewAccess creates a new Access instance.
func NewAccess(redis *database.RedisCache) *Access {
	return &Access{redis: redis}
}

// getUserAccess: Redis first, serv-uam on miss, cached back for next time.
func (a *Access) getUserAccess(ctx context.Context, userID uint) (*userAccessData, bool) {
	if a.redis != nil && a.redis.IsCacheAvailable() {
		var cached userAccessData
		if err := a.redis.GetJSON(accessKey(userID), &cached); err == nil {
			return &cached, true
		}
	}

	data, err := a.fetchUserAccess(ctx, userID)
	if err != nil {
		return nil, false
	}

	if a.redis != nil && a.redis.IsCacheAvailable() {
		_ = a.redis.SetJSON(accessKey(userID), data, accessCacheTTL)
	}

	return data, true
}

// fetchUserAccess calls serv-uam's GET /users/:id/access.
func (a *Access) fetchUserAccess(ctx context.Context, userID uint) (*userAccessData, error) {
	url := UamBaseURL() + "/users/" + strconv.FormatUint(uint64(userID), 10) + "/access"
	resp, err := HTTPCall(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("uam: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("uam: unexpected status %d", resp.StatusCode)
	}

	var envelope struct {
		Data userAccessData `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("uam: decode response: %w", err)
	}

	return &envelope.Data, nil
}

// HasPermission checks if the caller has ANY of the specified permissions.
func (a *Access) HasPermission(ctx context.Context, permissions ...string) bool {
	userID := GetCallerID(ctx)
	if userID == 0 {
		return false
	}

	data, ok := a.getUserAccess(ctx, userID)
	if !ok {
		return false
	}

	for _, want := range permissions {
		for _, have := range data.Permissions {
			if have == want {
				return true
			}
		}
	}
	return false
}

// HasAllPermissions checks if the caller has ALL of the specified permissions.
// Used for GatewayRoute permission_match_mode = "all".
func (a *Access) HasAllPermissions(ctx context.Context, permissions ...string) bool {
	if len(permissions) == 0 {
		return true
	}

	userID := GetCallerID(ctx)
	if userID == 0 {
		return false
	}

	data, ok := a.getUserAccess(ctx, userID)
	if !ok {
		return false
	}

	have := make(map[string]bool, len(data.Permissions))
	for _, p := range data.Permissions {
		have[p] = true
	}
	for _, want := range permissions {
		if !have[want] {
			return false
		}
	}
	return true
}

// HasRole checks if the caller has the specified role.
func (a *Access) HasRole(ctx context.Context, role string) bool {
	userID := GetCallerID(ctx)
	if userID == 0 {
		return false
	}

	data, ok := a.getUserAccess(ctx, userID)
	if !ok {
		return false
	}

	for _, r := range data.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// GetUserAccess returns the caller's current roles and permissions (same cached source as
// HasPermission/HasRole) for callers needing the full list, not just a yes/no check.
func (a *Access) GetUserAccess(ctx context.Context) (roles, permissions []string, ok bool) {
	userID := GetCallerID(ctx)
	if userID == 0 {
		return nil, nil, false
	}

	data, ok := a.getUserAccess(ctx, userID)
	if !ok {
		return nil, nil, false
	}

	return data.Roles, data.Permissions, true
}

// Invalidate clears a user's cached access — call on role/permission change, not logout.
func (a *Access) Invalidate(userID uint) {
	if a.redis != nil && a.redis.IsCacheAvailable() {
		_ = a.redis.Delete(accessKey(userID))
	}
}

// InvalidateAll clears every cached access entry — for changes affecting many users at once.
func (a *Access) InvalidateAll() {
	if a.redis != nil && a.redis.IsCacheAvailable() {
		_ = a.redis.DeleteByPattern("access:*")
	}
}

// Seed writes roles/permissions straight to the access cache — call right after login with
// the serv-uam response already in hand, so the first permission check doesn't cache-miss on
// data we just fetched.
func (a *Access) Seed(userID uint, roles, permissions []string) {
	if a.redis != nil && a.redis.IsCacheAvailable() {
		_ = a.redis.SetJSON(accessKey(userID), &userAccessData{Roles: roles, Permissions: permissions}, accessCacheTTL)
	}
}
