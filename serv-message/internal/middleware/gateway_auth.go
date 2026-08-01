package middleware

import (
	"context"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-message/internal/helpers"
)

// GatewayAuth reads the caller's identity from api-gateway headers (set
// after it has already verified the caller's token — this service trusts
// them as-is, it never validates a token itself, it only ever runs behind
// the gateway) and stores it in the request context. Roles/permissions are
// comma-separated (e.g. "admin,viewer" / "user.index,user.create"). Rejects
// the request with 401 when no valid X-User-Id is present — use this for
// routes that act on "the current user's own data".
func GatewayAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !applyGatewayHeaders(c) {
			helpers.Unauthorized(c, "Missing or invalid identity headers")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GatewayPublic reads the same identity headers as GatewayAuth when
// present, but never rejects the request when X-User-Id is missing/invalid
// — the caller id in context is simply 0. Use this for public routes that
// accept an optional caller identity (e.g. an internal create endpoint that
// records who triggered it when known, "system" otherwise) — never for
// routes reading/mutating "the current user's own data" (those must use
// GatewayAuth).
func GatewayPublic() gin.HandlerFunc {
	return func(c *gin.Context) {
		applyGatewayHeaders(c)
		c.Next()
	}
}

// applyGatewayHeaders parses the gateway identity headers, stores them in
// the request context and gin context, and returns whether X-User-Id was
// present and valid.
func applyGatewayHeaders(c *gin.Context) bool {
	userID, err := strconv.ParseUint(c.GetHeader(helpers.HeaderUserID), 10, 64)
	ok := err == nil

	email := c.GetHeader(helpers.HeaderUserEmail)
	name := c.GetHeader(helpers.HeaderUserName)
	roles := splitHeaderList(c.GetHeader(helpers.HeaderUserRoles))
	permissions := splitHeaderList(c.GetHeader(helpers.HeaderUserPermissions))

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, helpers.KeyUserID, uint(userID))
	ctx = context.WithValue(ctx, helpers.KeyEmail, email)
	ctx = context.WithValue(ctx, helpers.KeyName, name)
	ctx = context.WithValue(ctx, helpers.KeyRoles, roles)
	ctx = context.WithValue(ctx, helpers.KeyPermissions, permissions)
	c.Request = c.Request.WithContext(ctx)

	c.Set("user_id", uint(userID))
	c.Set("email", email)

	return ok
}

func splitHeaderList(v string) []string {
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
