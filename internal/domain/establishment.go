package domain

import (
	"context"
	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
	"time"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/establishment_repository.go -package mock_domain . EstablishmentRepository
type EstablishmentRepository interface {
	List(ctx context.Context) ([]*Establishment, error)
}

// Establishment represents a place that can be checked in and reviewed.
type Establishment struct {
	ID        string    `json:"id" validade:"required"`
	Name      string    `json:"name" validade:"required"`
	Slug      string    `json:"slug"`
	Category  string    `json:"category" validade:"required"`
	Address   string    `json:"address,omitempty"`
	Lat       *float64  `json:"lat,omitempty"`
	Lng       *float64  `json:"lng,omitempty"`
	QRCode    string    `json:"qr_code,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *Establishment) Validate() error {
	if err := validator.Validate(e); err != nil {
		return customerror.NewValidationError(err.Error())
	}
	return nil
}
