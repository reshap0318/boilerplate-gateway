package repositories

import "gorm.io/gorm"

type Repositories struct {
	DB           *gorm.DB
	TxManager    *TransactionManager
	Notification *NotificationRepository
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
	return &Repositories{
		DB:           db,
		TxManager:    NewTransactionManager(db),
		Notification: NewNotificationRepository(db),
	}, nil
}
