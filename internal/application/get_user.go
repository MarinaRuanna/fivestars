package application

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
)

// GetUserUseCase implements user retrieval business logic.
// ⭐ Completely isolated from HTTP and database specifics.
type GetUserUseCase struct {
	userRepo domain.UserRepository
}

// NewGetUserUseCase creates a new GetUserUseCase.
func NewGetUserUseCase(userRepo domain.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{userRepo: userRepo}
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

// Execute runs the user retrieval logic.
func (uc *GetUserUseCase) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	// 1. VALIDATE
	if input.UserID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// 2. FETCH USER
	user, err := uc.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// 3. MAP TO OUTPUT (domain → output DTO)
	return &GetUserOutput{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Level:     user.Level,
		AvatarURL: user.AvatarURL,
	}, nil
}
