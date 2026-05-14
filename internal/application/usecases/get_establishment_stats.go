package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/get_establishment_stats.go -package mock_usecases . GetEstablishmentStatsUseCase
type GetEstablishmentStatsUseCase interface {
	Execute(ctx context.Context, establishmentID string) (*domain.EstablishmentStats, error)
}

type getEstablishmentStatsUseCase struct {
	establishmentRepo domain.EstablishmentRepository
}

func NewGetEstablishmentStatsUseCase(establishmentRepo domain.EstablishmentRepository) GetEstablishmentStatsUseCase {
	return &getEstablishmentStatsUseCase{establishmentRepo: establishmentRepo}
}

func (uc *getEstablishmentStatsUseCase) Execute(ctx context.Context, establishmentID string) (*domain.EstablishmentStats, error) {
	if establishmentID == "" {
		return nil, customerror.NewValidationError("establishment ID is required")
	}

	stats, err := uc.establishmentRepo.GetStats(ctx, establishmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch establishment stats: %w", err)
	}
	if stats == nil {
		return nil, customerror.NewNotFoundError("establishment stats not found")
	}
	if err := stats.Validate(); err != nil {
		return nil, err
	}

	return stats, nil
}
