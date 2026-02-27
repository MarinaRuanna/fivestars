package controller

import (
	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/pkg/validator"
	"time"
)

type createCheckinDTO struct {
	EstablishmentID string   `json:"establishment_id" validate:"required,uuid4"`
	Lat             *float64 `json:"lat" validate:"required"`
	Lng             *float64 `json:"lng" validate:"required"`
}

type checkinResponseDTO struct {
	ID              string    `json:"checkin_id"`
	UserID          string    `json:"user_id"`
	EstablishmentID string    `json:"establishment_id"`
	Lat             float64   `json:"lat"`
	Lng             float64   `json:"lng"`
	CheckedAt       time.Time `json:"checked_at"`
}

func (d *checkinResponseDTO) Validate() error {
	if err := validator.Validate(d); err != nil {
		return customerror.NewValidationError(err.Error())
	}
	return nil
}

func ToDomainCheckin(dto *createCheckinDTO, userID string) (*domain.Checkin, error) {
	if dto.Lat == nil || dto.Lng == nil {
		return nil, customerror.NewValidationError("lat/lng required")
	}

	now := time.Now().UTC()

	checkin := &domain.Checkin{
		UserID:          userID,
		EstablishmentID: dto.EstablishmentID,
		Lat:             *dto.Lat,
		Lng:             *dto.Lng,
		CheckedAt:       now,
	}

	err := checkin.Validate()
	if err != nil {
		return nil, err
	}

	return checkin, nil

}

func ToCheckinDTO(checkin *domain.Checkin) (checkinResponseDTO, error) {
	dto := checkinResponseDTO{
		ID:              checkin.ID,
		UserID:          checkin.UserID,
		EstablishmentID: checkin.EstablishmentID,
		Lat:             checkin.Lat,
		Lng:             checkin.Lng,
		CheckedAt:       checkin.CheckedAt,
	}
	if err := dto.Validate(); err != nil {
		return checkinResponseDTO{}, err
	}
	return dto, nil
}
