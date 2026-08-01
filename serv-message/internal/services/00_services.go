package services

import (
	"github.com/reshap0318/serv-message/internal/database"
	"github.com/reshap0318/serv-message/internal/helpers"
	"github.com/reshap0318/serv-message/internal/pkg/email"
	"github.com/reshap0318/serv-message/internal/repositories"
)

// ServicesConfig holds all dependencies for Services.
// Add new dependencies here without changing NewServices signature.
type ServicesConfig struct {
	Repo   *repositories.Repositories
	Redis  *database.RedisCache
	Logger *helpers.Logger
	Access *helpers.Access
	Email  *email.EmailClient
}

// Services holds all service dependencies.
type Services struct {
	repo        *repositories.Repositories
	RedisClient *database.RedisCache
	Logger      *helpers.Logger
	// Access checks the caller's roles/permissions (from api-gateway headers).
	// No route-level enforcement anymore — call it directly wherever business
	// logic needs a check, e.g. `if !s.Access.HasPermission(ctx, "notification.delete") { ... }`.
	Access *helpers.Access
	Email  *email.EmailClient
}

// NewServices creates and initializes all services.
func NewServices(cfg *ServicesConfig) *Services {
	return &Services{
		repo:        cfg.Repo,
		RedisClient: cfg.Redis,
		Logger:      cfg.Logger,
		Access:      cfg.Access,
		Email:       cfg.Email,
	}
}
