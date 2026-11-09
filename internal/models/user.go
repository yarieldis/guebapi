package models

// User represents a user in our system
type User struct {
	Username string `json:"username" dynamodbav:"username"`
	Password string `json:"password" dynamodbav:"password"`
}
