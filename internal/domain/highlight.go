package domain

import (
	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
	"time"
)

const MaxHighlightsPerEstablishment = 3

type Highlight struct {
	EstablishmentID string `validate:"required,uuid4"`
	ReviewID        string `validate:"required,uuid4"`
	CreatedByUserID string `validate:"required,uuid4"`
	CreatedAt       time.Time
}

type EstablishmentStats struct {
	AverageRating          float64
	TotalReviews           int
	TotalLikes             int
	HighlightedReviewCount int
}

func NewHighlight(input Highlight) (*Highlight, error) {
	highlight := input
	if highlight.CreatedAt.IsZero() {
		highlight.CreatedAt = time.Now().UTC()
	}

	if err := highlight.Validate(); err != nil {
		return nil, err
	}

	return &highlight, nil
}

func (h *Highlight) Validate() error {
	if err := validator.Validate(h); err != nil {
		return customerror.NewValidationError(err.Error())
	}

	return nil
}

func (s *EstablishmentStats) Validate() error {
	if s.AverageRating < 0 || s.AverageRating > 5 {
		return customerror.NewValidationError("average rating must be between 0 and 5")
	}
	if s.TotalReviews < 0 || s.TotalLikes < 0 || s.HighlightedReviewCount < 0 {
		return customerror.NewValidationError("establishment stats counts must be non-negative")
	}

	return nil
}
