package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/delete_highlight.go -package mock_usecases . DeleteHighlightUseCase
type DeleteHighlightUseCase interface {
	Execute(ctx context.Context, userID, establishmentID, reviewID string) error
}

type deleteHighlightUseCase struct {
	highlightRepo     domain.HighlightRepository
	establishmentRepo domain.EstablishmentRepository
	operatorPolicy    EstablishmentOperatorPolicy
}

func NewDeleteHighlightUseCase(
	highlightRepo domain.HighlightRepository,
	establishmentRepo domain.EstablishmentRepository,
	operatorPolicy EstablishmentOperatorPolicy,
) DeleteHighlightUseCase {
	return &deleteHighlightUseCase{
		highlightRepo:     highlightRepo,
		establishmentRepo: establishmentRepo,
		operatorPolicy:    operatorPolicy,
	}
}

func (uc *deleteHighlightUseCase) Execute(ctx context.Context, userID, establishmentID, reviewID string) error {
	if userID == "" {
		return customerror.NewUnauthorizedError("user not authenticated")
	}
	if establishmentID == "" {
		return customerror.NewValidationError("establishment ID is required")
	}
	if reviewID == "" {
		return customerror.NewValidationError("review ID is required")
	}

	establishment, err := uc.establishmentRepo.GetByID(ctx, establishmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch establishment: %w", err)
	}
	if establishment == nil {
		return customerror.NewNotFoundError("establishment not found")
	}

	allowed, err := uc.operatorPolicy.CanManageEstablishment(ctx, userID, establishmentID)
	if err != nil {
		return fmt.Errorf("failed to authorize establishment operator: %w", err)
	}
	if !allowed {
		return customerror.NewForbiddenError("user cannot manage this establishment")
	}

	if err := uc.highlightRepo.Delete(ctx, establishmentID, reviewID); err != nil {
		return fmt.Errorf("failed to delete highlight: %w", err)
	}

	return nil
}
