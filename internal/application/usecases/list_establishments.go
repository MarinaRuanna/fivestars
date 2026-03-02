package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/list_establishments.go -package mock_usecases . ListEstablishmentsUseCase
type ListEstablishmentsUseCase interface {
	Execute(ctx context.Context) ([]domain.Establishment, error)
}

type listEstablishmentsUseCase struct {
	estabRepo domain.EstablishmentRepository
}

func NewListEstablishmentsUseCase(estabRepo domain.EstablishmentRepository) ListEstablishmentsUseCase {
	return &listEstablishmentsUseCase{
		estabRepo: estabRepo,
	}
}

func (s *listEstablishmentsUseCase) Execute(ctx context.Context) ([]domain.Establishment, error) {
	establishments, err := s.estabRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list establishments: %w", err)
	}

	if establishments == nil {
		return []domain.Establishment{}, nil
	}

	return establishments, nil
}
