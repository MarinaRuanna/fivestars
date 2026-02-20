package config

import (
	"os"
	"strconv"
)

// Config holds application configuration from environment.
type Config struct {
	DatabaseURL string
	Port        int
	JWTSecret   string // usado na Fase 2
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
