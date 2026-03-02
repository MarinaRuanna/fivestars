package usecases_test

import (
	"context"
	"errors"
	"fivestars/internal/application/usecases"
	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/domain/domain_fakes"
	"fivestars/internal/domain/mock_domain"
	"fivestars/internal/infra/auth"
	"fivestars/internal/infra/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_LoginUserUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	jwtCfg := config.JWTConfig{Secret: "test-secret"}

	t.Run("should login successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		userRepo := mock_domain.NewMockUserRepository(ctrl)
		uc := usecases.NewLoginUserUseCase(userRepo, jwtCfg)

		hash, err := auth.HashPassword("123456")
		require.NoError(t, err)

		input := domain.UserCredentials{Email: "user@example.com", Password: "123456"}
		user := domain_fakes.NewUserBuilder().WithEmail(input.Email).WithPasswordHash(hash).WithID("7aa6fb0c-976c-4f60-8710-9a10295b4868").Build()

		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(&user, nil)

		result, err := uc.Execute(ctx, input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, user.ID, result.UserID)
		assert.NotEmpty(t, result.Token)
	})

	t.Run("should return validation error when input is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		userRepo := mock_domain.NewMockUserRepository(ctrl)
		uc := usecases.NewLoginUserUseCase(userRepo, jwtCfg)

		input := domain.UserCredentials{Email: "invalid", Password: "123"}

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return unauthorized when user does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		userRepo := mock_domain.NewMockUserRepository(ctrl)
		uc := usecases.NewLoginUserUseCase(userRepo, jwtCfg)

		input := domain.UserCredentials{Email: "missing@example.com", Password: "123456"}
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(nil, nil)

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.UnauthorizedErrorType)
	})

	t.Run("should return unauthorized when password does not match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		userRepo := mock_domain.NewMockUserRepository(ctrl)
		uc := usecases.NewLoginUserUseCase(userRepo, jwtCfg)

		hash, err := auth.HashPassword("correct-password")
		require.NoError(t, err)

		input := domain.UserCredentials{Email: "user@example.com", Password: "wrong-password"}
		user := domain_fakes.NewUserBuilder().WithEmail(input.Email).WithPasswordHash(hash).Build()
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(&user, nil)

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.UnauthorizedErrorType)
	})

	t.Run("should return wrapped error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		userRepo := mock_domain.NewMockUserRepository(ctrl)
		uc := usecases.NewLoginUserUseCase(userRepo, jwtCfg)

		input := domain.UserCredentials{Email: "user@example.com", Password: "123456"}
		repoErr := errors.New("db unavailable")
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(nil, repoErr)

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to fetch user")
		assert.ErrorIs(t, err, repoErr)
	})
}
