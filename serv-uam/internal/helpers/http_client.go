package helpers

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HTTPClient is the shared client used for all service-to-service calls.
var HTTPClient = &http.Client{Timeout: 10 * time.Second}

// HTTPCall builds and sends a request to another internal service, copying the
// X-Trace-Id / X-User-* headers from ctx (as set by middleware.TraceID and
// middleware.GatewayAuth) so the call chain stays correlated and the callee
// trusts the same caller identity.
func HTTPCall(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	forwardCallerHeaders(req, ctx)
	return HTTPClient.Do(req)
}

func forwardCallerHeaders(req *http.Request, ctx context.Context) {
	if traceID := GetTraceID(ctx); traceID != "" {
		req.Header.Set(TraceIDHeader, traceID)
	}
	if uid := GetCallerID(ctx); uid != 0 {
		req.Header.Set(HeaderUserID, strconv.FormatUint(uint64(uid), 10))
	}
	if email := GetCallerEmail(ctx); email != "" {
		req.Header.Set(HeaderUserEmail, email)
	}
	if name := GetCallerName(ctx); name != "" {
		req.Header.Set(HeaderUserName, name)
	}
	if roles := GetCallerRoles(ctx); len(roles) > 0 {
		req.Header.Set(HeaderUserRoles, strings.Join(roles, ","))
	}
	if perms := GetCallerPermissions(ctx); len(perms) > 0 {
		req.Header.Set(HeaderUserPermissions, strings.Join(perms, ","))
	}
}
