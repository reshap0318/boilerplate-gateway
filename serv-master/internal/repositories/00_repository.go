package repositories

import "gorm.io/gorm"

type Repositories struct {
	DB        *gorm.DB
	TxManager *TransactionManager
	Category  *CategoryRepository
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
	return &Repositories{
		DB:        db,
		TxManager: NewTransactionManager(db),
		Category:  NewCategoryRepository(db),
	}, nil
}
