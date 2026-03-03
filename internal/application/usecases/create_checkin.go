package usecases

import (
	"context"
	"fmt"
	"time"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/create_checkin.go -package mock_usecases . CreateCheckinUseCase
type CreateCheckinUseCase interface {
	Execute(ctx context.Context, input domain.Checkin) (*domain.Checkin, error)
}

type createCheckinUseCase struct {
	checkinRepo  domain.CheckinRepository
	estabRepo    domain.EstablishmentRepository
	radiusMeters float64
}

func NewCreateCheckinUseCase(checkinRepo domain.CheckinRepository, estabRepo domain.EstablishmentRepository, radiusMeters float64) CreateCheckinUseCase {
	return &createCheckinUseCase{checkinRepo: checkinRepo, estabRepo: estabRepo, radiusMeters: radiusMeters}
}

func (uc *createCheckinUseCase) Execute(ctx context.Context, input domain.Checkin) (*domain.Checkin, error) {
	if input.UserID == "" {
		return nil, customerror.NewValidationError("user ID is required")
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// ensure establishment exists
	estab, err := uc.estabRepo.GetByID(ctx, input.EstablishmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch establishment: %w", err)
	}
	if estab == nil {
		return nil, customerror.NewNotFoundError("establishment not found")
	}
	if estab.Lat == nil || estab.Lng == nil {
		return nil, customerror.NewValidationError("establishment has no location")
	}

	// check distance using repository (PostGIS)
	distance, err := uc.estabRepo.DistanceTo(ctx, input.EstablishmentID, input.Lat, input.Lng)
	if err != nil {
		return nil, fmt.Errorf("failed to compute distance: %w", err)
	}
	if distance > uc.radiusMeters {
		return nil, customerror.NewValidationError("user too far from establishment")
	}

	// enforce 1 check-in per day for same user+establishment
	existing, err := uc.checkinRepo.FindTodayByUserAndEstablishment(ctx, input.UserID, input.EstablishmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing checkin: %w", err)
	}
	if existing != nil {
		return nil, customerror.NewConflictError("check-in already performed today for this establishment")
	}

	// set timestamps
	if input.CheckedAt.IsZero() {
		input.CheckedAt = time.Now().UTC()
	}
	if input.CreatedAt.IsZero() {
		input.CreatedAt = time.Now().UTC()
	}

	// persist
	if err := uc.checkinRepo.Create(ctx, &input); err != nil {
		if errorType, ok := customerror.TypeOf(err); ok && errorType == customerror.ConflictErrorType {
			return nil, customerror.NewConflictError("check-in already performed today for this establishment")
		}
		return nil, fmt.Errorf("failed to create checkin: %w", err)
	}

	return &input, nil
}
