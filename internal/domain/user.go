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

type User struct {
	ID           string
	Email        string `validate:"required,email"`
	PasswordHash string
	Name         string `validate:"required,min=1"`
	AvatarURL    string `validate:"omitempty,url"`
	Level        int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(email, passwordHash, name string) (*User, error) {
	now := time.Now().UTC()
	user := &User{
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Validate() error {
	if err := validator.Validate(u); err != nil {
		return customerror.NewValidationError(err.Error())
	}
	return nil
}
