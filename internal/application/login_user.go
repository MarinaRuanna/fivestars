package application

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/auth"
	"fivestars/internal/infra/config"
)

// LoginUserUseCase implements user login business logic.
// ⭐ Completely isolated from HTTP and database specifics.
type LoginUserUseCase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

// NewLoginUserUseCase creates a new LoginUserUseCase.
func NewLoginUserUseCase(userRepo domain.UserRepository, jwtSecret config.JWTConfig) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret.Secret,
	}
}

// LoginUserInput DTO for the use case input.
type LoginUserInput struct {
	Email    string
	Password string
}

// LoginUserOutput DTO for the use case output.
type LoginUserOutput struct {
	UserID string
	Token  string
}

// Execute runs the login logic.
func (uc *LoginUserUseCase) Execute(ctx context.Context, input LoginUserInput) (*LoginUserOutput, error) {
	// 1. VALIDATE
	if input.Email == "" || input.Password == "" {
		return nil, customerror.NewValidationError("email and password are required")
	}

	// 2. FETCH USER
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if user == nil {
		return nil, customerror.NewUnauthorizedError("invalid email or password")
	}

	// 3. VERIFY PASSWORD (auth utility)
	if !auth.CheckPassword(user.PasswordHash, input.Password) {
		return nil, customerror.NewUnauthorizedError("invalid email or password")
	}

	// 4. GENERATE TOKEN (auth utility)
	token, err := auth.NewToken(user.ID, uc.jwtSecret, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 5. RETURN OUTPUT (formatted for HTTP, but decoupled from HTTP layer)
	return &LoginUserOutput{
		UserID: user.ID,
		Token:  token,
	}, nil
}
