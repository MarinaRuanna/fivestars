package domain_fakes

import (
	"fivestars/internal/domain"
	"time"
)

type ReviewBuilder struct {
	*domain.Builder[domain.Review]
}

func NewReviewBuilder() *ReviewBuilder {
	now := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
	review := &domain.Review{
		ID:              "44444444-4444-4444-8444-444444444444",
		UserID:          "11111111-1111-4111-8111-111111111111",
		EstablishmentID: "22222222-2222-4222-8222-222222222222",
		CheckinID:       "33333333-3333-4333-8333-333333333333",
		Rating:          5,
		Text:            "Atendimento excelente e ambiente ótimo.",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	return &ReviewBuilder{Builder: domain.NewBuilder[domain.Review](*review)}
}

func (b *ReviewBuilder) WithUserID(userID string) *ReviewBuilder {
	b.Builder.Value.UserID = userID
	return b
}

func (b *ReviewBuilder) WithCheckinID(checkinID string) *ReviewBuilder {
	b.Builder.Value.CheckinID = checkinID
	return b
}

func (b *ReviewBuilder) WithEstablishmentID(establishmentID string) *ReviewBuilder {
	b.Builder.Value.EstablishmentID = establishmentID
	return b
}

func (b *ReviewBuilder) WithRating(rating int) *ReviewBuilder {
	b.Builder.Value.Rating = rating
	return b
}

func (b *ReviewBuilder) WithText(text string) *ReviewBuilder {
	b.Builder.Value.Text = text
	return b
}

func (b *ReviewBuilder) WithoutTimestamps() *ReviewBuilder {
	b.Builder.Value.CreatedAt = time.Time{}
	b.Builder.Value.UpdatedAt = time.Time{}
	return b
}
