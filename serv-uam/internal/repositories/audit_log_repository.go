package repositories

import (
	"gorm.io/gorm"

	"github.com/reshap0318/serv-uam/internal/models"
)

// AuditLogRepository extends GenericRepository — Create + FindAllWithOpts cover
// everything needed (insert-only, list-only), no bespoke queries required.
type AuditLogRepository struct {
	*GenericRepository[models.AuditLog]
}

// NewAuditLogRepository creates a new AuditLogRepository.
func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{
		GenericRepository: NewGenericRepository(db, &models.AuditLog{}),
	}
}
