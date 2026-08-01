package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/reshap0318/serv-uam/internal/helpers"
	"github.com/reshap0318/serv-uam/internal/models"
)

// UserRepository extends GenericRepository with User-specific queries.
type UserRepository struct {
	*GenericRepository[models.User]
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		GenericRepository: NewGenericRepository(db, &models.User{}),
	}
}

// FindByEmail finds a user by email, preloading roles and permissions —
// used for login (password check) and access resolution.
func (r *UserRepository) FindByEmail(tx *gorm.DB, email string) (*models.User, error) {
	db := r.getDB(tx)

	var user models.User
	err := db.Preload("Roles.Permissions").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// SyncRoles replaces a user's role set with exactly the given IDs. Uses
// GORM's Association API (not Create/Update) so it only touches the
// user_has_roles join rows — never the roles table itself.
func (r *UserRepository) SyncRoles(tx *gorm.DB, user *models.User, roleIDs []uint) error {
	db := r.getDB(tx)

	roles := make([]models.Role, len(roleIDs))
	for i, id := range roleIDs {
		roles[i] = models.Role{ID: id}
	}

	return db.Model(user).Association("Roles").Replace(roles)
}
