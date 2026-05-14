package usecases_test

import (
	"context"
	"errors"
	"testing"

	"fivestars/internal/application/usecases"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/domain/domain_fakes"
	"fivestars/internal/domain/mock_domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_ClaimEstablishmentOwnershipUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should claim establishment successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewClaimEstablishmentOwnershipUseCase(establishmentRepo)

		establishment := domain_fakes.NewEstablishmentBuilder().WithOwnerID("11111111-1111-4111-8111-111111111111").Build()
		establishmentRepo.EXPECT().
			ClaimOwnership(ctx, establishment.ID, "11111111-1111-4111-8111-111111111111", "qr-cafe-central").
			Return(&establishment, nil)

		result, err := uc.Execute(ctx, "11111111-1111-4111-8111-111111111111", establishment.ID, "qr-cafe-central")

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "11111111-1111-4111-8111-111111111111", result.OwnerID)
	})

	t.Run("should return unauthorized when userID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewClaimEstablishmentOwnershipUseCase(establishmentRepo)

		result, err := uc.Execute(ctx, "", "22222222-2222-4222-8222-222222222222", "qr-cafe-central")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.UnauthorizedErrorType)
	})

	t.Run("should return validation error when establishmentID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewClaimEstablishmentOwnershipUseCase(establishmentRepo)

		result, err := uc.Execute(ctx, "11111111-1111-4111-8111-111111111111", "", "qr-cafe-central")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return validation error when qr_code is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewClaimEstablishmentOwnershipUseCase(establishmentRepo)

		result, err := uc.Execute(ctx, "11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222", "")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return not found when establishment does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewClaimEstablishmentOwnershipUseCase(establishmentRepo)

		establishmentRepo.EXPECT().
			ClaimOwnership(ctx, "22222222-2222-4222-8222-222222222222", "11111111-1111-4111-8111-111111111111", "qr-cafe-central").
			Return(nil, nil)

		result, err := uc.Execute(ctx, "11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222", "qr-cafe-central")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.NotFoundErrorType)
	})

	t.Run("should return wrapped error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewClaimEstablishmentOwnershipUseCase(establishmentRepo)

		repoErr := errors.New("claim failed")
		establishmentRepo.EXPECT().
			ClaimOwnership(ctx, "22222222-2222-4222-8222-222222222222", "11111111-1111-4111-8111-111111111111", "qr-cafe-central").
			Return(nil, repoErr)

		result, err := uc.Execute(ctx, "11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222", "qr-cafe-central")

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to claim establishment ownership")
		assert.ErrorIs(t, err, repoErr)
	})
}
