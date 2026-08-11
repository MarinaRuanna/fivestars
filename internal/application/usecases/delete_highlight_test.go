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

func Test_DeleteHighlightUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete highlight successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewDeleteHighlightUseCase(highlightRepo, establishmentRepo, allowOperatorPolicy{allowed: true})

		review := domain_fakes.NewReviewBuilder().Build()
		establishmentRepo.EXPECT().GetByID(ctx, review.EstablishmentID).Return(&domain.Establishment{ID: review.EstablishmentID, Name: "Cafe", Slug: "cafe", Category: "cafe"}, nil)
		highlightRepo.EXPECT().Delete(ctx, review.EstablishmentID, review.ID).Return(nil)

		err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		require.NoError(t, err)
	})

	t.Run("should return unauthorized when user cannot manage establishment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewDeleteHighlightUseCase(highlightRepo, establishmentRepo, allowOperatorPolicy{allowed: false})

		review := domain_fakes.NewReviewBuilder().Build()
		establishmentRepo.EXPECT().GetByID(ctx, review.EstablishmentID).Return(&domain.Establishment{ID: review.EstablishmentID, Name: "Cafe", Slug: "cafe", Category: "cafe"}, nil)

		err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		requireCustomErrorType(t, err, customerror.ForbiddenErrorType)
	})

	t.Run("should return validation error when establishmentID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewDeleteHighlightUseCase(highlightRepo, establishmentRepo, allowOperatorPolicy{allowed: true})

		err := uc.Execute(ctx, "11111111-1111-4111-8111-111111111111", "", "44444444-4444-4444-8444-444444444444")

		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return not found when establishment does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewDeleteHighlightUseCase(highlightRepo, establishmentRepo, allowOperatorPolicy{allowed: true})

		establishmentRepo.EXPECT().GetByID(ctx, "22222222-2222-4222-8222-222222222222").Return(nil, nil)

		err := uc.Execute(ctx, "11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222", "44444444-4444-4444-8444-444444444444")

		requireCustomErrorType(t, err, customerror.NotFoundErrorType)
	})

	t.Run("should return wrapped error when authorization fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		policyErr := errors.New("policy failed")
		uc := usecases.NewDeleteHighlightUseCase(highlightRepo, establishmentRepo, allowOperatorPolicy{err: policyErr})

		review := domain_fakes.NewReviewBuilder().Build()
		establishmentRepo.EXPECT().GetByID(ctx, review.EstablishmentID).Return(&domain.Establishment{ID: review.EstablishmentID, Name: "Cafe", Slug: "cafe", Category: "cafe"}, nil)

		err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to authorize establishment operator")
		assert.ErrorIs(t, err, policyErr)
	})

	t.Run("should return wrapped error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		establishmentRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewDeleteHighlightUseCase(highlightRepo, establishmentRepo, allowOperatorPolicy{allowed: true})

		review := domain_fakes.NewReviewBuilder().Build()
		establishmentRepo.EXPECT().GetByID(ctx, review.EstablishmentID).Return(&domain.Establishment{ID: review.EstablishmentID, Name: "Cafe", Slug: "cafe", Category: "cafe"}, nil)
		repoErr := errors.New("delete failed")
		highlightRepo.EXPECT().Delete(ctx, review.EstablishmentID, review.ID).Return(repoErr)

		err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to delete highlight")
		assert.ErrorIs(t, err, repoErr)
	})
}
