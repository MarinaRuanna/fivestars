package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"fivestars/internal/application/usecases"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/domain/domain_fakes"
	"fivestars/internal/domain/mock_domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_CreateReviewUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should create review successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		uc := usecases.NewCreateReviewUseCase(reviewRepo, checkinRepo)

		now := time.Now().UTC()
		checkin := domain_fakes.NewCheckinBuilder().WithTimestamps(now.Add(-24*time.Hour), now.Add(-24*time.Hour)).Build()
		review := domain_fakes.NewReviewBuilder().WithCheckinID(checkin.ID).WithUserID(checkin.UserID).WithoutTimestamps().Build()

		checkinRepo.EXPECT().GetByID(ctx, review.CheckinID).Return(&checkin, nil)
		reviewRepo.EXPECT().GetByCheckinID(ctx, review.CheckinID).Return(nil, nil)
		reviewRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		result, err := uc.Execute(ctx, review)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, checkin.EstablishmentID, result.EstablishmentID)
		assert.False(t, result.CreatedAt.IsZero())
		assert.False(t, result.UpdatedAt.IsZero())
	})

	t.Run("should return validation error when userID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		uc := usecases.NewCreateReviewUseCase(reviewRepo, checkinRepo)

		review := domain_fakes.NewReviewBuilder().WithUserID("").Build()

		result, err := uc.Execute(ctx, review)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return not found when checkin does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		uc := usecases.NewCreateReviewUseCase(reviewRepo, checkinRepo)

		review := domain_fakes.NewReviewBuilder().Build()
		checkinRepo.EXPECT().GetByID(ctx, review.CheckinID).Return(nil, nil)

		result, err := uc.Execute(ctx, review)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.NotFoundErrorType)
	})

	t.Run("should return unauthorized when checkin belongs to another user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		uc := usecases.NewCreateReviewUseCase(reviewRepo, checkinRepo)

		review := domain_fakes.NewReviewBuilder().Build()
		checkin := domain_fakes.NewCheckinBuilder().WithUserID("99999999-9999-4999-8999-999999999999").Build()

		checkinRepo.EXPECT().GetByID(ctx, review.CheckinID).Return(&checkin, nil)

		result, err := uc.Execute(ctx, review)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.UnauthorizedErrorType)
	})

	t.Run("should return validation error when review window expired", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		uc := usecases.NewCreateReviewUseCase(reviewRepo, checkinRepo)

		now := time.Now().UTC()
		checkin := domain_fakes.NewCheckinBuilder().WithTimestamps(now.Add(-10*24*time.Hour), now.Add(-10*24*time.Hour)).Build()
		review := domain_fakes.NewReviewBuilder().WithCheckinID(checkin.ID).WithUserID(checkin.UserID).Build()

		checkinRepo.EXPECT().GetByID(ctx, review.CheckinID).Return(&checkin, nil)

		result, err := uc.Execute(ctx, review)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return conflict when review already exists for checkin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		uc := usecases.NewCreateReviewUseCase(reviewRepo, checkinRepo)

		now := time.Now().UTC()
		checkin := domain_fakes.NewCheckinBuilder().WithTimestamps(now.Add(-24*time.Hour), now.Add(-24*time.Hour)).Build()
		review := domain_fakes.NewReviewBuilder().WithCheckinID(checkin.ID).WithUserID(checkin.UserID).Build()
		existing := domain_fakes.NewReviewBuilder().Build()

		checkinRepo.EXPECT().GetByID(ctx, review.CheckinID).Return(&checkin, nil)
		reviewRepo.EXPECT().GetByCheckinID(ctx, review.CheckinID).Return(&existing, nil)

		result, err := uc.Execute(ctx, review)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ConflictErrorType)
	})

	t.Run("should return wrapped error when checkin lookup fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		uc := usecases.NewCreateReviewUseCase(reviewRepo, checkinRepo)

		review := domain_fakes.NewReviewBuilder().Build()
		repoErr := errors.New("db failed")
		checkinRepo.EXPECT().GetByID(ctx, review.CheckinID).Return(nil, repoErr)

		result, err := uc.Execute(ctx, review)

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to fetch checkin")
		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("should return wrapped error when create fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		uc := usecases.NewCreateReviewUseCase(reviewRepo, checkinRepo)

		now := time.Now().UTC()
		checkin := domain_fakes.NewCheckinBuilder().WithTimestamps(now.Add(-24*time.Hour), now.Add(-24*time.Hour)).Build()
		review := domain_fakes.NewReviewBuilder().WithCheckinID(checkin.ID).WithUserID(checkin.UserID).Build()
		repoErr := errors.New("insert failed")

		checkinRepo.EXPECT().GetByID(ctx, review.CheckinID).Return(&checkin, nil)
		reviewRepo.EXPECT().GetByCheckinID(ctx, review.CheckinID).Return(nil, nil)
		reviewRepo.EXPECT().Create(ctx, gomock.Any()).Return(repoErr)

		result, err := uc.Execute(ctx, review)

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to create review")
		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("should return conflict when database detects duplicated review", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		reviewRepo := mock_domain.NewMockReviewRepository(ctrl)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		uc := usecases.NewCreateReviewUseCase(reviewRepo, checkinRepo)

		now := time.Now().UTC()
		checkin := domain_fakes.NewCheckinBuilder().WithTimestamps(now.Add(-24*time.Hour), now.Add(-24*time.Hour)).Build()
		review := domain_fakes.NewReviewBuilder().WithCheckinID(checkin.ID).WithUserID(checkin.UserID).Build()

		checkinRepo.EXPECT().GetByID(ctx, review.CheckinID).Return(&checkin, nil)
		reviewRepo.EXPECT().GetByCheckinID(ctx, review.CheckinID).Return(nil, nil)
		reviewRepo.EXPECT().Create(ctx, gomock.Any()).Return(customerror.NewConflictError("review exists"))

		result, err := uc.Execute(ctx, review)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ConflictErrorType)
		assert.ErrorContains(t, err, "review already exists for this checkin")
	})
}
