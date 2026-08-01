package repositories

import (
	"gorm.io/gorm"

	"github.com/reshap0318/serv-uam/internal/models"
)

// RoleRepository extends GenericRepository with Role-specific queries.
type RoleRepository struct {
	*GenericRepository[models.Role]
}

// NewRoleRepository creates a new RoleRepository.
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		GenericRepository: NewGenericRepository(db, &models.Role{}),
	}
}

// SyncPermissions replaces a role's permission set with exactly the given
// IDs. Uses GORM's Association API (not Create/Update) so it only touches
// the role_has_permissions join rows — never the permissions table itself.
func (r *RoleRepository) SyncPermissions(tx *gorm.DB, role *models.Role, permissionIDs []uint) error {
	db := r.getDB(tx)

	permissions := make([]models.Permission, len(permissionIDs))
	for i, id := range permissionIDs {
		permissions[i] = models.Permission{ID: id}
	}

	return db.Model(role).Association("Permissions").Replace(permissions)
}
