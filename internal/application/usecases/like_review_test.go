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

func Test_LikeReviewUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should like review successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		likeRepo := mock_domain.NewMockReviewLikeRepository(ctrl)
		uc := usecases.NewLikeReviewUseCase(reviewRepo, likeRepo)

		review := domain_fakes.NewReviewBuilder().Build()
		reviewRepo.EXPECT().GetByID(ctx, review.ID).Return(&review, nil)
		likeRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		err := uc.Execute(ctx, review.UserID, review.ID)

		require.NoError(t, err)
	})

	t.Run("should return unauthorized when userID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		likeRepo := mock_domain.NewMockReviewLikeRepository(ctrl)
		uc := usecases.NewLikeReviewUseCase(reviewRepo, likeRepo)

		err := uc.Execute(ctx, "", "review")

		requireCustomErrorType(t, err, customerror.UnauthorizedErrorType)
	})

	t.Run("should return validation error when reviewID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		likeRepo := mock_domain.NewMockReviewLikeRepository(ctrl)
		uc := usecases.NewLikeReviewUseCase(reviewRepo, likeRepo)

		err := uc.Execute(ctx, "user", "")

		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return not found when review does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		likeRepo := mock_domain.NewMockReviewLikeRepository(ctrl)
		uc := usecases.NewLikeReviewUseCase(reviewRepo, likeRepo)

		reviewRepo.EXPECT().GetByID(ctx, "review").Return(nil, nil)

		err := uc.Execute(ctx, "user", "review")

		requireCustomErrorType(t, err, customerror.NotFoundErrorType)
	})

	t.Run("should return conflict when review already liked", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		likeRepo := mock_domain.NewMockReviewLikeRepository(ctrl)
		uc := usecases.NewLikeReviewUseCase(reviewRepo, likeRepo)

		review := domain_fakes.NewReviewBuilder().Build()
		reviewRepo.EXPECT().GetByID(ctx, review.ID).Return(&review, nil)
		likeRepo.EXPECT().Create(ctx, gomock.Any()).Return(customerror.NewConflictError("already liked"))

		err := uc.Execute(ctx, review.UserID, review.ID)

		requireCustomErrorType(t, err, customerror.ConflictErrorType)
	})

	t.Run("should return wrapped error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		likeRepo := mock_domain.NewMockReviewLikeRepository(ctrl)
		uc := usecases.NewLikeReviewUseCase(reviewRepo, likeRepo)

		review := domain_fakes.NewReviewBuilder().Build()
		repoErr := errors.New("db failed")
		reviewRepo.EXPECT().GetByID(ctx, review.ID).Return(&review, nil)
		likeRepo.EXPECT().Create(ctx, gomock.Any()).Return(repoErr)

		err := uc.Execute(ctx, review.UserID, review.ID)

		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to like review")
		assert.ErrorIs(t, err, repoErr)
	})
}
