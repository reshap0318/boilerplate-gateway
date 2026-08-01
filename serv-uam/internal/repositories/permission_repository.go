package repositories

import (
	"gorm.io/gorm"

	"github.com/reshap0318/serv-uam/internal/models"
)

// PermissionRepository extends GenericRepository with Permission-specific queries.
type PermissionRepository struct {
	*GenericRepository[models.Permission]
}

// NewPermissionRepository creates a new PermissionRepository.
func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		GenericRepository: NewGenericRepository(db, &models.Permission{}),
	}
}
