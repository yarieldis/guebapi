package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yarieldis/guebapi/internal/api/middleware"
	"github.com/yarieldis/guebapi/internal/repository/user"
	"github.com/yarieldis/guebapi/internal/service/auth"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter() (*gin.Engine, auth.Service, user.Repository) {
	repo := user.NewMemoryRepositoryWithData(map[string]string{
		"john": "password123",
	})
	authSvc := auth.NewJWTService(repo, "test-secret", 72*time.Hour)

	authHandler := NewAuthHandler(authSvc)
	profileHandler := NewProfileHandler(repo)

	router := gin.New()

	public := router.Group("/api")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
	}

	protected := router.Group("/api/protected")
	protected.Use(middleware.Auth(authSvc))
	{
		protected.GET("/profile", profileHandler.Profile)
		protected.POST("/update", profileHandler.Update)
	}

	return router, authSvc, repo
}

func TestAuthHandler_Login_Success(t *testing.T) {
	router, _, _ := setupTestRouter()

	body := map[string]string{"username": "john", "password": "password123"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["token"] == "" {
		t.Error("expected token in response")
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	router, _, _ := setupTestRouter()

	body := map[string]string{"username": "john", "password": "wrongpassword"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthHandler_Register_Success(t *testing.T) {
	router, _, _ := setupTestRouter()

	body := map[string]string{"username": "newuser", "password": "newpass"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthHandler_Register_DuplicateUsername(t *testing.T) {
	router, _, _ := setupTestRouter()

	body := map[string]string{"username": "john", "password": "anypass"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestProfileHandler_Profile_Success(t *testing.T) {
	router, authSvc, _ := setupTestRouter()

	token, _ := authSvc.GenerateToken("john")

	req := httptest.NewRequest(http.MethodGet, "/api/protected/profile", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["username"] != "john" {
		t.Errorf("expected username 'john', got '%s'", response["username"])
	}
}

func TestProfileHandler_Update_Success(t *testing.T) {
	router, authSvc, _ := setupTestRouter()

	token, _ := authSvc.GenerateToken("john")

	body := map[string]string{"password": "newpassword"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/protected/update", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
