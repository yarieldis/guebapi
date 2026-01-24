package config

import "time"

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	JWT    JWTConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host string
	Port string
}

// JWTConfig holds JWT-specific configuration
type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

// Address returns the server address in host:port format
func (c *ServerConfig) Address() string {
	return c.Host + ":" + c.Port
}
