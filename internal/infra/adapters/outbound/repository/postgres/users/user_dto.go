package users

import (
	"fmt"
	"time"

	"fivestars/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"
)

// UserDTO representa uma linha da tabela users.
type UserDTO struct {
	ID           pgtype.UUID `json:"id_user" validate:"required,uuid4"`
	Email        string      `json:"email" validate:"required,email"`
	PasswordHash string      `json:"-"` // nunca expor em JSON
	Name         string      `json:"name" validate:"required"`
	AvatarURL    *string     `json:"avatar_url,omitempty"`
	Level        int         `json:"level" validate:"required"`
	CreatedAt    time.Time   `json:"created_at" validate:"required"`
	UpdatedAt    time.Time   `json:"updated_at" validate:"required"`
}

// ToDomain converte o DTO de persistência em entidade de domínio.
func (dto *UserDTO) ToDomain() (*domain.User, error) {
	user := &domain.User{
		ID:           uuidToString(dto.ID),
		Email:        dto.Email,
		PasswordHash: dto.PasswordHash,
		Name:         dto.Name,
		Level:        dto.Level,
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
	}
	if dto.AvatarURL != nil {
		user.AvatarURL = *dto.AvatarURL
	}

	err := user.Validate()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	b := u.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
