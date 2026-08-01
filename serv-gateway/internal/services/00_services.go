package services

import (
	"github.com/reshap0318/serv-gateway/internal/database"
	"github.com/reshap0318/serv-gateway/internal/helpers"
	"github.com/reshap0318/serv-gateway/internal/repositories"
)

// ServicesConfig holds all dependencies for Services.
// Add new dependencies here without changing NewServices signature.
type ServicesConfig struct {
	Repo   *repositories.Repositories
	Redis  *database.RedisCache
	Logger *helpers.Logger
}

// Services holds all service dependencies.
type Services struct {
	repo        *repositories.Repositories
	RedisClient *database.RedisCache
	Logger      *helpers.Logger
	JWKSManager *JWKSManager
	Access      *helpers.Access
	RouteCache  RouteCacheRefresher
}

// NewServices creates and initializes all services.
func NewServices(cfg *ServicesConfig) *Services {
	return &Services{
		repo:        cfg.Repo,
		RedisClient: cfg.Redis,
		Logger:      cfg.Logger,
	}
}
