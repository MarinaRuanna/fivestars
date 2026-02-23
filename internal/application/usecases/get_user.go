package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/get_user.go -package mock_usecases . GetUserUseCase
type GetUserUseCase interface {
	Execute(ctx context.Context, userID string) (*domain.User, error)
}

type getUserUseCase struct {
	userRepo domain.UserRepository
}

// NewGetUserUseCase creates a new GetUserUseCase.
func NewGetUserUseCase(userRepo domain.UserRepository) GetUserUseCase {
	return &getUserUseCase{userRepo: userRepo}
}

// GetUserInput DTO for the use case input.
type GetUserInput struct {
	UserID string
}

// GetUserOutput DTO for the use case output.
type GetUserOutput struct {
	ID        string
	Email     string
	Name      string
	Level     int
	AvatarURL string
}

func (uc *getUserUseCase) Execute(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if user == nil {
		return nil, customerror.NewNotFoundError("user not found")
	}

	err = user.Validate()
	if err != nil {
		return nil, err
	}

	return user, nil
}
