package establishments

import (
	"time"

	"fivestars/internal/domain"
)

type EstablishmentDTO struct {
	ID        string    `json:"establishment_id" validate:"required,uuid4"`
	OwnerID   string    `json:"owner_id"`
	Name      string    `json:"name" validate:"required"`
	Slug      string    `json:"slug"`
	Category  string    `json:"category" validate:"required"`
	Address   string    `json:"address,omitempty"`
	Lat       *float64  `json:"lat,omitempty"`
	Lng       *float64  `json:"lng,omitempty"`
	QRCode    string    `json:"qr_code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EstablishmentStatsDTO struct {
	AverageRating          float64 `json:"average_rating"`
	TotalReviews           int     `json:"total_reviews"`
	TotalLikes             int     `json:"total_likes"`
	HighlightedReviewCount int     `json:"highlighted_review_count"`
}

func (r *EstablishmentDTO) ToDomain() (*domain.Establishment, error) {
	estab := &domain.Establishment{
		ID:        r.ID,
		OwnerID:   r.OwnerID,
		Name:      r.Name,
		Slug:      r.Slug,
		Category:  r.Category,
		Address:   r.Address,
		Lat:       r.Lat,
		Lng:       r.Lng,
		QRCode:    r.QRCode,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
	if err := estab.Validate(); err != nil {
		return nil, err
	}
	return estab, nil
}

func (d *EstablishmentStatsDTO) ToDomain() (*domain.EstablishmentStats, error) {
	stats := &domain.EstablishmentStats{
		AverageRating:          d.AverageRating,
		TotalReviews:           d.TotalReviews,
		TotalLikes:             d.TotalLikes,
		HighlightedReviewCount: d.HighlightedReviewCount,
	}

	if err := stats.Validate(); err != nil {
		return nil, err
	}

	return stats, nil
}
