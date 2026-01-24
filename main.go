package main

import (
	"log"

	"github.com/yarieldis/guebapi/config"
	"github.com/yarieldis/guebapi/internal/api/handlers"
	"github.com/yarieldis/guebapi/internal/api/router"
	"github.com/yarieldis/guebapi/internal/repository/user"
	"github.com/yarieldis/guebapi/internal/service/auth"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize repository with default users
	userRepo := user.NewMemoryRepositoryWithData(map[string]string{
		"john": "password123",
		"jane": "password456",
	})

	// Initialize services
	authService := auth.NewJWTService(userRepo, cfg.JWT.SecretKey, cfg.JWT.TokenDuration)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	profileHandler := handlers.NewProfileHandler(userRepo)

	// Setup router
	r := router.SetupRouter(authHandler, profileHandler, authService)

	// Start server
	addr := cfg.Server.Address()
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
