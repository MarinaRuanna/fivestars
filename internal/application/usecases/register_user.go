package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/auth"
	"fivestars/internal/infra/config"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/register_user.go -package mock_usecases . RegisterUserUseCase
type RegisterUserUseCase interface {
	Execute(ctx context.Context, inputLogin domain.UserRegistration) (*domain.AuthenticationResult, error)
}

type registerUserUseCase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewRegisterUserUseCase(userRepo domain.UserRepository, jwtSecret config.JWTConfig) RegisterUserUseCase {
	return &registerUserUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret.Secret,
	}
}

func (uc *registerUserUseCase) Execute(ctx context.Context, inputLogin domain.UserRegistration) (*domain.AuthenticationResult, error) {
	error := inputLogin.Validate()
	if error != nil {
		return nil, error
	}

	existingUser, err := uc.userRepo.GetByEmail(ctx, inputLogin.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, customerror.NewConflictError("user with this email already exists")
	}

	hashedPassword, err := auth.HashPassword(inputLogin.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := domain.NewUser(inputLogin.Email, hashedPassword, inputLogin.Name)
	if err != nil {
		return nil, err
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	createdUser, err := uc.userRepo.GetByEmail(ctx, inputLogin.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created user: %w", err)
	}

	token, err := auth.NewToken(createdUser.ID, uc.jwtSecret, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return domain.NewAuthenticationResult(createdUser.ID, createdUser.Name, token)
}
