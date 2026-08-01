package helpers

import "context"

// ContextKey represents a type for context keys.
type ContextKey string

const (
	// KeyUserID is the context key for the caller's ID (set by middleware.GatewayAuth from api-gateway headers).
	KeyUserID ContextKey = "user_id"
	// KeyEmail is the context key for the caller's email.
	KeyEmail ContextKey = "email"
	// KeyName is the context key for the caller's name.
	KeyName ContextKey = "name"
	// KeyRoles is the context key for the caller's roles.
	KeyRoles ContextKey = "roles"
	// KeyPermissions is the context key for the caller's permissions.
	KeyPermissions ContextKey = "permissions"
	// KeyTraceID is the context key for the request's trace id (set by middleware.TraceID).
	KeyTraceID ContextKey = "trace_id"
)

// Identity/trace headers the api-gateway sets on every request, and that
// service-to-service calls (see HTTPCall) forward downstream so the
// callee trusts the same caller context. Kept alongside the context keys
// above since they're two representations of the same wire contract.
const (
	HeaderUserID          = "X-User-Id"
	HeaderUserEmail       = "X-User-Email"
	HeaderUserName        = "X-User-Name"
	HeaderUserRoles       = "X-User-Roles"
	HeaderUserPermissions = "X-User-Permissions"
	TraceIDHeader         = "X-Trace-Id"
)

// GetCallerID extracts the caller's user ID from the context.
// Returns 0 if not present.
func GetCallerID(ctx context.Context) uint {
	uid, ok := ctx.Value(KeyUserID).(uint)
	if !ok {
		return 0
	}
	return uid
}

// GetCallerEmail extracts the caller's email from the context.
func GetCallerEmail(ctx context.Context) string {
	email, _ := ctx.Value(KeyEmail).(string)
	return email
}

// GetCallerName extracts the caller's name from the context.
func GetCallerName(ctx context.Context) string {
	name, _ := ctx.Value(KeyName).(string)
	return name
}

// GetCallerRoles extracts the caller's roles from the context.
func GetCallerRoles(ctx context.Context) []string {
	roles, _ := ctx.Value(KeyRoles).([]string)
	return roles
}

// GetCallerPermissions extracts the caller's permissions from the context.
func GetCallerPermissions(ctx context.Context) []string {
	permissions, _ := ctx.Value(KeyPermissions).([]string)
	return permissions
}

// GetTraceID extracts the request's trace id from the context.
func GetTraceID(ctx context.Context) string {
	traceID, _ := ctx.Value(KeyTraceID).(string)
	return traceID
}
