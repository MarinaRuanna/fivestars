package application

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/auth"
	"fivestars/internal/infra/config"
)

// RegisterUserUseCase implements user registration business logic.
// ⭐ Completely isolated from HTTP and database specifics.
type RegisterUserUseCase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

// NewRegisterUserUseCase creates a new RegisterUserUseCase.
func NewRegisterUserUseCase(userRepo domain.UserRepository, jwtSecret config.JWTConfig) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret.Secret,
	}
}

// RegisterUserInput DTO for the use case input.
type RegisterUserInput struct {
	Email    string
	Password string
	Name     string
}

// RegisterUserOutput DTO for the use case output.
type RegisterUserOutput struct {
	UserID string
	Token  string
	Name   string
}

// Execute runs the registration logic.
func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	// 1. VALIDATE
	if input.Email == "" || input.Password == "" || input.Name == "" {
		return nil, customerror.NewValidationError("email, password, and name are required")
	}

	if len(input.Password) < 6 {
		return nil, customerror.NewValidationError("password must be at least 6 characters")
	}

	// 2. CHECK DUPLICATE (domain-level constraint)
	existingUser, _ := uc.userRepo.GetByEmail(ctx, input.Email)
	if existingUser != nil {
		return nil, customerror.NewConflictError(fmt.Sprintf("user with email %s already exists", input.Email))
	}

	// 3. HASH PASSWORD (auth utility, not HTTP layer)
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 4. CREATE USER (domain entity)
	user := &domain.User{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Name:         input.Name,
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Get user ID for token generation
	createdUser, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created user: %w", err)
	}

	// 5. GENERATE TOKEN (auth utility)
	token, err := auth.NewToken(createdUser.ID, uc.jwtSecret, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 6. RETURN OUTPUT (formatted for HTTP, but decoupled from HTTP layer)
	return &RegisterUserOutput{
		UserID: createdUser.ID,
		Token:  token,
		Name:   createdUser.Name,
	}, nil
}
