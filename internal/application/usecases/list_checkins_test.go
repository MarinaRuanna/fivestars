package usecases_test

import (
	"context"
	"errors"
	"fivestars/internal/application/usecases"
	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/domain/domain_fakes"
	"fivestars/internal/domain/mock_domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_ListCheckinsUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	userID := "11111111-1111-4111-8111-111111111111"

	t.Run("should return checkins successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		uc := usecases.NewListCheckinsUseCase(checkinRepo)

		checkin := domain_fakes.NewCheckinBuilder().WithUserID(userID).Build()
		checkinRepo.EXPECT().ListByUser(ctx, userID).Return([]domain.Checkin{checkin}, nil)

		result, err := uc.Execute(ctx, userID)

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, checkin.ID, result[0].ID)
	})

	t.Run("should return validation error when user id is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		uc := usecases.NewListCheckinsUseCase(checkinRepo)

		result, err := uc.Execute(ctx, "")

		require.Nil(t, result)
		requireCustomErrorType(t, err, customerror.ValidationErrorType)
	})

	t.Run("should return repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		checkinRepo := mock_domain.NewMockCheckinRepository(ctrl)
		uc := usecases.NewListCheckinsUseCase(checkinRepo)

		repoErr := errors.New("query failed")
		checkinRepo.EXPECT().ListByUser(ctx, userID).Return(nil, repoErr)

		result, err := uc.Execute(ctx, userID)

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
	})
}
