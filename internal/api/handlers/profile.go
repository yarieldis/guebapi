package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yarieldis/guebapi/internal/models"
	"github.com/yarieldis/guebapi/internal/repository/user"
)

// ProfileHandler handles profile-related requests
type ProfileHandler struct {
	userRepo user.Repository
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(userRepo user.Repository) *ProfileHandler {
	return &ProfileHandler{
		userRepo: userRepo,
	}
}

// Profile returns the current user's profile
func (h *ProfileHandler) Profile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"message":  "Welcome to your profile!",
	})
}

// Update handles user profile update requests
func (h *ProfileHandler) Update(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var u models.User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.userRepo.UpdatePassword(c.Request.Context(), username.(string), u.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}
