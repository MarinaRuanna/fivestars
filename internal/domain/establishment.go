package domain

import (
	"context"
	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
	"time"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/establishment_repository.go -package mock_domain . EstablishmentRepository
type EstablishmentRepository interface {
	List(ctx context.Context) ([]Establishment, error)
}

type Establishment struct {
	ID        string `validate:"required"`
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
