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

func Test_GetEstablishmentStatsUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should return establishment stats successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewGetEstablishmentStatsUseCase(establishmentRepo)

		stats := &domain.EstablishmentStats{
			AverageRating:          4.5,
			TotalReviews:           12,
			TotalLikes:             48,
			HighlightedReviewCount: 3,
		}
		establishmentRepo.EXPECT().GetStats(ctx, "22222222-2222-4222-8222-222222222222").Return(stats, nil)

		result, err := uc.Execute(ctx, "22222222-2222-4222-8222-222222222222")

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 4.5, result.AverageRating)
		assert.Equal(t, 12, result.TotalReviews)
	})

	t.Run("should return validation error when establishmentID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewGetEstablishmentStatsUseCase(establishmentRepo)

		result, err := uc.Execute(ctx, "")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return not found when stats are missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewGetEstablishmentStatsUseCase(establishmentRepo)

		establishmentRepo.EXPECT().GetStats(ctx, "missing").Return(nil, nil)

		result, err := uc.Execute(ctx, "missing")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.NotFoundErrorType)
	})

	t.Run("should return wrapped error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewGetEstablishmentStatsUseCase(establishmentRepo)

		repoErr := errors.New("stats failed")
		establishmentRepo.EXPECT().GetStats(ctx, "22222222-2222-4222-8222-222222222222").Return(nil, repoErr)

		result, err := uc.Execute(ctx, "22222222-2222-4222-8222-222222222222")

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to fetch establishment stats")
		assert.ErrorIs(t, err, repoErr)
	})
}
