package usecases_test

import (
	"context"
	"errors"
	"fivestars/internal/application/usecases"
	"fivestars/internal/domain"
	"fivestars/internal/domain/domain_fakes"
	"fivestars/internal/domain/mock_domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_ListEstablishmentsUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should return establishments successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewListEstablishmentsUseCase(repo)

		estab := domain_fakes.NewEstablishmentBuilder().Build()
		repo.EXPECT().List(ctx).Return([]domain.Establishment{estab}, nil)

		result, err := uc.Execute(ctx)

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, estab.ID, result[0].ID)
	})

	t.Run("should return empty slice when repository returns nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewListEstablishmentsUseCase(repo)

		repo.EXPECT().List(ctx).Return(nil, nil)

		result, err := uc.Execute(ctx)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("should return wrapped error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mock_domain.NewMockEstablishmentRepository(ctrl)
		uc := usecases.NewListEstablishmentsUseCase(repo)

		repoErr := errors.New("db failed")
		repo.EXPECT().List(ctx).Return(nil, repoErr)

		result, err := uc.Execute(ctx)

		require.Nil(t, result)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to list establishments")
		assert.ErrorIs(t, err, repoErr)
	})
}
