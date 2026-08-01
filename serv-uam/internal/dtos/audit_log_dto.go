package dtos

import (
	"time"

	"github.com/reshap0318/serv-uam/internal/models"
)

// AuditLogRequest is the insert payload api-gateway sends on a config change.
// The actor is read from the X-User-Id header (see helpers.HeaderUserID), not this body.
type AuditLogRequest struct {
	Action      string                 `json:"action" validate:"required"`
	EntityType  string                 `json:"entity_type" validate:"required"`
	EntityID    uint                   `json:"entity_id"`
	Description string                 `json:"description"`
	Payloads    map[string]interface{} `json:"payloads"`
}

// AuditLogDTO is the response shape for a single audit log entry.
type AuditLogDTO struct {
	ID          uint                   `json:"id"`
	ActorID     *uint                  `json:"actor_id"`
	ActorName   string                 `json:"actor_name"`
	ActorEmail  string                 `json:"actor_email"`
	Action      string                 `json:"action"`
	EntityType  string                 `json:"entity_type"`
	EntityID    uint                   `json:"entity_id"`
	Description string                 `json:"description"`
	Payloads    map[string]interface{} `json:"payloads,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// ToAuditLogDTO converts an AuditLog model to its response DTO.
func ToAuditLogDTO(a *models.AuditLog) AuditLogDTO {
	dto := AuditLogDTO{
		ID:          a.ID,
		ActorID:     a.ActorID,
		Action:      a.Action,
		EntityType:  a.EntityType,
		EntityID:    a.EntityID,
		Description: a.Description,
		Payloads:    a.Payloads,
		CreatedAt:   a.CreatedAt,
	}
	if a.Actor != nil {
		dto.ActorName = a.Actor.Name
		dto.ActorEmail = a.Actor.Email
	}
	return dto
}

// ToAuditLogDTOList converts a slice of AuditLog models to response DTOs.
func ToAuditLogDTOList(logs []models.AuditLog) []AuditLogDTO {
	dtos := make([]AuditLogDTO, len(logs))
	for i, l := range logs {
		dtos[i] = ToAuditLogDTO(&l)
	}
	return dtos
}
