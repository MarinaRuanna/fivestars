package usecases

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/like_review.go -package mock_usecases . LikeReviewUseCase
type LikeReviewUseCase interface {
	Execute(ctx context.Context, userID, reviewID string) error
}

type likeReviewUseCase struct {
	reviewRepo domain.ReviewRepository
	likeRepo   domain.ReviewLikeRepository
}

func NewLikeReviewUseCase(reviewRepo domain.ReviewRepository, likeRepo domain.ReviewLikeRepository) LikeReviewUseCase {
	return &likeReviewUseCase{reviewRepo: reviewRepo, likeRepo: likeRepo}
}

func (uc *likeReviewUseCase) Execute(ctx context.Context, userID, reviewID string) error {
	if userID == "" {
		return customerror.NewUnauthorizedError("user not authenticated")
	}
	if reviewID == "" {
		return customerror.NewValidationError("review ID is required")
	}

	review, err := uc.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("failed to fetch review: %w", err)
	}
	if review == nil {
		return customerror.NewNotFoundError("review not found")
	}

	like := &domain.ReviewLike{UserID: userID, ReviewID: reviewID}
	if err := like.Validate(); err != nil {
		return err
	}

	if err := uc.likeRepo.Create(ctx, like); err != nil {
		if errorType, ok := customerror.TypeOf(err); ok && errorType == customerror.ConflictErrorType {
			return customerror.NewConflictError("review already liked")
		}
		return fmt.Errorf("failed to like review: %w", err)
	}

	return nil
}
