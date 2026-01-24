package auth

import (
	"context"
	"testing"
	"time"

	"github.com/yarieldis/guebapi/internal/models"
	"github.com/yarieldis/guebapi/internal/repository/user"
)

func TestJWTService_Login(t *testing.T) {
	repo := user.NewMemoryRepositoryWithData(map[string]string{
		"john": "password123",
	})
	svc := NewJWTService(repo, "test-secret", 72*time.Hour)
	ctx := context.Background()

	// Test successful login
	token, err := svc.Login(ctx, "john", "password123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected token, got empty string")
	}

	// Test wrong password
	_, err = svc.Login(ctx, "john", "wrongpassword")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}

	// Test non-existing user
	_, err = svc.Login(ctx, "nonexistent", "password")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestJWTService_Register(t *testing.T) {
	repo := user.NewMemoryRepository()
	svc := NewJWTService(repo, "test-secret", 72*time.Hour)
	ctx := context.Background()

	// Test successful registration
	err := svc.Register(ctx, &models.User{Username: "newuser", Password: "newpass"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify user can login
	_, err = svc.Login(ctx, "newuser", "newpass")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test duplicate registration
	err = svc.Register(ctx, &models.User{Username: "newuser", Password: "anotherpass"})
	if err != user.ErrUserAlreadyExists {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestJWTService_ValidateToken(t *testing.T) {
	repo := user.NewMemoryRepository()
	svc := NewJWTService(repo, "test-secret", 72*time.Hour)

	// Generate a token
	token, err := svc.GenerateToken("john")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Test valid token
	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if claims.Username != "john" {
		t.Errorf("expected username 'john', got '%s'", claims.Username)
	}

	// Test invalid token
	_, err = svc.ValidateToken("invalid-token")
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestJWTService_GenerateToken(t *testing.T) {
	repo := user.NewMemoryRepository()
	svc := NewJWTService(repo, "test-secret", 72*time.Hour)

	token, err := svc.GenerateToken("testuser")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}

	// Verify the token is valid
	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Errorf("generated token should be valid: %v", err)
	}
	if claims.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", claims.Username)
	}
}
