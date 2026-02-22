package application

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
)

type EstablishmentService interface {
	ListEstablishments(ctx context.Context) ([]domain.Establishment, error)
}

type establishmentService struct {
	estabRepo domain.EstablishmentRepository
}

func NewEstablishmentService(estabRepo domain.EstablishmentRepository) EstablishmentService {
	return &establishmentService{
		estabRepo: estabRepo,
	}
}

func (s *establishmentService) ListEstablishments(ctx context.Context) ([]domain.Establishment, error) {
	// 1. FETCH ALL ESTABLISHMENTS
	establishments, err := s.estabRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list establishments: %w", err)
	}

	if establishments == nil {
		establishments = []*domain.Establishment{}
	}

	// 2. MAP TO OUTPUT (domain → output DTO)
	outputs := make([]domain.Establishment, len(establishments))
	for i, estab := range establishments {
		outputs[i] = domain.Establishment{
			ID:       estab.ID,
			Name:     estab.Name,
			Category: estab.Category,
		}
	}

	// 3. RETURN OUTPUT
	return outputs, nil
}
