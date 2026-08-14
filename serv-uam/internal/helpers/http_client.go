package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HTTPClient is the shared client used for all service-to-service calls.
var HTTPClient = &http.Client{Timeout: 10 * time.Second}

// HTTPError wraps a non-2xx response from HTTPCall. Body/Message carry the callee's payload so
// callers that need to branch on status or read its error message can, without HTTPCall's own
// []byte return value having to stay non-nil on error. Message is a best-effort pre-parse of
// the callee's `{"message": "..."}` envelope (every service in this system responds with that
// shape), empty if the body wasn't JSON or had no such field.
type HTTPError struct {
	Status  int
	Body    []byte
	Message string
}

func (e *HTTPError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("unexpected status %d: %s", e.Status, e.Message)
	}
	return fmt.Sprintf("unexpected status %d: %s", e.Status, e.Body)
}

// HTTPCall builds and sends a request to another internal service, copying the
// X-Trace-Id / X-User-* headers from ctx (as set by middleware.TraceID and
// middleware.GatewayAuth) so the call chain stays correlated and the callee
// trusts the same caller identity. The response body is closed internally and,
// on success, returned as raw bytes — decode it yourself (e.g. json.Unmarshal)
// if you need structured data. On a non-2xx status the bytes are nil and the
// error is an *HTTPError instead, which still carries the callee's response
// body (and pre-parsed Message) so nothing is lost.
func HTTPCall(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	forwardCallerHeaders(req, ctx)

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		httpErr := &HTTPError{Status: resp.StatusCode, Body: data}
		var envelope struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(data, &envelope) == nil {
			httpErr.Message = envelope.Message
		}
		return nil, httpErr
	}

	return data, nil
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
