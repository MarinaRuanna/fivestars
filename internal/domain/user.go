package domain

import (
	"context"
	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
	"time"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/user_repository.go -package mock_domain . UserRepository
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, userID string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

// User represents a platform user.
type User struct {
	ID           string    `json:"user_id" validate:"required,uuid4"`
	Email        string    `json:"email" validate:"required,email"`
	PasswordHash string    `json:"-"` // nunca expor em JSON
	Name         string    `json:"name" validate:"required,min=1"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	Level        int       `json:"level"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if err := validator.Validate(u); err != nil {
		return customerror.NewValidationError(err.Error())
	}
	return nil
}
