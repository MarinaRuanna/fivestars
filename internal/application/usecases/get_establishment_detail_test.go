package usecases_test

import (
	"context"
	"errors"
	"testing"

	"fivestars/internal/application/usecases"
	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/domain/domain_fakes"
	"fivestars/internal/domain/mock_domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_GetEstablishmentDetailUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should return establishment detail with highlights", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewGetEstablishmentDetailUseCase(establishmentRepo, highlightRepo, reviewRepo)

		establishment := domain.Establishment{
			ID:       "22222222-2222-4222-8222-222222222222",
			Name:     "Cafe Central",
			Category: "food",
		}
		review := domain_fakes.NewReviewBuilder().WithEstablishmentID(establishment.ID).Build()
		highlights := []domain.Highlight{
			{
				EstablishmentID: establishment.ID,
				ReviewID:        review.ID,
				CreatedByUserID: review.UserID,
				CreatedAt:       review.CreatedAt,
			},
		}

		establishmentRepo.EXPECT().GetByID(ctx, establishment.ID).Return(&establishment, nil)
		highlightRepo.EXPECT().ListByEstablishment(ctx, establishment.ID).Return(highlights, nil)
		reviewRepo.EXPECT().GetByID(ctx, review.ID).Return(&review, nil)

		result, err := uc.Execute(ctx, establishment.ID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, establishment.ID, result.Establishment.ID)
		require.Len(t, result.Highlights, 1)
		assert.Equal(t, review.ID, result.Highlights[0].ID)
	})

	t.Run("should return validation error when establishmentID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewGetEstablishmentDetailUseCase(establishmentRepo, highlightRepo, reviewRepo)

		result, err := uc.Execute(ctx, "")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return not found when establishment does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewGetEstablishmentDetailUseCase(establishmentRepo, highlightRepo, reviewRepo)

		establishmentRepo.EXPECT().GetByID(ctx, "missing").Return(nil, nil)

		result, err := uc.Execute(ctx, "missing")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.NotFoundErrorType)
	})

	t.Run("should return empty highlights when none are selected", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewGetEstablishmentDetailUseCase(establishmentRepo, highlightRepo, reviewRepo)

		establishment := domain.Establishment{
			ID:       "22222222-2222-4222-8222-222222222222",
			Name:     "Cafe Central",
			Category: "food",
		}

		establishmentRepo.EXPECT().GetByID(ctx, establishment.ID).Return(&establishment, nil)
		highlightRepo.EXPECT().ListByEstablishment(ctx, establishment.ID).Return(nil, nil)

		result, err := uc.Execute(ctx, establishment.ID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Highlights)
	})

	t.Run("should return wrapped error when highlight lookup fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewGetEstablishmentDetailUseCase(establishmentRepo, highlightRepo, reviewRepo)

		establishment := domain.Establishment{
			ID:       "22222222-2222-4222-8222-222222222222",
			Name:     "Cafe Central",
			Category: "food",
		}
		repoErr := errors.New("highlight list failed")

		establishmentRepo.EXPECT().GetByID(ctx, establishment.ID).Return(&establishment, nil)
		highlightRepo.EXPECT().ListByEstablishment(ctx, establishment.ID).Return(nil, repoErr)

		result, err := uc.Execute(ctx, establishment.ID)

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to list establishment highlights")
		assert.ErrorIs(t, err, repoErr)
	})
}
