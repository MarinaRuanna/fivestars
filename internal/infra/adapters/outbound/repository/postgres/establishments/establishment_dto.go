package establishments

import (
	"time"

	"fivestars/internal/domain"
)

type EstablishmentDTO struct {
	ID        string    `json:"establishment_id" validate:"required,uuid4"`
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

func (r *EstablishmentDTO) ToDomain() (*domain.Establishment, error) {
	estab := &domain.Establishment{
		ID:        r.ID,
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
