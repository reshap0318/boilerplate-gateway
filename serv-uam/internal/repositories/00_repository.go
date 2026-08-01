package repositories

import "gorm.io/gorm"

type Repositories struct {
	DB         *gorm.DB
	TxManager  *TransactionManager
	User       *UserRepository
	Role       *RoleRepository
	Permission *PermissionRepository
	AuditLog   *AuditLogRepository
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
	return &Repositories{
		DB:         db,
		TxManager:  NewTransactionManager(db),
		User:       NewUserRepository(db),
		Role:       NewRoleRepository(db),
		Permission: NewPermissionRepository(db),
		AuditLog:   NewAuditLogRepository(db),
	}, nil
}
