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

type allowOperatorPolicy struct {
	allowed bool
	err     error
}

func (p allowOperatorPolicy) CanManageEstablishment(ctx context.Context, userID, establishmentID string) (bool, error) {
	return p.allowed, p.err
}

func Test_CreateHighlightUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should create highlight successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewCreateHighlightUseCase(highlightRepo, reviewRepo, allowOperatorPolicy{allowed: true})

		review := domain_fakes.NewReviewBuilder().Build()

		reviewRepo.EXPECT().GetByID(ctx, review.ID).Return(&review, nil)
		highlightRepo.EXPECT().Exists(ctx, review.EstablishmentID, review.ID).Return(false, nil)
		highlightRepo.EXPECT().CountByEstablishment(ctx, review.EstablishmentID).Return(0, nil)
		highlightRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		result, err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, review.EstablishmentID, result.EstablishmentID)
		assert.Equal(t, review.ID, result.ReviewID)
		assert.Equal(t, review.UserID, result.CreatedByUserID)
		assert.False(t, result.CreatedAt.IsZero())
	})

	t.Run("should return unauthorized when user cannot manage establishment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewCreateHighlightUseCase(highlightRepo, reviewRepo, allowOperatorPolicy{allowed: false})

		review := domain_fakes.NewReviewBuilder().Build()

		result, err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ForbiddenErrorType)
	})

	t.Run("should return not found when review does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewCreateHighlightUseCase(highlightRepo, reviewRepo, allowOperatorPolicy{allowed: true})

		reviewRepo.EXPECT().GetByID(ctx, "missing").Return(nil, nil)

		result, err := uc.Execute(ctx, "11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222", "missing")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.NotFoundErrorType)
	})

	t.Run("should return validation error when review belongs to another establishment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewCreateHighlightUseCase(highlightRepo, reviewRepo, allowOperatorPolicy{allowed: true})

		review := domain_fakes.NewReviewBuilder().Build()
		reviewRepo.EXPECT().GetByID(ctx, review.ID).Return(&review, nil)

		result, err := uc.Execute(ctx, review.UserID, "99999999-9999-4999-8999-999999999999", review.ID)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return conflict when highlight already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewCreateHighlightUseCase(highlightRepo, reviewRepo, allowOperatorPolicy{allowed: true})

		review := domain_fakes.NewReviewBuilder().Build()
		reviewRepo.EXPECT().GetByID(ctx, review.ID).Return(&review, nil)
		highlightRepo.EXPECT().Exists(ctx, review.EstablishmentID, review.ID).Return(true, nil)

		result, err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ConflictErrorType)
	})

	t.Run("should return conflict when highlight limit is reached", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		uc := usecases.NewCreateHighlightUseCase(highlightRepo, reviewRepo, allowOperatorPolicy{allowed: true})

		review := domain_fakes.NewReviewBuilder().Build()
		reviewRepo.EXPECT().GetByID(ctx, review.ID).Return(&review, nil)
		highlightRepo.EXPECT().Exists(ctx, review.EstablishmentID, review.ID).Return(false, nil)
		highlightRepo.EXPECT().CountByEstablishment(ctx, review.EstablishmentID).Return(domain.MaxHighlightsPerEstablishment, nil)

		result, err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ConflictErrorType)
	})

	t.Run("should return wrapped error when policy lookup fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		highlightRepo := mock_domain.NewMockHighlightRepository(ctrl)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		policyErr := errors.New("policy failed")
		uc := usecases.NewCreateHighlightUseCase(highlightRepo, reviewRepo, allowOperatorPolicy{err: policyErr})

		review := domain_fakes.NewReviewBuilder().Build()

		result, err := uc.Execute(ctx, review.UserID, review.EstablishmentID, review.ID)

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to authorize establishment operator")
		assert.ErrorIs(t, err, policyErr)
	})
}
