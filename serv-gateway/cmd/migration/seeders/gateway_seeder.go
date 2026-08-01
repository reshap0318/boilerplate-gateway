package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/serv-gateway/internal/models"
	"gorm.io/gorm"
)

// SeedUamService registers serv-uam as an upstream Service under base path /uam. Adjust
// BaseURL to wherever serv-uam actually runs in your environment.
func SeedUamService(db *gorm.DB) uint {
	fmt.Println("Seeding gateway service (serv-uam)...")

	const name = "serv-uam"

	var existing models.GatewayService
	err := db.Where("name = ?", name).First(&existing).Error
	if err == nil {
		fmt.Printf("  ⊘ Service %s already exists, skipping\n", name)
		return existing.ID
	}

	s := models.GatewayService{
		Name:         name,
		BaseURL:      "http://localhost:8081",
		BasePath:     "/uam",
		Protocol:     "http",
		IsActive:     true,
		HealthStatus: "unknown",
	}

	if err := db.Create(&s).Error; err != nil {
		log.Printf("Failed to create service %s: %v", name, err)
		return 0
	}

	fmt.Printf("✓ Seeded gateway service: %s (ID %d)\n", name, s.ID)
	return s.ID
}

// SeedUamRoutes registers serv-uam's User/Role/Permission CRUD routes, each guarded by the
// same permission names already used for the equivalent capability elsewhere (see
// permission_seeder.go) — moving where something is implemented doesn't rename it.
func SeedUamRoutes(db *gorm.DB, serviceID uint) {
	if serviceID == 0 {
		log.Println("serv-uam Service ID is 0, skipping route seeding")
		return
	}

	fmt.Println("Seeding gateway routes (serv-uam)...")

	type routeSeed struct {
		Method      string
		PathPattern string
		MatchMode   string
		Permissions []string
		Public      bool
	}

	// path_pattern is relative to the Service's base_path ("/uam", see SeedUamService), so
	// these resolve to e.g. "/uam/users", "/uam/users/:id".
	routes := []routeSeed{
		// "My own profile" — no specific permission, just needs to be authenticated (empty
		// Permissions, Public: false), same pattern as serv-message's notification routes.
		{"GET", "/me", "any", nil, false},
		{"PUT", "/me", "any", nil, false},

		{"POST", "/users", "any", []string{"user.create"}, false},
		{"GET", "/users", "any", []string{"user.index"}, false},
		{"GET", "/users/:id", "any", []string{"user.index"}, false},
		{"PUT", "/users/:id", "any", []string{"user.edit"}, false},
		{"PUT", "/users/:id/status", "any", []string{"user.edit"}, false},
		{"POST", "/users/:id/unlock", "any", []string{"user.edit"}, false},
		{"DELETE", "/users/:id", "any", []string{"user.delete"}, false},

		{"POST", "/roles", "any", []string{"role.create"}, false},
		{"GET", "/roles", "any", []string{"role.index"}, false},
		{"GET", "/roles/:id", "any", []string{"role.index"}, false},
		{"GET", "/roles/:id/permissions", "any", []string{"role.index"}, false},
		{"PUT", "/roles/:id", "any", []string{"role.edit"}, false},
		{"DELETE", "/roles/:id", "any", []string{"role.delete"}, false},

		{"POST", "/permissions", "any", []string{"permission.create"}, false},
		{"GET", "/permissions", "any", []string{"permission.index"}, false},
		{"GET", "/permissions/:id", "any", []string{"permission.index"}, false},
		{"PUT", "/permissions/:id", "any", []string{"permission.edit"}, false},
		{"DELETE", "/permissions/:id", "any", []string{"permission.delete"}, false},

		{"GET", "/audit-logs", "any", []string{"audit.index"}, false},

		// Forgot/reset-password are NOT proxied here — api-gateway forwards them to serv-uam
		// server-to-server as fixed Management API routes (see
		// serv-gateway/internal/services/auth_service.go), same as login/refresh.
	}

	count := 0
	for _, r := range routes {
		var existing models.GatewayRoute
		err := db.Where("service_id = ? AND method = ? AND path_pattern = ?", serviceID, r.Method, r.PathPattern).First(&existing).Error
		if err == nil {
			fmt.Printf("  ⊘ Route %s %s already exists, skipping\n", r.Method, r.PathPattern)
			continue
		}

		route := models.GatewayRoute{
			ServiceID:           serviceID,
			Method:              r.Method,
			PathPattern:         r.PathPattern,
			PermissionMatchMode: r.MatchMode,
			Permissions:         r.Permissions,
			Public:              r.Public,
			IsActive:            true,
		}

		if err := db.Create(&route).Error; err != nil {
			log.Printf("Failed to create route %s %s: %v", r.Method, r.PathPattern, err)
			continue
		}

		count++
	}

	fmt.Printf("✓ Seeded %d gateway routes\n", count)
}

