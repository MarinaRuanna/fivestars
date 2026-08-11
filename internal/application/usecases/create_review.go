package usecases

import (
	"context"
	"fmt"
	"time"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
)

const reviewWindow = 5 * 24 * time.Hour

//go:generate go run go.uber.org/mock/mockgen -destination mock_usecases/create_review.go -package mock_usecases . CreateReviewUseCase
type CreateReviewUseCase interface {
	Execute(ctx context.Context, input domain.Review) (*domain.Review, error)
}

type createReviewUseCase struct {
	reviewRepo  domain.ReviewRepository
	checkinRepo domain.CheckinRepository
}

func NewCreateReviewUseCase(reviewRepo domain.ReviewRepository, checkinRepo domain.CheckinRepository) CreateReviewUseCase {
	return &createReviewUseCase{reviewRepo: reviewRepo, checkinRepo: checkinRepo}
}

func (uc *createReviewUseCase) Execute(ctx context.Context, input domain.Review) (*domain.Review, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	checkin, err := uc.checkinRepo.GetByID(ctx, input.CheckinID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch checkin: %w", err)
	}
	if checkin == nil {
		return nil, customerror.NewNotFoundError("checkin not found")
	}
	if checkin.UserID != input.UserID {
		return nil, customerror.NewUnauthorizedError("checkin does not belong to user")
	}

	if err := input.EnsureWithinWindow(checkin, reviewWindow); err != nil {
		return nil, err
	}

	existing, err := uc.reviewRepo.GetByCheckinID(ctx, input.CheckinID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing review: %w", err)
	}
	if existing != nil {
		return nil, customerror.NewConflictError("review already exists for this checkin")
	}

	input.EstablishmentID = checkin.EstablishmentID

	review, err := domain.NewReview(input)
	if err != nil {
		return nil, err
	}

	if err := uc.reviewRepo.Create(ctx, review); err != nil {
		if errorType, ok := customerror.TypeOf(err); ok && errorType == customerror.ConflictErrorType {
			return nil, customerror.NewConflictError("review already exists for this checkin")
		}
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}
