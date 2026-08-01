package models

import "time"

// AuditLog records a config change made in api-gateway (service/route CRUD).
// Insert-only — no soft delete, no updates.
type AuditLog struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ActorID     *uint  `gorm:"column:actor_id" json:"actor_id"`
	Action      string `gorm:"size:50;not null" json:"action"`
	EntityType  string `gorm:"column:entity_type;size:100;not null" json:"entity_type"`
	EntityID    uint   `gorm:"column:entity_id" json:"entity_id"`
	Description string `gorm:"size:1000" json:"description"`
	// Payloads holds {"before": {...}, "after": {...}} — create sends after only, delete
	// sends before only, update sends both. Keys/shape are whatever the caller passed, not
	// enforced here — free-form JSON, same field for every action.
	Payloads  map[string]interface{} `gorm:"column:payloads;serializer:json" json:"payloads,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	// Actor is populated by preload — actor_id points at this same DB's users table now
	// that User lives in serv-uam, unlike the old gateway-side AuditLog whose Actor relation
	// broke once User moved out.
	Actor *User `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}

// TableName specifies the table name for AuditLog model.
func (AuditLog) TableName() string {
	return "audit_logs"
}
