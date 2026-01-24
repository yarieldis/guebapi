package models

// User represents a user in our system
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
