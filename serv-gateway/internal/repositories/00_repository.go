package repositories

import "gorm.io/gorm"

type Repositories struct {
	DB             *gorm.DB
	TxManager      *TransactionManager
	GatewayService *GatewayServiceRepository
	GatewayRoute   *GatewayRouteRepository
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
	txManager := NewTransactionManager(db)
	gatewayServiceRepo := NewGatewayServiceRepository(db)
	gatewayRouteRepo := NewGatewayRouteRepository(db)

	return &Repositories{
		DB:             db,
		TxManager:      txManager,
		GatewayService: gatewayServiceRepo,
		GatewayRoute:   gatewayRouteRepo,
	}, nil
}
