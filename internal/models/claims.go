package models

import "github.com/golang-jwt/jwt/v4"

// Claims represents JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
