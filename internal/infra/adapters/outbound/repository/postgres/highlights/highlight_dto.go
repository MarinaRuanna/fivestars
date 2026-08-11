package highlights

import (
	"time"

	"fivestars/internal/domain"
)

type HighlightDTO struct {
	EstablishmentID string    `json:"establishment_id"`
	ReviewID        string    `json:"review_id"`
	CreatedByUserID string    `json:"created_by_user_id"`
	CreatedAt       time.Time `json:"created_at"`
}

func (d *HighlightDTO) ToDomain() (*domain.Highlight, error) {
	highlight := &domain.Highlight{
		EstablishmentID: d.EstablishmentID,
		ReviewID:        d.ReviewID,
		CreatedByUserID: d.CreatedByUserID,
		CreatedAt:       d.CreatedAt,
	}

	if err := highlight.Validate(); err != nil {
		return nil, err
	}

	return highlight, nil
}

func FromDomain(highlight *domain.Highlight) (*HighlightDTO, error) {
	if err := highlight.Validate(); err != nil {
		return nil, err
	}

	dto := &HighlightDTO{
		EstablishmentID: highlight.EstablishmentID,
		ReviewID:        highlight.ReviewID,
		CreatedByUserID: highlight.CreatedByUserID,
		CreatedAt:       highlight.CreatedAt,
	}

	if dto.CreatedAt.IsZero() {
		dto.CreatedAt = time.Now().UTC()
	}

	return dto, nil
}
