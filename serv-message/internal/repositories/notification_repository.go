package repositories

import (
	"time"

	"gorm.io/gorm"

	"github.com/reshap0318/serv-message/internal/helpers"
	"github.com/reshap0318/serv-message/internal/models"
)

// NotificationRepository provides database operations for Notification model.
type NotificationRepository struct {
	*GenericRepository[models.Notification]
}

// NewNotificationRepository creates a new NotificationRepository.
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		GenericRepository: NewGenericRepository(db, &models.Notification{}),
	}
}

// CountUnread counts unread notifications belonging to userID.
func (r *NotificationRepository) CountUnread(tx *gorm.DB, userID uint) (int64, error) {
	db := r.getDB(tx)
	var count int64
	if err := db.Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// MarkAllAsRead sets read_at = now for every unread notification belonging to userID.
func (r *NotificationRepository) MarkAllAsRead(tx *gorm.DB, userID uint) error {
	db := r.getDB(tx)
	return db.Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", time.Now()).Error
}

// DeleteOwned soft-deletes a notification by id, scoped to its owner.
// Returns helpers.ErrNotFound if no row matched (wrong id or wrong owner).
func (r *NotificationRepository) DeleteOwned(tx *gorm.DB, id, userID uint) error {
	db := r.getDB(tx)
	res := db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return helpers.ErrNotFound
	}
	return nil
}

// DeleteAllByUser soft-deletes every notification belonging to userID.
func (r *NotificationRepository) DeleteAllByUser(tx *gorm.DB, userID uint) error {
	db := r.getDB(tx)
	return db.Where("user_id = ?", userID).Delete(&models.Notification{}).Error
}
