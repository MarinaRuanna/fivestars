package controller

import (
	"fivestars/internal/domain"
	"time"
)

// UserResponse é o contrato da API para o perfil do usuário (GET /users/me).
// Não expõe password_hash.
type UserResponse struct {
	ID        string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Level     int    `json:"level"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UserFromDomain converte entidade de domínio em DTO de resposta.
func UserFromDomain(u *domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Level:     u.Level,
		CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
