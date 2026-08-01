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
