package usecases

import (
	"context"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/list_checkins.go -package mock_usecases . ListCheckinsUseCase
type ListCheckinsUseCase interface {
	Execute(ctx context.Context, userID string) ([]domain.Checkin, error)
}

type listCheckinsUseCase struct {
	checkinRepo domain.CheckinRepository
}

func NewListCheckinsUseCase(checkinRepo domain.CheckinRepository) ListCheckinsUseCase {
	return &listCheckinsUseCase{checkinRepo: checkinRepo}
}

func (uc *listCheckinsUseCase) Execute(ctx context.Context, userID string) ([]domain.Checkin, error) {
	if userID == "" {
		return nil, customerror.NewValidationError("user ID is required")
	}

	return uc.checkinRepo.ListByUser(ctx, userID)
}