// SeedMessageService registers serv-message as an upstream Service under base path
// /message. Adjust BaseURL to wherever serv-message actually runs in your environment.
func SeedMessageService(db *gorm.DB) uint {
	fmt.Println("Seeding gateway service (serv-message)...")

	const name = "serv-message"

	var existing models.GatewayService
	err := db.Where("name = ?", name).First(&existing).Error
	if err == nil {
		fmt.Printf("  ⊘ Service %s already exists, skipping\n", name)
		return existing.ID
	}

	s := models.GatewayService{
		Name:         name,
		BaseURL:      "http://localhost:8082",
		BasePath:     "/message",
		Protocol:     "http",
		IsActive:     true,
		HealthStatus: "unknown",
	}

	if err := db.Create(&s).Error; err != nil {
		log.Printf("Failed to create service %s: %v", name, err)
		return 0
	}

	fmt.Printf("✓ Seeded gateway service: %s (ID %d)\n", name, s.ID)
	return s.ID
}

// SeedMessageRoutes registers serv-message's notification routes, except POST (create) —
// that one is called service-to-service (e.g. from serv-uam) rather than through the
// gateway on behalf of an end user, so it's deliberately left unregistered here. The rest
// act on "the caller's own notifications" and need no specific permission — just an
// authenticated caller (empty Permissions, Public: false; see proxy/handler.go).
func SeedMessageRoutes(db *gorm.DB, serviceID uint) {
	if serviceID == 0 {
		log.Println("serv-message Service ID is 0, skipping route seeding")
		return
	}

	fmt.Println("Seeding gateway routes (serv-message)...")

	type routeSeed struct {
		Method      string
		PathPattern string
		MatchMode   string
		Permissions []string
		Public      bool
	}

	// path_pattern is relative to the Service's base_path ("/message", see
	// SeedMessageService), so these resolve to e.g. "/message/notifications".
	routes := []routeSeed{
		{"GET", "/notifications", "any", nil, false},
		{"GET", "/notifications/unread-count", "any", nil, false},
		{"PATCH", "/notifications/read-all", "any", nil, false},
		{"PATCH", "/notifications/:id/read", "any", nil, false},
		{"DELETE", "/notifications", "any", nil, false},
		{"DELETE", "/notifications/:id", "any", nil, false},
	}

	count := 0
	for _, r := range routes {
		var existing models.GatewayRoute
		err := db.Where("service_id = ? AND method = ? AND path_pattern = ?", serviceID, r.Method, r.PathPattern).First(&existing).Error
		if err == nil {
			fmt.Printf("  ⊘ Route %s %s already exists, skipping\n", r.Method, r.PathPattern)
			continue
		}

		route := models.GatewayRoute{
			ServiceID:           serviceID,
			Method:              r.Method,
			PathPattern:         r.PathPattern,
			PermissionMatchMode: r.MatchMode,
			Permissions:         r.Permissions,
			Public:              r.Public,
			IsActive:            true,
		}

		if err := db.Create(&route).Error; err != nil {
			log.Printf("Failed to create route %s %s: %v", r.Method, r.PathPattern, err)
			continue
		}

		count++
	}

	fmt.Printf("✓ Seeded %d gateway routes\n", count)
}
