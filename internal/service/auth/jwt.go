package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/yarieldis/guebapi/internal/models"
	"github.com/yarieldis/guebapi/internal/repository/user"
)

// JWTService implements Service using JWT tokens
type JWTService struct {
	userRepo      user.Repository
	secretKey     []byte
	tokenDuration time.Duration
}

// NewJWTService creates a new JWT authentication service
func NewJWTService(userRepo user.Repository, secretKey string, tokenDuration time.Duration) *JWTService {
	return &JWTService{
		userRepo:      userRepo,
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
	}
}

func (s *JWTService) Login(ctx context.Context, username, password string) (string, error) {
	storedUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if storedUser.Password != password {
		return "", ErrInvalidCredentials
	}

	return s.GenerateToken(username)
}

func (s *JWTService) Register(ctx context.Context, u *models.User) error {
	return s.userRepo.Create(ctx, u)
}

func (s *JWTService) ValidateToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return s.secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *JWTService) GenerateToken(username string) (string, error) {
	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}
