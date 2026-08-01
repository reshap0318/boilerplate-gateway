package services

import (
	"context"

	"github.com/reshap0318/serv-uam/internal/dtos"
	"github.com/reshap0318/serv-uam/internal/helpers"
	"github.com/reshap0318/serv-uam/internal/models"
	"github.com/reshap0318/serv-uam/internal/repositories"
)

// AuditLogCreate inserts a config-change record — called both externally (api-gateway hits
// POST /audit-logs, actor comes from the X-User-Id header the handler puts into ctx) and
// internally (role/permission services call this directly, actor comes from GatewayAuth's
// ctx). Either way the actor is read from ctx here, not passed in, so callers just parse
// ctx. The insert happens in a goroutine so every caller gets a fast response without waiting
// on the DB write — actor is nil for system-initiated changes (no caller id in ctx).
func (s *Services) AuditLogCreate(ctx context.Context, req dtos.AuditLogRequest) {
	var actorID *uint
	if uid := helpers.GetCallerID(ctx); uid != 0 {
		actorID = &uid
	}

	go func() {
		log := &models.AuditLog{
			ActorID:     actorID,
			Action:      req.Action,
			EntityType:  req.EntityType,
			EntityID:    req.EntityID,
			Description: req.Description,
			Payloads:    req.Payloads,
		}
		if _, err := s.repo.AuditLog.Create(nil, log); err != nil {
			s.Logger.LogWarn("AuditLogCreate", "Failed to insert audit log: %v", err)
		}
	}()
}

// AuditLogGetAll returns a paginated list of audit log entries, newest first.
func (s *Services) AuditLogGetAll(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.AuditLogDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "DESC"
	}
	opts.Preloads = append(opts.Preloads, "Actor")

	result, err := s.repo.AuditLog.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	return &repositories.PagedResult[dtos.AuditLogDTO]{
		Data:       dtos.ToAuditLogDTOList(result.Data),
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}
