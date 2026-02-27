package checkins

import (
	"time"

	"fivestars/internal/domain"
)

type CheckinDTO struct {
	ID              string    `json:"checkin_id"`
	UserID          string    `json:"user_id"`
	EstablishmentID string    `json:"establishment_id"`
	Lat             float64   `json:"lat"`
	Lng             float64   `json:"lng"`
	CheckedAt       time.Time `json:"checked_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func (d *CheckinDTO) ToDomain() *domain.Checkin {
	return &domain.Checkin{
		ID:              d.ID,
		UserID:          d.UserID,
		EstablishmentID: d.EstablishmentID,
		Lat:             d.Lat,
		Lng:             d.Lng,
		CheckedAt:       d.CheckedAt,
		CreatedAt:       d.CreatedAt,
	}
}

func FromDomain(c *domain.Checkin) *CheckinDTO {
	return &CheckinDTO{
		ID:              c.ID,
		UserID:          c.UserID,
		EstablishmentID: c.EstablishmentID,
		Lat:             c.Lat,
		Lng:             c.Lng,
		CheckedAt:       c.CheckedAt,
		CreatedAt:       c.CreatedAt,
	}
}
