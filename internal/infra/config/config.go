package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application configuration from environment.
type Config struct {
	DatabaseURL string
	Port        int
	JWTSecret   string
}

// Load reads config from environment variables.
func Load() *Config {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        port,
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
}

// Validate checks if required config values are set
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid PORT: must be between 1 and 65535")
	}
	return nil
}
