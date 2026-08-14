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

// HTTPClient is the shared client for all internal-service calls (e.g. to serv-uam).
var HTTPClient = &http.Client{Timeout: 10 * time.Second}

// HTTPError wraps a non-2xx response from HTTPCall. Body/Message carry the callee's payload
// so callers that need to branch on status (e.g. 401 → surface the callee's own error message)
// or decode its error payload can do so — see auth_service.go — without HTTPCall's own []byte
// return value having to stay non-nil on error. Message is a best-effort pre-parse of the
// callee's `{"message": "..."}` envelope (every service in this system responds with that
// shape — see serv-uam's response_helper.go), empty if the body wasn't JSON or had no such
// field.
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

// HTTPCall sends a request to another internal service, forwarding this request's trace id
// and full caller identity from ctx (X-Trace-Id, X-User-Id/-Email/-Name/-Roles/-Permissions)
// so the callee sees the same correlated, trusted context — no separate params needed,
// everything comes from ctx. The response body is closed internally and, on success, returned
// as raw bytes — decode it yourself (e.g. json.Unmarshal) if you need structured data. On a
// non-2xx status the bytes are nil and the error is an *HTTPError instead, which still carries
// the callee's response body (and pre-parsed Message) so nothing is lost.
func HTTPCall(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
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
