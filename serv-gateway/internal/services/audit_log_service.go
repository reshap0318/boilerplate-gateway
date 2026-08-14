package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/reshap0318/serv-gateway/internal/helpers"
)

// recordAuditLog reports a config change to serv-uam — best-effort: a failed audit call must
// never fail the request that already succeeded. serv-uam's AuditLogCreate responds
// immediately and inserts in its own goroutine, so this stays a plain synchronous call rather
// than needing its own goroutine here too. helpers.HTTPCall forwards the caller id from ctx
// as X-User-Id automatically (serv-uam reads the actor from that header, not the body).
func (s *Services) recordAuditLog(ctx context.Context, action, entityType string, entityID uint, description string, payloads map[string]interface{}) {
	body, err := json.Marshal(map[string]interface{}{
		"action":      action,
		"entity_type": entityType,
		"entity_id":   entityID,
		"description": description,
		"payloads":    payloads,
	})
	if err != nil {
		s.Logger.LogWarn(ctx, "recordAuditLog", "Failed to marshal payload: %v", err)
		return
	}

	url := helpers.UamBaseURL() + "/audit-logs"
	if _, err := helpers.HTTPCall(ctx, http.MethodPost, url, bytes.NewReader(body)); err != nil {
		s.Logger.LogWarn(ctx, "recordAuditLog", "Request failed: %v", err)
	}
}
