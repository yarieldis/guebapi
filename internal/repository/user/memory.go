package user

import (
	"context"
	"sync"

	"github.com/yarieldis/guebapi/internal/models"
)

// MemoryRepository implements Repository with in-memory storage
type MemoryRepository struct {
	mu    sync.RWMutex
	users map[string]string // username -> password
}

// NewMemoryRepository creates a new in-memory user repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[string]string),
	}
}

// NewMemoryRepositoryWithData creates a new in-memory user repository with initial data
func NewMemoryRepositoryWithData(users map[string]string) *MemoryRepository {
	copied := make(map[string]string)
	for k, v := range users {
		copied[k] = v
	}
	return &MemoryRepository{
		users: copied,
	}
}

func (r *MemoryRepository) Create(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Username]; exists {
		return ErrUserAlreadyExists
	}

	r.users[user.Username] = user.Password
	return nil
}

func (r *MemoryRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	password, exists := r.users[username]
	if !exists {
		return nil, ErrUserNotFound
	}

	return &models.User{
		Username: username,
		Password: password,
	}, nil
}

func (r *MemoryRepository) UpdatePassword(ctx context.Context, username, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[username]; !exists {
		return ErrUserNotFound
	}

	r.users[username] = password
	return nil
}

func (r *MemoryRepository) Exists(ctx context.Context, username string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.users[username]
	return exists, nil
}
