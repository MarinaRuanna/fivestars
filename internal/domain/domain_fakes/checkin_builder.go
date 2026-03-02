package domain_fakes

import (
	"fivestars/internal/domain"
	"time"
)

type CheckinBuilder struct {
	*domain.Builder[domain.Checkin]
}

func NewCheckinBuilder() *CheckinBuilder {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	checkin := &domain.Checkin{
		ID:              "33333333-3333-4333-8333-333333333333",
		UserID:          "11111111-1111-4111-8111-111111111111",
		EstablishmentID: "22222222-2222-4222-8222-222222222222",
		Lat:             -23.55052,
		Lng:             -46.633308,
		CheckedAt:       now,
		CreatedAt:       now,
	}

	return &CheckinBuilder{Builder: domain.NewBuilder[domain.Checkin](*checkin)}
}

func (b *CheckinBuilder) WithUserID(userID string) *CheckinBuilder {
	b.Builder.Value.UserID = userID
	return b
}

func (b *CheckinBuilder) WithEstablishmentID(establishmentID string) *CheckinBuilder {
	b.Builder.Value.EstablishmentID = establishmentID
	return b
}

func (b *CheckinBuilder) WithLocation(lat, lng float64) *CheckinBuilder {
	b.Builder.Value.Lat = lat
	b.Builder.Value.Lng = lng
	return b
}

func (b *CheckinBuilder) WithTimestamps(checkedAt, createdAt time.Time) *CheckinBuilder {
	b.Builder.Value.CheckedAt = checkedAt
	b.Builder.Value.CreatedAt = createdAt
	return b
}

func (b *CheckinBuilder) WithoutTimestamps() *CheckinBuilder {
	b.Builder.Value.CheckedAt = time.Time{}
	b.Builder.Value.CreatedAt = time.Time{}
	return b
}
