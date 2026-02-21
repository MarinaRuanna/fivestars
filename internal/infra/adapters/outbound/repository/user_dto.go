package repository

import (
	"fmt"
	"time"

	"fivestars/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"
)

// UserRow representa uma linha da tabela users.
type UserRow struct {
	ID           pgtype.UUID
	Email        string
	PasswordHash string
	Name         string
	AvatarURL    *string
	Level        int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ToDomain converte o DTO de persistência em entidade de domínio.
func (r *UserRow) ToDomain() *domain.User {
	u := &domain.User{
		ID:           uuidToString(r.ID),
		Email:        r.Email,
		PasswordHash: r.PasswordHash,
		Name:         r.Name,
		Level:        r.Level,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	if r.AvatarURL != nil {
		u.AvatarURL = *r.AvatarURL
	}
	return u
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	b := u.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
