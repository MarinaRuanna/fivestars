package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/security"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/claim_establishment_ownership.go -package mock_usecases . ClaimEstablishmentOwnershipUseCase
type ClaimEstablishmentOwnershipUseCase interface {
	Execute(ctx context.Context, userID, establishmentID, claimQRCode string) (*domain.Establishment, error)
}

type claimEstablishmentOwnershipUseCase struct {
	establishmentRepo domain.EstablishmentRepository
	hasher            security.ClaimCodeHasher
}

func NewClaimEstablishmentOwnershipUseCase(
	establishmentRepo domain.EstablishmentRepository,
	hasher security.ClaimCodeHasher,
) ClaimEstablishmentOwnershipUseCase {
	return &claimEstablishmentOwnershipUseCase{
		establishmentRepo: establishmentRepo,
		hasher:            hasher,
	}
}

func (uc *claimEstablishmentOwnershipUseCase) Execute(ctx context.Context, userID, establishmentID, claimCode string) (*domain.Establishment, error) {
	if userID == "" {
		return nil, customerror.NewUnauthorizedError("user not authenticated")
	}
	if establishmentID == "" {
		return nil, customerror.NewValidationError("establishment ID is required")
	}
	if claimCode == "" {
		return nil, customerror.NewValidationError("claim_code is required")
	}

	claimCodeHash := uc.hasher.Hash(claimCode)

	establishment, err := uc.establishmentRepo.ClaimOwnership(ctx, establishmentID, userID, claimCodeHash)
	if err != nil {
		return nil, fmt.Errorf("failed to claim establishment ownership: %w", err)
	}
	if establishment == nil {
		return nil, customerror.NewNotFoundError("establishment not found")
	}

	return establishment, nil
}
