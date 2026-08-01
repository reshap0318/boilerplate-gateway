package services

import "fmt"

// invalidateUserAccess clears the gateway's cached access:{userID} — call after a mutation
// that changes what a single user can do. Requires both services to share the same Redis.
func (s *Services) invalidateUserAccess(userID uint) {
	if s.RedisClient == nil || !s.RedisClient.IsCacheAvailable() {
		return
	}
	_ = s.RedisClient.Delete(fmt.Sprintf("access:%d", userID))
}

// invalidateAllAccess clears every cached access entry — for changes affecting many users.
func (s *Services) invalidateAllAccess() {
	if s.RedisClient == nil || !s.RedisClient.IsCacheAvailable() {
		return
	}
	_ = s.RedisClient.DeleteByPattern("access:*")
}

// revokeAllSessions deletes every session:{userID}:{jti} key for a user — every access and
// refresh token they currently hold stops passing api-gateway's AuthValidateToken on its very
// next use, regardless of expiry. Call on suspend: unlike invalidateUserAccess (which only
// affects the next proxied permission check), this also cuts off already-authenticated
// Management API requests. Requires both services to share the same Redis (see
// invalidateUserAccess).
func (s *Services) revokeAllSessions(userID uint) {
	if s.RedisClient == nil || !s.RedisClient.IsCacheAvailable() {
		return
	}
	_ = s.RedisClient.DeleteByPattern(fmt.Sprintf("session:%d:*", userID))
}
