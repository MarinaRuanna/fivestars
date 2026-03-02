package domain

import (
	"context"
	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
	"time"
)

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/checkin_repository.go -package mock_domain . CheckinRepository
type CheckinRepository interface {
	Create(ctx context.Context, c *Checkin) error
	ListByUser(ctx context.Context, userID string) ([]Checkin, error)
	FindTodayByUserAndEstablishment(ctx context.Context, userID, establishmentID string) (*Checkin, error)
}

type Checkin struct {
	ID              string
	UserID          string  `validate:"required,uuid4"`
	EstablishmentID string  `validate:"required,uuid4"`
	Lat             float64 `validate:"required"`
	Lng             float64 `validate:"required"`
	CheckedAt       time.Time
	CreatedAt       time.Time
}

func (c *Checkin) Validate() error {
	if err := validator.Validate(c); err != nil {
		return customerror.NewValidationError(err.Error())
	}
	return nil
}
