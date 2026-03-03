package domain_fakes

import (
	"fivestars/internal/domain"
	"time"
)

type EstablishmentBuilder struct {
	*domain.Builder[domain.Establishment]
}

func NewEstablishmentBuilder() *EstablishmentBuilder {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	lat := -23.55052
	lng := -46.633308
	establishment := &domain.Establishment{
		ID:        "22222222-2222-4222-8222-222222222222",
		Name:      "Cafe Central",
		Slug:      "cafe-central",
		Category:  "cafe",
		Address:   "Av. Paulista, 1000",
		Lat:       &lat,
		Lng:       &lng,
		QRCode:    "qr-cafe-central",
		CreatedAt: now,
		UpdatedAt: now,
	}

	return &EstablishmentBuilder{Builder: domain.NewBuilder[domain.Establishment](*establishment)}
}

func (b *EstablishmentBuilder) WithID(id string) *EstablishmentBuilder {
	b.Builder.Value.ID = id
	return b
}

func (b *EstablishmentBuilder) WithName(name string) *EstablishmentBuilder {
	b.Builder.Value.Name = name
	return b
}

func (b *EstablishmentBuilder) WithCategory(category string) *EstablishmentBuilder {
	b.Builder.Value.Category = category
	return b
}

func (b *EstablishmentBuilder) WithLocation(lat, lng float64) *EstablishmentBuilder {
	b.Builder.Value.Lat = &lat
	b.Builder.Value.Lng = &lng
	return b
}

func (b *EstablishmentBuilder) WithoutLocation() *EstablishmentBuilder {
	b.Builder.Value.Lat = nil
	b.Builder.Value.Lng = nil
	return b
}
