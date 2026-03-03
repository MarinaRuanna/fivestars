package usecases_test

import (
	"context"
	"errors"
	"fivestars/internal/application/usecases"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/domain/domain_fakes"
	"fivestars/internal/domain/mock_domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_CreateCheckinUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	radius := 100.0

	t.Run("should create checkin successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		estabRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewCreateCheckinUseCase(checkinRepo, estabRepo, radius)

		input := domain_fakes.NewCheckinBuilder().WithoutTimestamps().Build()
		estab := domain_fakes.NewEstablishmentBuilder().WithID(input.EstablishmentID).Build()

		estabRepo.EXPECT().GetByID(ctx, input.EstablishmentID).Return(&estab, nil)
		estabRepo.EXPECT().DistanceTo(ctx, input.EstablishmentID, input.Lat, input.Lng).Return(10.0, nil)
		checkinRepo.EXPECT().FindTodayByUserAndEstablishment(ctx, input.UserID, input.EstablishmentID).Return(nil, nil)
		checkinRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		result, err := uc.Execute(ctx, input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, input.UserID, result.UserID)
		assert.False(t, result.CheckedAt.IsZero())
		assert.False(t, result.CreatedAt.IsZero())
	})

	t.Run("should return validation error when userID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		estabRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewCreateCheckinUseCase(checkinRepo, estabRepo, radius)

		input := domain_fakes.NewCheckinBuilder().WithUserID("").Build()

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return not found when establishment does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		estabRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewCreateCheckinUseCase(checkinRepo, estabRepo, radius)

		input := domain_fakes.NewCheckinBuilder().Build()
		estabRepo.EXPECT().GetByID(ctx, input.EstablishmentID).Return(nil, nil)

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.NotFoundErrorType)
	})

	t.Run("should return wrapped error when fetching establishment fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		estabRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewCreateCheckinUseCase(checkinRepo, estabRepo, radius)

		input := domain_fakes.NewCheckinBuilder().Build()
		repoErr := errors.New("db failed")
		estabRepo.EXPECT().GetByID(ctx, input.EstablishmentID).Return(nil, repoErr)

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to fetch establishment")
		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("should return conflict when checkin already exists today", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		estabRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewCreateCheckinUseCase(checkinRepo, estabRepo, radius)

		input := domain_fakes.NewCheckinBuilder().Build()
		estab := domain_fakes.NewEstablishmentBuilder().WithID(input.EstablishmentID).Build()
		existing := domain_fakes.NewCheckinBuilder().WithTimestamps(time.Now().UTC(), time.Now().UTC()).Build()

		estabRepo.EXPECT().GetByID(ctx, input.EstablishmentID).Return(&estab, nil)
		estabRepo.EXPECT().DistanceTo(ctx, input.EstablishmentID, input.Lat, input.Lng).Return(5.0, nil)
		checkinRepo.EXPECT().FindTodayByUserAndEstablishment(ctx, input.UserID, input.EstablishmentID).Return(&existing, nil)

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ConflictErrorType)
	})

	t.Run("should return wrapped error when create fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		estabRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewCreateCheckinUseCase(checkinRepo, estabRepo, radius)

		input := domain_fakes.NewCheckinBuilder().Build()
		estab := domain_fakes.NewEstablishmentBuilder().WithID(input.EstablishmentID).Build()
		repoErr := errors.New("insert failed")

		estabRepo.EXPECT().GetByID(ctx, input.EstablishmentID).Return(&estab, nil)
		estabRepo.EXPECT().DistanceTo(ctx, input.EstablishmentID, input.Lat, input.Lng).Return(5.0, nil)
		checkinRepo.EXPECT().FindTodayByUserAndEstablishment(ctx, input.UserID, input.EstablishmentID).Return(nil, nil)
		checkinRepo.EXPECT().Create(ctx, gomock.Any()).Return(repoErr)

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to create checkin")
		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("should return conflict when database detects duplicated checkin day", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		estabRepo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewCreateCheckinUseCase(checkinRepo, estabRepo, radius)

		input := domain_fakes.NewCheckinBuilder().Build()
		estab := domain_fakes.NewEstablishmentBuilder().WithID(input.EstablishmentID).Build()

		estabRepo.EXPECT().GetByID(ctx, input.EstablishmentID).Return(&estab, nil)
		estabRepo.EXPECT().DistanceTo(ctx, input.EstablishmentID, input.Lat, input.Lng).Return(5.0, nil)
		checkinRepo.EXPECT().FindTodayByUserAndEstablishment(ctx, input.UserID, input.EstablishmentID).Return(nil, nil)
		checkinRepo.EXPECT().Create(ctx, gomock.Any()).Return(customerror.NewConflictError("checkin already exists"))

		result, err := uc.Execute(ctx, input)

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ConflictErrorType)
		assert.ErrorContains(t, err, "check-in already performed today for this establishment")
	})
}
