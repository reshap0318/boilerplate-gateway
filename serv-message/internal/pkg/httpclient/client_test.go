package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reshap0318/serv-message/internal/helpers"
)

func TestCallForwardsHeaders(t *testing.T) {
	var got http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx := context.Background()
	ctx = context.WithValue(ctx, helpers.KeyTraceID, "trace-123")
	ctx = context.WithValue(ctx, helpers.KeyUserID, uint(42))
	ctx = context.WithValue(ctx, helpers.KeyEmail, "user@example.com")
	ctx = context.WithValue(ctx, helpers.KeyRoles, []string{"admin", "viewer"})

	if _, err := Call(ctx, http.MethodGet, srv.URL, nil); err != nil {
		t.Fatalf("Call failed: %v", err)
	}

	cases := map[string]string{
		helpers.TraceIDHeader:   "trace-123",
		helpers.HeaderUserID:    "42",
		helpers.HeaderUserEmail: "user@example.com",
		helpers.HeaderUserRoles: "admin,viewer",
	}
	for header, want := range cases {
		if g := got.Get(header); g != want {
			t.Errorf("header %s = %q, want %q", header, g, want)
		}
	}

	if got.Get(helpers.HeaderUserName) != "" {
		t.Errorf("expected empty %s header when not set in context", helpers.HeaderUserName)
	}
}
