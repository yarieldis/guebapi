package main

import (
	"log"

	"github.com/yarieldis/guebapi/config"
	"github.com/yarieldis/guebapi/internal/app"
	"github.com/yarieldis/guebapi/internal/repository/user"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Local development uses the in-memory repository with seeded users
	userRepo := user.NewMemoryRepositoryWithData(map[string]string{
		"john": "password123",
		"jane": "password456",
	})

	// Build the router
	r := app.BuildRouter(cfg, userRepo)

	// Start server
	addr := cfg.Server.Address()
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
