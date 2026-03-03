package usecases_test

import (
	"context"
	"errors"
	"fivestars/internal/application/usecases"
	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/domain/domain_fakes"
	"fivestars/internal/domain/mock_domain"
	"fivestars/internal/infra/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_RegisterUserUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	jwtCfg := config.JWTConfig{Secret: "test-secret"}

	t.Run("should register user successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		userRepo := mock_domain.NewMockUserRepository(ctrl)
		uc := usecases.NewRegisterUserUseCase(userRepo, jwtCfg)

		input := domain.UserRegistration{Email: "new@example.com", Password: "123456", Name: "New User"}
		createdUser := domain_fakes.NewUserBuilder().WithEmail(input.Email).WithName(input.Name).WithID("a8f96f57-7b80-4796-8ae8-d5d2ffa2f0bf").Build()

		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(nil, nil)
		userRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, u *domain.User) error {
			require.NotEmpty(t, u.PasswordHash)
			assert.NotEqual(t, input.Password, u.PasswordHash)
			return nil
		})
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(&createdUser, nil)

		result, err := uc.Execute(ctx, input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, createdUser.ID, result.UserID)
		assert.NotEmpty(t, result.Token)
	})

	t.Run("should return validation error when input is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		userRepo := mock_domain.NewMockUserRepository(ctrl)
		uc := usecases.NewRegisterUserUseCase(userRepo, jwtCfg)

		input := domain.UserRegistration{Email: "invalid-email", Password: "123", Name: ""}

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return conflict when user already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		userRepo := mock_domain.NewMockUserRepository(ctrl)
		uc := usecases.NewRegisterUserUseCase(userRepo, jwtCfg)

		input := domain.UserRegistration{Email: "existing@example.com", Password: "123456", Name: "Existing User"}
		existing := domain_fakes.NewUserBuilder().WithEmail(input.Email).Build()

		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(&existing, nil)

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ConflictErrorType)
	})

	t.Run("should return wrapped error when checking existing user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		userRepo := mock_domain.NewMockUserRepository(ctrl)
		uc := usecases.NewRegisterUserUseCase(userRepo, jwtCfg)

		input := domain.UserRegistration{Email: "new@example.com", Password: "123456", Name: "New User"}
		repoErr := errors.New("db down")
		userRepo.EXPECT().GetByEmail(ctx, input.Email).Return(nil, repoErr)

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to check existing user")
		assert.ErrorIs(t, err, repoErr)
	})
}
