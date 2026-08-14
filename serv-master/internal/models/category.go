package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a category record — the reference CRUD example for this boilerplate.
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description *string        `gorm:"size:1000" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Category model.
func (Category) TableName() string {
	return "categories"
}
