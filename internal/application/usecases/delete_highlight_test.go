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

func Test_DeleteHighlightUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete highlight successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		uc := usecases.NewDeleteHighlightUseCase(highlightRepo, allowOperatorPolicy{allowed: true})

		review := domain_fakes.NewReviewBuilder().Build()
		highlightRepo.EXPECT().Delete(ctx, review.EstablishmentID, review.ID).Return(nil)

		err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		require.NoError(t, err)
	})

	t.Run("should return unauthorized when user cannot manage establishment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		uc := usecases.NewDeleteHighlightUseCase(highlightRepo, allowOperatorPolicy{allowed: false})

		review := domain_fakes.NewReviewBuilder().Build()

		err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		requireCustomErrorType(t, err, customerror.ForbiddenErrorType)
	})

	t.Run("should return validation error when establishmentID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		uc := usecases.NewDeleteHighlightUseCase(highlightRepo, allowOperatorPolicy{allowed: true})

		err := uc.Execute(ctx, "11111111-1111-4111-8111-111111111111", "", "44444444-4444-4444-8444-444444444444")

		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return wrapped error when authorization fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		policyErr := errors.New("policy failed")
		uc := usecases.NewDeleteHighlightUseCase(highlightRepo, allowOperatorPolicy{err: policyErr})

		review := domain_fakes.NewReviewBuilder().Build()

		err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to authorize establishment operator")
		assert.ErrorIs(t, err, policyErr)
	})

	t.Run("should return wrapped error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		uc := usecases.NewDeleteHighlightUseCase(highlightRepo, allowOperatorPolicy{allowed: true})

		review := domain_fakes.NewReviewBuilder().Build()
		repoErr := errors.New("delete failed")
		highlightRepo.EXPECT().Delete(ctx, review.EstablishmentID, review.ID).Return(repoErr)

		err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to delete highlight")
		assert.ErrorIs(t, err, repoErr)
	})
}
