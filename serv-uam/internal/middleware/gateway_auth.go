package middleware

import (
	"context"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/serv-uam/internal/helpers"
)

// GatewayAuth reads the caller's identity from api-gateway headers (set
// after it has already verified the caller's token — this service trusts
// them as-is, it never validates a token itself, it only ever runs behind
// the gateway) and stores it in the request context. Roles/permissions are
// comma-separated (e.g. "admin,viewer" / "user.index,user.create").
func GatewayAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.GetHeader(helpers.HeaderUserID), 10, 64)
		if err != nil {
			helpers.Unauthorized(c, "Missing or invalid identity headers")
			c.Abort()
			return
		}

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

		c.Next()
	}
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
