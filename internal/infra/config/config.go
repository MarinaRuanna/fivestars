package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config holds application configuration from environment.
type Config struct {
	AppPort          int            `envconfig:"appPort" required:"true" default:"8080"`
	DatabasePostgres PostgresConfig `envconfig:"postgres" required:"true"`
	JWTSecret        JWTConfig      `envconfig:"jwt" required:"true"`
	CORS             CORSConfig     `envconfig:"cors"`
}

type PostgresConfig struct {
	Host     string `envconfig:"host" required:"true"`
	Port     int    `envconfig:"port" required:"true"`
	User     string `envconfig:"user" required:"true"`
	Password string `envconfig:"password" required:"true"`
	Database string `envconfig:"database" required:"true"`
	SSLMode  string `envconfig:"sslmode"`
	MaxConns int32  `envconfig:"maxConns"`
	MinConns int32  `envconfig:"minConns"`
}

type JWTConfig struct {
	Secret string `envconfig:"secret" required:"true"`
}

type CORSConfig struct {
	AllowedOrigins []string `envconfig:"allowed_origins" default:"http://localhost:3000,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:5173"`
}

func (pc PostgresConfig) DSN() string {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&pool_max_conns=%d&pool_min_conns=%d",
		pc.User,
		pc.Password,
		pc.Host,
		pc.Port,
		pc.Database,
		pc.SSLMode,
		pc.MaxConns,
		pc.MinConns)

	return dsn
}

// Load reads config from environment variables.
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}

// Validate checks if required config values are set
func (c *Config) Validate() error {
	if c.DatabasePostgres == (PostgresConfig{}) {
		return fmt.Errorf("DATABASE_POSTGRES config is required")
	}
	if c.JWTSecret == (JWTConfig{}) {
		return fmt.Errorf("JWT_SECRET config is required")
	}
	return nil
}
