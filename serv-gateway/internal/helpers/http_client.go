package helpers

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HTTPClient is the shared client for all internal-service calls (e.g. to serv-uam).
var HTTPClient = &http.Client{Timeout: 10 * time.Second}

// HTTPCall sends a request to another internal service, forwarding this request's trace id
// and full caller identity from ctx (X-Trace-Id, X-User-Id/-Email/-Name/-Roles/-Permissions)
// so the callee sees the same correlated, trusted context — no separate params needed,
// everything comes from ctx.
func HTTPCall(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if traceID := GetTraceID(ctx); traceID != "" {
		req.Header.Set("X-Trace-Id", traceID)
	}
	if callerID := GetCallerID(ctx); callerID != 0 {
		req.Header.Set("X-User-Id", strconv.FormatUint(uint64(callerID), 10))
	}
	if email := GetCallerEmail(ctx); email != "" {
		req.Header.Set("X-User-Email", email)
	}
	if name := GetCallerName(ctx); name != "" {
		req.Header.Set("X-User-Name", name)
	}
	if roles := GetCallerRoles(ctx); len(roles) > 0 {
		req.Header.Set("X-User-Roles", strings.Join(roles, ","))
	}
	if permissions := GetCallerPermissions(ctx); len(permissions) > 0 {
		req.Header.Set("X-User-Permissions", strings.Join(permissions, ","))
	}
	return HTTPClient.Do(req)
}
