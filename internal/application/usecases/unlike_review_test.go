package usecases_test

import (
	"context"
	"errors"
	"testing"

	"fivestars/internal/application/usecases"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/domain/mock_domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_UnlikeReviewUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should unlike review successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		likeRepo := mock_domain.NewMockReviewLikeRepository(ctrl)
		uc := usecases.NewUnlikeReviewUseCase(likeRepo)

		likeRepo.EXPECT().Delete(ctx, "user", "review").Return(nil)

		err := uc.Execute(ctx, "user", "review")

		require.NoError(t, err)
	})

	t.Run("should return unauthorized when userID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		likeRepo := mock_domain.NewMockReviewLikeRepository(ctrl)
		uc := usecases.NewUnlikeReviewUseCase(likeRepo)

		err := uc.Execute(ctx, "", "review")

		requireCustomErrorType(t, err, customerror.UnauthorizedErrorType)
	})

	t.Run("should return validation error when reviewID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		likeRepo := mock_domain.NewMockReviewLikeRepository(ctrl)
		uc := usecases.NewUnlikeReviewUseCase(likeRepo)

		err := uc.Execute(ctx, "user", "")

		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return wrapped error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		likeRepo := mock_domain.NewMockReviewLikeRepository(ctrl)
		uc := usecases.NewUnlikeReviewUseCase(likeRepo)

		repoErr := errors.New("db failed")
		likeRepo.EXPECT().Delete(ctx, "user", "review").Return(repoErr)

		err := uc.Execute(ctx, "user", "review")

		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to unlike review")
		assert.ErrorIs(t, err, repoErr)
	})
}
