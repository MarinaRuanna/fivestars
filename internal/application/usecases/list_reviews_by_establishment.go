package usecases

import (
	"context"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/list_reviews_by_establishment.go -package mock_usecases . ListReviewsByEstablishmentUseCase
type ListReviewsByEstablishmentUseCase interface {
	Execute(ctx context.Context, input domain.ReviewListOptions, establishmentID string) ([]domain.Review, error)
}

type listReviewsByEstablishmentUseCase struct {
	reviewRepo domain.ReviewRepository
}

func NewListReviewsByEstablishmentUseCase(reviewRepo domain.ReviewRepository) ListReviewsByEstablishmentUseCase {
	return &listReviewsByEstablishmentUseCase{reviewRepo: reviewRepo}
}

func (uc *listReviewsByEstablishmentUseCase) Execute(ctx context.Context, input domain.ReviewListOptions, establishmentID string) ([]domain.Review, error) {
	if establishmentID == "" {
		return nil, customerror.NewValidationError("establishment ID is required")
	}

	if input.Limit == 0 {
		input.Limit = 20
	}
	if input.Limit < 0 || input.Offset < 0 {
		return nil, customerror.NewValidationError("invalid pagination values")
	}
	if input.Limit > 100 {
		return nil, customerror.NewValidationError("limit exceeds maximum")
	}
	if input.MinRating != nil && (*input.MinRating < 1 || *input.MinRating > 5) {
		return nil, customerror.NewValidationError("invalid min rating")
	}
	if input.MaxRating != nil && (*input.MaxRating < 1 || *input.MaxRating > 5) {
		return nil, customerror.NewValidationError("invalid max rating")
	}
	if input.MinRating != nil && input.MaxRating != nil && *input.MinRating > *input.MaxRating {
		return nil, customerror.NewValidationError("min rating greater than max rating")
	}

	list, err := uc.reviewRepo.ListByEstablishment(ctx, establishmentID, input)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return []domain.Review{}, nil
	}
	return list, nil
}
