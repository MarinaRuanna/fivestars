package domain

import (
	"context"
	"time"

	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/review_like_repository.go -package mock_domain . ReviewLikeRepository
type ReviewLikeRepository interface {
	Create(ctx context.Context, like *ReviewLike) error
	Delete(ctx context.Context, userID, reviewID string) error
}

type ReviewLike struct {
	UserID    string `validate:"required,uuid4"`
	ReviewID  string `validate:"required,uuid4"`
	CreatedAt time.Time
}

func (l *ReviewLike) Validate() error {
	if err := validator.Validate(l); err != nil {
		return customerror.NewValidationError(err.Error())
	}
	return nil
}
