package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yarieldis/guebapi/internal/api/handlers"
	"github.com/yarieldis/guebapi/internal/api/middleware"
	"github.com/yarieldis/guebapi/internal/service/auth"
)

// SetupRouter configures and returns the Gin router
func SetupRouter(authHandler *handlers.AuthHandler, profileHandler *handlers.ProfileHandler, authService auth.Service) *gin.Engine {
	router := gin.Default()

	// Public routes
	public := router.Group("/api")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
	}

	// Protected routes
	protected := router.Group("/api/protected")
	protected.Use(middleware.Auth(authService))
	{
		protected.GET("/profile", profileHandler.Profile)
		protected.POST("/update", profileHandler.Update)
	}

	return router
}
