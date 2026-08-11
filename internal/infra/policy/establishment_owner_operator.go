package policy

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
)

type EstablishmentOwnerOperator struct {
	establishmentRepo domain.EstablishmentRepository
}

func NewEstablishmentOwnerOperator(establishmentRepo domain.EstablishmentRepository) EstablishmentOwnerOperator {
	return EstablishmentOwnerOperator{establishmentRepo: establishmentRepo}
}

func (p EstablishmentOwnerOperator) CanManageEstablishment(ctx context.Context, userID, establishmentID string) (bool, error) {
	establishment, err := p.establishmentRepo.GetByID(ctx, establishmentID)
	if err != nil {
		return false, fmt.Errorf("fetch establishment for operator check: %w", err)
	}
	if establishment == nil {
		return false, nil
	}

	return establishment.OwnerID != "" && establishment.OwnerID == userID, nil
}
