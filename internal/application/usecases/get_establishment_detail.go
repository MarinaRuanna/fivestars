package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/get_establishment_detail.go -package mock_usecases . GetEstablishmentDetailUseCase
type GetEstablishmentDetailUseCase interface {
	Execute(ctx context.Context, establishmentID string) (*domain.EstablishmentDetail, error)
}

type getEstablishmentDetailUseCase struct {
	establishmentRepo domain.EstablishmentRepository
	highlightRepo     domain.HighlightRepository
	reviewRepo        domain.ReviewRepository
}

func NewGetEstablishmentDetailUseCase(
	establishmentRepo domain.EstablishmentRepository,
	highlightRepo domain.HighlightRepository,
	reviewRepo domain.ReviewRepository,
) GetEstablishmentDetailUseCase {
	return &getEstablishmentDetailUseCase{
		establishmentRepo: establishmentRepo,
		highlightRepo:     highlightRepo,
		reviewRepo:        reviewRepo,
	}
}

func (uc *getEstablishmentDetailUseCase) Execute(ctx context.Context, establishmentID string) (*domain.EstablishmentDetail, error) {
	if establishmentID == "" {
		return nil, customerror.NewValidationError("establishment ID is required")
	}

	establishment, err := uc.establishmentRepo.GetByID(ctx, establishmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch establishment: %w", err)
	}
	if establishment == nil {
		return nil, customerror.NewNotFoundError("establishment not found")
	}

	highlights, err := uc.highlightRepo.ListByEstablishment(ctx, establishmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list establishment highlights: %w", err)
	}

	detail := &domain.EstablishmentDetail{
		Establishment: *establishment,
		Highlights:    []domain.Review{},
	}

	for _, highlight := range highlights {
		review, err := uc.reviewRepo.GetByID(ctx, highlight.ReviewID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch highlighted review: %w", err)
		}
		if review == nil {
			continue
		}

		detail.Highlights = append(detail.Highlights, *review)
	}

	return detail, nil
}
