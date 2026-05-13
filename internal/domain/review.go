package domain

import (
	"context"
	"strings"
	"time"

	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
)

const MinReviewBodyLen = 10

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/review_repository.go -package mock_domain . ReviewRepository
type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	GetByID(ctx context.Context, reviewID string) (*Review, error)
	GetByCheckinID(ctx context.Context, checkinID string) (*Review, error)
	ListByEstablishment(ctx context.Context, establishmentID string, options ReviewListOptions) ([]Review, error)
}

type Review struct {
	ID              string
	UserID          string `validate:"required,uuid4"`
	EstablishmentID string `validate:"required,uuid4"`
	CheckinID       string `validate:"required,uuid4"`
	Rating          int    `validate:"required,min=1,max=5"`
	Text            string `validate:"required"`
	LikeCount       int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ReviewListOptions struct {
	Limit     int
	Offset    int
	MinRating *int
	MaxRating *int
}

func NewReview(input Review) (*Review, error) {
	r := input
	now := time.Now().UTC()
	if r.CreatedAt.IsZero() {
		r.CreatedAt = now
	}
	if r.UpdatedAt.IsZero() {
		r.UpdatedAt = now
	}
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return &r, nil
}

func (r *Review) Validate() error {
	if err := validator.Validate(r); err != nil {
		return customerror.NewValidationError(err.Error())
	}

	if len(strings.TrimSpace(r.Text)) < MinReviewBodyLen {
		return customerror.NewValidationError("review text is too short")
	}

	return nil
}

func (r *Review) EnsureWithinWindow(checkin *Checkin, window time.Duration) error {
	now := time.Now().UTC()
	if now.After(checkin.CheckedAt.Add(window)) {
		return customerror.NewValidationError("review window expired")
	}
	return nil
}
