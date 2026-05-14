package usecases_test

import (
	"context"
	"errors"
	"testing"

	"fivestars/internal/application/usecases"
	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/domain/mock_domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_CreateEstablishmentUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should create establishment successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewCreateEstablishmentUseCase(establishmentRepo)

		lat := -23.55052
		lng := -46.633308
		input := domain.Establishment{
			Name:     "Cafe Central",
			Slug:     "cafe-central",
			Category: "cafe",
			Address:  "Av. Paulista, 1000",
			Lat:      &lat,
			Lng:      &lng,
		}
		establishmentRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		result, err := uc.Execute(ctx, "11111111-1111-4111-8111-111111111111", input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "11111111-1111-4111-8111-111111111111", result.OwnerID)
		assert.Equal(t, "Cafe Central", result.Name)
	})

	t.Run("should return unauthorized when userID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewCreateEstablishmentUseCase(establishmentRepo)

		result, err := uc.Execute(ctx, "", domain.Establishment{Name: "Cafe", Slug: "cafe", Category: "cafe"})

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.UnauthorizedErrorType)
	})

	t.Run("should return validation error when input is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewCreateEstablishmentUseCase(establishmentRepo)

		result, err := uc.Execute(ctx, "11111111-1111-4111-8111-111111111111", domain.Establishment{})

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return wrapped error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewCreateEstablishmentUseCase(establishmentRepo)

		input := domain.Establishment{Name: "Cafe Central", Slug: "cafe-central", Category: "cafe"}
		repoErr := errors.New("insert failed")
		establishmentRepo.EXPECT().Create(ctx, gomock.Any()).Return(repoErr)

		result, err := uc.Execute(ctx, "11111111-1111-4111-8111-111111111111", input)

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to create establishment")
		assert.ErrorIs(t, err, repoErr)
	})
}
