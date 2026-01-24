package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yarieldis/guebapi/internal/repository/user"
	"github.com/yarieldis/guebapi/internal/service/auth"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuth_MissingToken(t *testing.T) {
	repo := user.NewMemoryRepository()
	authSvc := auth.NewJWTService(repo, "test-secret", 72*time.Hour)

	router := gin.New()
	router.Use(Auth(authSvc))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	repo := user.NewMemoryRepository()
	authSvc := auth.NewJWTService(repo, "test-secret", 72*time.Hour)

	router := gin.New()
	router.Use(Auth(authSvc))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuth_ValidToken(t *testing.T) {
	repo := user.NewMemoryRepository()
	authSvc := auth.NewJWTService(repo, "test-secret", 72*time.Hour)

	token, err := authSvc.GenerateToken("testuser")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	router := gin.New()
	router.Use(Auth(authSvc))
	router.GET("/protected", func(c *gin.Context) {
		username, _ := c.Get("username")
		c.JSON(http.StatusOK, gin.H{"username": username})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
