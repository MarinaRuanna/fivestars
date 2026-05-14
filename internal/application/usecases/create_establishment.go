package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/create_establishment.go -package mock_usecases . CreateEstablishmentUseCase
type CreateEstablishmentUseCase interface {
	Execute(ctx context.Context, userID string, input domain.Establishment) (*domain.Establishment, error)
}

type createEstablishmentUseCase struct {
	establishmentRepo domain.EstablishmentRepository
}

func NewCreateEstablishmentUseCase(establishmentRepo domain.EstablishmentRepository) CreateEstablishmentUseCase {
	return &createEstablishmentUseCase{establishmentRepo: establishmentRepo}
}

func (uc *createEstablishmentUseCase) Execute(ctx context.Context, userID string, input domain.Establishment) (*domain.Establishment, error) {
	if userID == "" {
		return nil, customerror.NewUnauthorizedError("user not authenticated")
	}

	establishment, err := domain.NewEstablishment(domain.Establishment{
		OwnerID:  userID,
		Name:     input.Name,
		Slug:     input.Slug,
		Category: input.Category,
		Address:  input.Address,
		Lat:      input.Lat,
		Lng:      input.Lng,
		QRCode:   input.QRCode,
	})
	if err != nil {
		return nil, err
	}

	if err := uc.establishmentRepo.Create(ctx, establishment); err != nil {
		return nil, fmt.Errorf("failed to create establishment: %w", err)
	}

	return establishment, nil
}
