package models

import (
	"time"

	"gorm.io/gorm"
)

// GatewayRoute represents a routing rule (method + path pattern) under a GatewayService.
type GatewayRoute struct {
	ID                  uint   `gorm:"primaryKey" json:"id"`
	ServiceID           uint   `gorm:"column:service_id;not null" json:"service_id"`
	Method              string `gorm:"size:10;not null" json:"method"`
	PathPattern         string `gorm:"column:path_pattern;size:500;not null" json:"path_pattern"`
	PermissionMatchMode string `gorm:"column:permission_match_mode;size:10;not null;default:any" json:"permission_match_mode"`
	// Public defaults to false so existing/new routes stay protected unless explicitly opted
	// into being public — a route can only skip auth by deliberate choice, never by omission.
	Public             bool           `gorm:"column:public;not null;default:false" json:"public"`
	RateLimitPerMinute *int           `gorm:"column:rate_limit_per_minute" json:"rate_limit_per_minute"`
	IsActive           bool           `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
	Service            GatewayService `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	// Permission names directly, no FK — Permission now lives in serv-uam's own DB.
	Permissions []string `gorm:"column:permissions;serializer:json" json:"permissions"`
}

// TableName specifies the table name for GatewayRoute model.
func (GatewayRoute) TableName() string {
	return "gateway_routes"
}
