package domain

import "time"

// User represents a platform user (Fase 1 mínima; expandir na Fase 2).
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // nunca expor em JSON
	Name         string    `json:"name"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	Level        int       `json:"level"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
