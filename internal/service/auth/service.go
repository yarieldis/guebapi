package auth

import (
	"context"
	"errors"

	"github.com/yarieldis/guebapi/internal/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

// Service defines the interface for authentication operations
type Service interface {
	// Login validates credentials and returns a JWT token
	Login(ctx context.Context, username, password string) (string, error)

	// Register creates a new user
	Register(ctx context.Context, user *models.User) error

	// ValidateToken validates a JWT token and returns the claims
	ValidateToken(token string) (*models.Claims, error)

	// GenerateToken generates a JWT token for a username
	GenerateToken(username string) (string, error)
}
