package repositories

import (
	"github.com/reshap0318/serv-master/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository provides database operations for Category model.
type CategoryRepository struct {
	*GenericRepository[models.Category]
}

// NewCategoryRepository creates a new CategoryRepository.
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		GenericRepository: NewGenericRepository(db, &models.Category{}),
	}
}
