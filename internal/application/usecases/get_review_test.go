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

func Test_GetReviewUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should return review when found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewGetReviewUseCase(reviewRepo)

		review := domain_fakes.NewReviewBuilder().Build()
		reviewRepo.EXPECT().GetByID(ctx, review.ID).Return(&review, nil)

		result, err := uc.Execute(ctx, review.ID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, review.ID, result.ID)
	})

	t.Run("should return validation error when reviewID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewGetReviewUseCase(reviewRepo)

		result, err := uc.Execute(ctx, "")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return not found when review does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewGetReviewUseCase(reviewRepo)

		reviewRepo.EXPECT().GetByID(ctx, "missing").Return(nil, nil)

		result, err := uc.Execute(ctx, "missing")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.NotFoundErrorType)
	})

	t.Run("should return wrapped error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewGetReviewUseCase(reviewRepo)

		repoErr := errors.New("db failed")
		reviewRepo.EXPECT().GetByID(ctx, "review").Return(nil, repoErr)

		result, err := uc.Execute(ctx, "review")

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to fetch review")
		assert.ErrorIs(t, err, repoErr)
	})
}
