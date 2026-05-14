package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/get_review.go -package mock_usecases . GetReviewUseCase
type GetReviewUseCase interface {
	Execute(ctx context.Context, reviewID string) (*domain.Review, error)
}

type getReviewUseCase struct {
	reviewRepo domain.ReviewRepository
}

func NewGetReviewUseCase(reviewRepo domain.ReviewRepository) GetReviewUseCase {
	return &getReviewUseCase{reviewRepo: reviewRepo}
}

func (uc *getReviewUseCase) Execute(ctx context.Context, reviewID string) (*domain.Review, error) {
	if reviewID == "" {
		return nil, customerror.NewValidationError("review ID is required")
	}

	review, err := uc.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch review: %w", err)
	}
	if review == nil {
		return nil, customerror.NewNotFoundError("review not found")
	}

	return review, nil
}
