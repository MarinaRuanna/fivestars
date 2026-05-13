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

func Test_ListReviewsByEstablishmentUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should list reviews successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewListReviewsByEstablishmentUseCase(reviewRepo)

		review := domain_fakes.NewReviewBuilder().Build()
		opts := domain.ReviewListOptions{Limit: 20}
		reviewRepo.EXPECT().ListByEstablishment(ctx, review.EstablishmentID, opts).Return([]domain.Review{review}, nil)

		result, err := uc.Execute(ctx, domain.ReviewListOptions{}, review.EstablishmentID)

		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, review.ID, result[0].ID)
	})

	t.Run("should return validation error when establishmentID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewListReviewsByEstablishmentUseCase(reviewRepo)

		result, err := uc.Execute(ctx, domain.ReviewListOptions{}, "")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return empty list when repository returns nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewListReviewsByEstablishmentUseCase(reviewRepo)

		opts := domain.ReviewListOptions{Limit: 20}
		reviewRepo.EXPECT().ListByEstablishment(ctx, "estab", opts).Return(nil, nil)

		result, err := uc.Execute(ctx, domain.ReviewListOptions{}, "estab")

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewListReviewsByEstablishmentUseCase(reviewRepo)

		repoErr := errors.New("db failed")
		opts := domain.ReviewListOptions{Limit: 20}
		reviewRepo.EXPECT().ListByEstablishment(ctx, "estab", opts).Return(nil, repoErr)

		result, err := uc.Execute(ctx, domain.ReviewListOptions{}, "estab")

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
	})
}
