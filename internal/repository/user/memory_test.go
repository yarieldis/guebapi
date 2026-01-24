package user

import (
	"context"
	"testing"

	"github.com/yarieldis/guebapi/internal/models"
)

func TestMemoryRepository_Create(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	user := &models.User{Username: "testuser", Password: "testpass"}

	// Test successful creation
	err := repo.Create(ctx, user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test duplicate creation
	err = repo.Create(ctx, user)
	if err != ErrUserAlreadyExists {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestMemoryRepository_GetByUsername(t *testing.T) {
	repo := NewMemoryRepositoryWithData(map[string]string{
		"john": "password123",
	})
	ctx := context.Background()

	// Test existing user
	user, err := repo.GetByUsername(ctx, "john")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if user.Username != "john" || user.Password != "password123" {
		t.Errorf("unexpected user data: %+v", user)
	}

	// Test non-existing user
	_, err = repo.GetByUsername(ctx, "nonexistent")
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestMemoryRepository_UpdatePassword(t *testing.T) {
	repo := NewMemoryRepositoryWithData(map[string]string{
		"john": "oldpassword",
	})
	ctx := context.Background()

	// Test successful update
	err := repo.UpdatePassword(ctx, "john", "newpassword")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify password was updated
	user, _ := repo.GetByUsername(ctx, "john")
	if user.Password != "newpassword" {
		t.Errorf("expected newpassword, got %s", user.Password)
	}

	// Test update non-existing user
	err = repo.UpdatePassword(ctx, "nonexistent", "pass")
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestMemoryRepository_Exists(t *testing.T) {
	repo := NewMemoryRepositoryWithData(map[string]string{
		"john": "password123",
	})
	ctx := context.Background()

	// Test existing user
	exists, err := repo.Exists(ctx, "john")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !exists {
		t.Error("expected user to exist")
	}

	// Test non-existing user
	exists, err = repo.Exists(ctx, "nonexistent")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if exists {
		t.Error("expected user not to exist")
	}
}
