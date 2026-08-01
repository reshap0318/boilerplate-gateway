package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Notification represents an in-app notification for a user.
type Notification struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `gorm:"not null;index" json:"user_id"` // recipient
	SenderID  uint            `json:"sender_id"`                     // 0 = system
	Type      string          `gorm:"size:100;not null" json:"type"`
	Title     string          `gorm:"size:255;not null" json:"title"`
	Message   string          `gorm:"type:text;not null" json:"message"`
	Data      json.RawMessage `gorm:"type:json" json:"data"` // arbitrary caller-supplied payload
	ReadAt    *time.Time      `json:"read_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName specifies the table name for Notification model.
func (Notification) TableName() string {
	return "notifications"
}
