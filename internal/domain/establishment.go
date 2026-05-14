package domain

import (
	"context"
	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
	"time"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/establishment_repository.go -package mock_domain . EstablishmentRepository
type EstablishmentRepository interface {
	Create(ctx context.Context, establishment *Establishment) error
	List(ctx context.Context) ([]Establishment, error)
	GetByID(ctx context.Context, id string) (*Establishment, error)
	GetStats(ctx context.Context, id string) (*EstablishmentStats, error)
	DistanceTo(ctx context.Context, id string, lat, lng float64) (float64, error)
}

type Establishment struct {
	ID        string `validate:"omitempty,uuid4"`
	OwnerID   string `validate:"omitempty,uuid4"`
	Name      string `validate:"required"`
	Slug      string
	Category  string `validate:"required"`
	Address   string
	Lat       *float64
	Lng       *float64
	QRCode    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *Establishment) Validate() error {
	if err := validator.Validate(e); err != nil {
		return customerror.NewValidationError(err.Error())
	}
	return nil
}

func NewEstablishment(input Establishment) (*Establishment, error) {
	establishment := input
	if err := establishment.Validate(); err != nil {
		return nil, err
	}

	return &establishment, nil
}
