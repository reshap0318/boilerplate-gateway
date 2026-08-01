package helpers

import "context"

// Access checks the caller's roles/permissions, sourced from the
// api-gateway's identity headers (see middleware.GatewayAuth) and stored in
// the request context. No DB/Redis lookups — the gateway is the source of
// truth for coarse-grained role/permission checks.
//
// This is also the extension point for future business-logic access rules
// (e.g. resource ownership, fine-grained overrides) that go beyond what the
// gateway sends — none exist yet, but service-layer code can call HasRole /
// HasPermission directly when that need shows up.
type Access struct{}

// NewAccess creates a new Access instance.
func NewAccess() *Access {
	return &Access{}
}

// HasPermission checks if the caller has ANY of the specified permissions.
func (a *Access) HasPermission(ctx context.Context, permissions ...string) bool {
	callerPerms := GetCallerPermissions(ctx)
	if len(callerPerms) == 0 {
		return false
	}
	set := make(map[string]bool, len(callerPerms))
	for _, p := range callerPerms {
		set[p] = true
	}
	for _, p := range permissions {
		if set[p] {
			return true
		}
	}
	return false
}

// HasRole checks if the caller has the specified role.
func (a *Access) HasRole(ctx context.Context, role string) bool {
	for _, r := range GetCallerRoles(ctx) {
		if r == role {
			return true
		}
	}
	return false
}
