package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/unlike_review.go -package mock_usecases . UnlikeReviewUseCase
type UnlikeReviewUseCase interface {
	Execute(ctx context.Context, userID, reviewID string) error
}

type unlikeReviewUseCase struct {
	likeRepo domain.ReviewLikeRepository
}

func NewUnlikeReviewUseCase(likeRepo domain.ReviewLikeRepository) UnlikeReviewUseCase {
	return &unlikeReviewUseCase{likeRepo: likeRepo}
}

func (uc *unlikeReviewUseCase) Execute(ctx context.Context, userID, reviewID string) error {
	if userID == "" {
		return customerror.NewUnauthorizedError("user not authenticated")
	}
	if reviewID == "" {
		return customerror.NewValidationError("review ID is required")
	}

	if err := uc.likeRepo.Delete(ctx, userID, reviewID); err != nil {
		return fmt.Errorf("failed to unlike review: %w", err)
	}

	return nil
}
