package app

import (
	"github.com/gin-gonic/gin"
	"github.com/yarieldis/guebapi/config"
	"github.com/yarieldis/guebapi/internal/api/handlers"
	"github.com/yarieldis/guebapi/internal/api/router"
	"github.com/yarieldis/guebapi/internal/repository/user"
	"github.com/yarieldis/guebapi/internal/service/auth"
)

// BuildRouter wires all dependencies together and returns the configured Gin engine.
func BuildRouter(cfg *config.Config, userRepo user.Repository) *gin.Engine {
	authService := auth.NewJWTService(userRepo, cfg.JWT.SecretKey, cfg.JWT.TokenDuration)
	authHandler := handlers.NewAuthHandler(authService)
	profileHandler := handlers.NewProfileHandler(userRepo)
	return router.SetupRouter(authHandler, profileHandler, authService)
}
