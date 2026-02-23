package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/auth"
	"fivestars/internal/infra/config"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/login_user.go -package mock_usecases . LoginUserUseCase
type LoginUserUseCase interface {
	Execute(ctx context.Context, input domain.UserCredentials) (*domain.AuthenticationResult, error)
}

type loginUserUseCase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewLoginUserUseCase(userRepo domain.UserRepository, jwtSecret config.JWTConfig) LoginUserUseCase {
	return &loginUserUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret.Secret,
	}
}

func (uc *loginUserUseCase) Execute(ctx context.Context, inputLogin domain.UserCredentials) (*domain.AuthenticationResult, error) {
	err := inputLogin.Validate()
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetByEmail(ctx, inputLogin.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if user == nil {
		return nil, customerror.NewNotFoundError("user not found")
	}

	if !auth.CheckPassword(user.PasswordHash, inputLogin.Password) {
		return nil, customerror.NewUnauthorizedError("invalid email or password")
	}

	token, err := auth.NewToken(user.ID, uc.jwtSecret, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return domain.NewAuthenticationResult(user.ID, user.Name, token)
}
