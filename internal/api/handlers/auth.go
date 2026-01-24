package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yarieldis/guebapi/internal/models"
	"github.com/yarieldis/guebapi/internal/repository/user"
	"github.com/yarieldis/guebapi/internal/service/auth"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService auth.Service
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login requests
func (h *AuthHandler) Login(c *gin.Context) {
	var u models.User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Register handles user registration requests
func (h *AuthHandler) Register(c *gin.Context) {
	var u models.User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.authService.Register(c.Request.Context(), &u)
	if err != nil {
		if err == user.ErrUserAlreadyExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user created successfully"})
}
