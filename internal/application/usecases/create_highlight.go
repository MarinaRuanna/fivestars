package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/create_highlight.go -package mock_usecases . CreateHighlightUseCase
type CreateHighlightUseCase interface {
	Execute(ctx context.Context, userID, establishmentID, reviewID string) (*domain.Highlight, error)
}

type EstablishmentOperatorPolicy interface {
	CanManageEstablishment(ctx context.Context, userID, establishmentID string) (bool, error)
}

type createHighlightUseCase struct {
	highlightRepo  domain.HighlightRepository
	reviewRepo     domain.ReviewRepository
	operatorPolicy EstablishmentOperatorPolicy
}

func NewCreateHighlightUseCase(
	highlightRepo domain.HighlightRepository,
	reviewRepo domain.ReviewRepository,
	operatorPolicy EstablishmentOperatorPolicy,
) CreateHighlightUseCase {
	return &createHighlightUseCase{
		highlightRepo:  highlightRepo,
		reviewRepo:     reviewRepo,
		operatorPolicy: operatorPolicy,
	}
}

func (uc *createHighlightUseCase) Execute(ctx context.Context, userID, establishmentID, reviewID string) (*domain.Highlight, error) {
	if userID == "" {
		return nil, customerror.NewUnauthorizedError("user not authenticated")
	}
	if establishmentID == "" {
		return nil, customerror.NewValidationError("establishment ID is required")
	}
	if reviewID == "" {
		return nil, customerror.NewValidationError("review ID is required")
	}

	allowed, err := uc.operatorPolicy.CanManageEstablishment(ctx, userID, establishmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to authorize establishment operator: %w", err)
	}
	if !allowed {
		return nil, customerror.NewForbiddenError("user cannot manage this establishment")
	}

	review, err := uc.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch review: %w", err)
	}
	if review == nil {
		return nil, customerror.NewNotFoundError("review not found")
	}
	if review.EstablishmentID != establishmentID {
		return nil, customerror.NewValidationError("review does not belong to this establishment")
	}

	exists, err := uc.highlightRepo.Exists(ctx, establishmentID, reviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing highlight: %w", err)
	}
	if exists {
		return nil, customerror.NewConflictError("review already highlighted")
	}

	count, err := uc.highlightRepo.CountByEstablishment(ctx, establishmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to count establishment highlights: %w", err)
	}
	if count >= domain.MaxHighlightsPerEstablishment {
		return nil, customerror.NewConflictError("highlight limit reached")
	}

	highlight, err := domain.NewHighlight(domain.Highlight{
		EstablishmentID: establishmentID,
		ReviewID:        reviewID,
		CreatedByUserID: userID,
	})
	if err != nil {
		return nil, err
	}

	if err := uc.highlightRepo.Create(ctx, highlight); err != nil {
		return nil, fmt.Errorf("failed to create highlight: %w", err)
	}

	return highlight, nil
}
