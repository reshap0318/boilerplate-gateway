package dtos

import (
	"encoding/json"
	"time"

	"github.com/reshap0318/serv-message/internal/models"
)

// NotificationCreateRequest represents the request to create a notification.
// UserID is the recipient; SenderID is not part of the body — it's filled
// by the handler from the forwarded X-User-Id header (0 if absent/system).
type NotificationCreateRequest struct {
	UserID  uint            `json:"user_id" validate:"required"`
	Type    string          `json:"type" validate:"required,max=100"`
	Title   string          `json:"title" validate:"required,max=255"`
	Message string          `json:"message" validate:"required"`
	Data    json.RawMessage `json:"data"`
}

// NotificationDTO represents notification data transfer object.
type NotificationDTO struct {
	ID        uint            `json:"id"`
	UserID    uint            `json:"user_id"`
	SenderID  uint            `json:"sender_id"`
	Type      string          `json:"type"`
	Title     string          `json:"title"`
	Message   string          `json:"message"`
	Data      json.RawMessage `json:"data"`
	ReadAt    *time.Time      `json:"read_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// ToNotificationDTO converts a Notification model to NotificationDTO.
func ToNotificationDTO(n *models.Notification) NotificationDTO {
	return NotificationDTO{
		ID:        n.ID,
		UserID:    n.UserID,
		SenderID:  n.SenderID,
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		Data:      n.Data,
		ReadAt:    n.ReadAt,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

// ToNotificationDTOList converts a slice of Notification models to NotificationDTOs.
func ToNotificationDTOList(notifications []models.Notification) []NotificationDTO {
	result := make([]NotificationDTO, len(notifications))
	for i, n := range notifications {
		result[i] = ToNotificationDTO(&n)
	}
	return result
}
