package user

import (
	"context"
	"errors"

	"github.com/yarieldis/guebapi/internal/models"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// Repository defines the interface for user storage operations
type Repository interface {
	Create(ctx context.Context, user *models.User) error
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	UpdatePassword(ctx context.Context, username, password string) error
	Exists(ctx context.Context, username string) (bool, error)
}
