package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/reshap0318/serv-uam/internal/di"
	"github.com/reshap0318/serv-uam/internal/helpers"
	"github.com/reshap0318/serv-uam/internal/middleware"
	"github.com/reshap0318/serv-uam/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	host := helpers.GetEnv("APP_HOST", "0.0.0.0")
	port := helpers.GetEnv("APP_PORT", "8080")
	trustedProxies := helpers.GetEnv("TRUSTED_PROXIES", "")

	gin.SetMode(helpers.GetEnv("GIN_MODE", "release"))

	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	r := gin.Default()

	if trustedProxies != "" {
		if err := r.SetTrustedProxies(strings.Split(trustedProxies, ",")); err != nil {
			log.Printf("Warning: failed to set trusted proxies: %v", err)
		}
	}

	// TraceID runs first so every response — including auth failures — carries a trace id.
	r.Use(middleware.TraceID())

	r.NoRoute(func(c *gin.Context) {
		helpers.NotFound(c, "Endpoint not found")
	})

	routes.RegisterAll(r, container.Handlers)

	addr := host + ":" + port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
