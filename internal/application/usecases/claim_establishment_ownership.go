package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/claim_establishment_ownership.go -package mock_usecases . ClaimEstablishmentOwnershipUseCase
type ClaimEstablishmentOwnershipUseCase interface {
	Execute(ctx context.Context, userID, establishmentID, claimQRCode string) (*domain.Establishment, error)
}

type claimEstablishmentOwnershipUseCase struct {
	establishmentRepo domain.EstablishmentRepository
}

func NewClaimEstablishmentOwnershipUseCase(establishmentRepo domain.EstablishmentRepository) ClaimEstablishmentOwnershipUseCase {
	return &claimEstablishmentOwnershipUseCase{establishmentRepo: establishmentRepo}
}

func (uc *claimEstablishmentOwnershipUseCase) Execute(ctx context.Context, userID, establishmentID, claimQRCode string) (*domain.Establishment, error) {
	if userID == "" {
		return nil, customerror.NewUnauthorizedError("user not authenticated")
	}
	if establishmentID == "" {
		return nil, customerror.NewValidationError("establishment ID is required")
	}
	if claimQRCode == "" {
		return nil, customerror.NewValidationError("qr_code is required")
	}

	establishment, err := uc.establishmentRepo.ClaimOwnership(ctx, establishmentID, userID, claimQRCode)
	if err != nil {
		return nil, fmt.Errorf("failed to claim establishment ownership: %w", err)
	}
	if establishment == nil {
		return nil, customerror.NewNotFoundError("establishment not found")
	}

	return establishment, nil
}
