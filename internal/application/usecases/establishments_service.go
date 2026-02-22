package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
)

// ListEstablishmentsUseCase implements establishment listing business logic.
// ⭐ Completely isolated from HTTP and database specifics.
type ListEstablishmentsUseCase struct {
	estabRepo domain.EstablishmentRepository
}

// NewListEstablishmentsUseCase creates a new ListEstablishmentsUseCase.
func NewListEstablishmentsUseCase(estabRepo domain.EstablishmentRepository) *ListEstablishmentsUseCase {
	return &ListEstablishmentsUseCase{
		estabRepo: estabRepo,
	}
}

// EstablishmentOutput DTO for the use case output.
type EstablishmentOutput struct {
	ID       string
	Name     string
	Category string
}

// Execute runs the establishment listing logic.
func (uc *ListEstablishmentsUseCase) Execute(ctx context.Context) (*[]EstablishmentOutput, error) {
	// 1. FETCH ALL ESTABLISHMENTS
	establishments, err := uc.estabRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list establishments: %w", err)
	}

	if establishments == nil {
		establishments = []*domain.Establishment{}
	}

	// 2. MAP TO OUTPUT (domain → output DTO)
	outputs := make([]EstablishmentOutput, len(establishments))
	for i, estab := range establishments {
		outputs[i] = EstablishmentOutput{
			ID:       estab.ID,
			Name:     estab.Name,
			Category: estab.Category,
		}
	}

	// 3. RETURN OUTPUT
	return &outputs, nil
}
