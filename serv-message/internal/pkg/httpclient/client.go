// Package httpclient provides a shared HTTP client for calling other
// internal services, forwarding the caller's identity/trace headers so the
// downstream service sees the same trusted context the api-gateway set on
// this request.
package httpclient

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/reshap0318/serv-message/internal/helpers"
)

// Client is the shared client used for all service-to-service calls.
var Client = &http.Client{Timeout: 10 * time.Second}

// Call builds and sends a request to another internal service, copying the
// X-Trace-Id / X-User-* headers from ctx (as set by middleware.TraceID and
// middleware.GatewayAuth) so the call chain stays correlated and the callee
// trusts the same caller identity.
func Call(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	forwardHeaders(req, ctx)
	return Client.Do(req)
}

func forwardHeaders(req *http.Request, ctx context.Context) {
	if traceID := helpers.GetTraceID(ctx); traceID != "" {
		req.Header.Set(helpers.TraceIDHeader, traceID)
	}
	if uid := helpers.GetCallerID(ctx); uid != 0 {
		req.Header.Set(helpers.HeaderUserID, strconv.FormatUint(uint64(uid), 10))
	}
	if email := helpers.GetCallerEmail(ctx); email != "" {
		req.Header.Set(helpers.HeaderUserEmail, email)
	}
	if name := helpers.GetCallerName(ctx); name != "" {
		req.Header.Set(helpers.HeaderUserName, name)
	}
	if roles := helpers.GetCallerRoles(ctx); len(roles) > 0 {
		req.Header.Set(helpers.HeaderUserRoles, strings.Join(roles, ","))
	}
	if perms := helpers.GetCallerPermissions(ctx); len(perms) > 0 {
		req.Header.Set(helpers.HeaderUserPermissions, strings.Join(perms, ","))
	}
}
