package reviews

import (
	"strings"
	"time"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
)

type ReviewDTO struct {
	ID              string    `json:"review_id"`
	UserID          string    `json:"user_id"`
	EstablishmentID string    `json:"establishment_id"`
	CheckinID       string    `json:"checkin_id"`
	Rating          int       `json:"rating"`
	Text            string    `json:"text"`
	LikeCount       int       `json:"like_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (r *ReviewDTO) Validate() error {
	if err := validator.Validate(r); err != nil {
		return customerror.NewValidationError(err.Error())
	}

	if len(strings.TrimSpace(r.Text)) < domain.MinReviewBodyLen {
		return customerror.NewValidationError("review text is too short")
	}

	return nil
}

func (d *ReviewDTO) ToDomain() (*domain.Review, error) {
	review := &domain.Review{
		ID:              d.ID,
		UserID:          d.UserID,
		EstablishmentID: d.EstablishmentID,
		CheckinID:       d.CheckinID,
		Rating:          d.Rating,
		Text:            d.Text,
		LikeCount:       d.LikeCount,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}

	err := review.Validate()
	if err != nil {
		return nil, err
	}

	return review, nil
}

func FromDomain(r *domain.Review) (*ReviewDTO, error) {
	dto := &ReviewDTO{
		ID:              r.ID,
		UserID:          r.UserID,
		EstablishmentID: r.EstablishmentID,
		CheckinID:       r.CheckinID,
		Rating:          r.Rating,
		Text:            r.Text,
		LikeCount:       r.LikeCount,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}

	err := dto.Validate()
	if err != nil {
		return nil, err
	}

	return dto, nil

}
