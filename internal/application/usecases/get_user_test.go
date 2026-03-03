package usecases_test

import (
	"context"
	"errors"
	"fivestars/internal/application/usecases"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/domain/domain_fakes"
	"fivestars/internal/domain/mock_domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_GetUserUseCase_Execute(t *testing.T) {
	userID := "user-id"

	t.Run("should return user successfully", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		ctx := context.Background()
		expectedUser := domain_fakes.NewUserBuilder().WithID(userID).Build()
		userRepo := mock_domain.NewMockUserRepository(ctrl)

		userRepo.EXPECT().GetByID(ctx, userID).Return(&expectedUser, nil)

		useCase := usecases.NewGetUserUseCase(userRepo)

		// When
		user, err := useCase.Execute(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, &expectedUser, user)
	})

	t.Run("should return error if user not found", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		ctx := context.Background()
		userID := "user-id"
		userRepo := mock_domain.NewMockUserRepository(ctrl)

		userRepo.EXPECT().GetByID(ctx, userID).Return(nil, nil)

		useCase := usecases.NewGetUserUseCase(userRepo)

		// When
		user, err := useCase.Execute(ctx, userID)

		// Then
		assert.Error(t, err)
		var customError *customerror.CustomError
		assert.ErrorAs(t, err, &customError)
		assert.Equal(t, customerror.NotFoundErrorType, customError.ErrorType())
		assert.Nil(t, user)
	})

	t.Run("should reurn unauthorized error when user id is empty", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		ctx := context.Background()
		userID := ""
		userRepo := mock_domain.NewMockUserRepository(ctrl)
		useCase := usecases.NewGetUserUseCase(userRepo)

		// When
		user, err := useCase.Execute(ctx, userID)

		// Then
		assert.Nil(t, user)
		assert.Error(t, err)
		var customError *customerror.CustomError
		assert.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.UnauthorizedErrorType, customError.ErrorType())
	})

	t.Run("should return an error when repository fails", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		ctx := context.Background()

		userRepo := mock_domain.NewMockUserRepository(ctrl)
		useCase := usecases.NewGetUserUseCase(userRepo)

		userRepo.EXPECT().GetByID(ctx, userID).Return(nil, assert.AnError)

		// When
		user, err := useCase.Execute(ctx, userID)

		// Then
		assert.Nil(t, user)
		assert.Error(t, err)
	})
}
