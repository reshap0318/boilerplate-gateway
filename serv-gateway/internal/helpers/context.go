package helpers

import "context"

// ContextKey represents a type for context keys.
type ContextKey string

const (
	// KeyUserID is the context key for storing the authenticated user's ID.
	KeyUserID ContextKey = "user_id"
	// KeyEmail is the context key for the caller's email.
	KeyEmail ContextKey = "email"
	// KeyName is the context key for the caller's name.
	KeyName ContextKey = "name"
	// KeyRoles is the context key for the caller's roles.
	KeyRoles ContextKey = "roles"
	// KeyPermissions is the context key for the caller's permissions.
	KeyPermissions ContextKey = "permissions"
	// KeyTraceID is the context key for storing the current request's trace ID.
	KeyTraceID ContextKey = "trace_id"
)

// ContextWithCaller returns a new context carrying the caller's identity from validated JWT
// claims — used by both JWTAuth (Management API) and the Dynamic Proxy Engine's authenticate()
// so every entry point populates the same context shape.
func ContextWithCaller(ctx context.Context, claims *JWTClaims) context.Context {
	ctx = context.WithValue(ctx, KeyUserID, claims.UserID)
	ctx = context.WithValue(ctx, KeyEmail, claims.Email)
	ctx = context.WithValue(ctx, KeyName, claims.Name)
	ctx = context.WithValue(ctx, KeyRoles, claims.Roles)
	ctx = context.WithValue(ctx, KeyPermissions, claims.Permissions)
	return ctx
}

// GetCallerID extracts the user ID from the context.
// Returns 0 if the user ID is not present in the context.
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

// GetTraceID extracts the trace ID from the context.
// Returns "" if the trace ID is not present in the context.
func GetTraceID(ctx context.Context) string {
	tid, ok := ctx.Value(KeyTraceID).(string)
	if !ok {
		return ""
	}
	return tid
}
